package main

import (
	"AutolineAssist/internal/db"
	_ "AutolineAssist/statik"
	"fmt"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/alexbrainman/odbc"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AutolineDS string
	Port       uint32
}

func Right(str string, num int) string {
	if num <= 0 {
		return ``
	}
	max := len(str)
	if num > max {
		num = max
	}
	num = max - num
	return str[num:]
}

func einvoiceinfo(ctx echo.Context) (err error) {
	docNumber := ctx.Param("DocNumber")
	ii, err := db.GetInvoiceInfo(docNumber)

	if err == nil {
		return ctx.JSON(http.StatusOK, echo.Map{"data": ii, "err": err})
	} else {
		return ctx.JSON(http.StatusOK, echo.Map{"data": nil, "err": err})

	}
}

func PORecords(ctx echo.Context) (err error) {

	fileType := ctx.Param("Type")
	loc := ctx.Param("Loc")
	poNumber := ctx.Param("PONumber")

	poRecords, err := db.GetPORecords(loc, poNumber)
	if err != nil {
		log.Errorf("error retrieving PORecords from DB, Error: %s", err.Error())
	}

	log.Infof("Kayıt sayısı: %d", len(poRecords))

	soStr := ""

	chassisNumbers := make(map[int32]string)
	for _, poRecord := range poRecords {
		headerLogMagic := poRecord.HeaderLogMagic
		chassisNumber := ""
		if headerLogMagic > 0 {
			if _, ok := chassisNumbers[headerLogMagic]; !ok {
				chassisNumber, err = db.GetChassisNumber(loc, headerLogMagic)
				log.Printf("chassisNumber: %s, headerLogMagic: %v\n", chassisNumber, headerLogMagic)
				if len(chassisNumber) < 17 {
					chassisNumber = ""
				}
				if len(chassisNumber) > 17 {
					chassisNumber = chassisNumber[:17]
				}
				if err != nil {
					chassisNumber = ""
					log.Errorf("error retrieving chassis number from DB, Error: %s", err.Error())
				} else {
					chassisNumbers[headerLogMagic] = chassisNumber
				}
			} else {
				chassisNumber = chassisNumbers[headerLogMagic]
			}
		}
		log.Printf("chassisNumber: %s, headerLogMagic: %v\n", chassisNumber, headerLogMagic)

		partNumber := strings.TrimSpace(poRecord.PartNumber)
		if partNumber[0] == 'M' {
			partNumber = partNumber[1:]
		}
		partNumber = strings.Replace(partNumber, " ", "", -1)
		partNumber = strings.Replace(partNumber, "/", "", -1)

		if fileType == "csv" {
			soStr = soStr + fmt.Sprintf("%s;%.0f;%s;%.6s\r\n", partNumber, poRecord.QuantityRequired, chassisNumber, fmt.Sprintf("%.0f", poRecord.WIPLineOrPONo))
		}
		if fileType == "so" {
			soStr = soStr + fmt.Sprintf("%s\t%.0f\t%.0f\n", partNumber, poRecord.QuantityRequired, poRecord.WIPLineOrPONo)
		}
		if fileType == "xfr" {
			soStr = soStr + fmt.Sprintf("%s||%03.0f|||\n", partNumber, poRecord.QuantityRequired)
		}
	}
	clear(chassisNumbers)
	soHederStr := ""
	if fileType == "xfr" {
		soHederStr = "Assist.XFR|AutolineAssist||\n"
	}
	if fileType == "csv" {
		soHederStr = "Material;Quantity;Vin Number;PO item\r\n"
	}

	return ctx.Blob(http.StatusOK, "text/csv", []byte(soHederStr+soStr))
}

func WIPRecords(ctx echo.Context) (err error) {
	fileType := ctx.Param("Type")
	loc := ctx.Param("Loc")
	wipNumber := ctx.Param("WIPNumber")

	log.Infof("Start %s\n", fileType)

	wipRecords, err := db.GetWIPRecords(loc, wipNumber)
	if err != nil {
		log.Errorf("error retrieving WIPRecords from DB, Error: %s", err.Error())
	}

	soStr := ""
	if fileType == "so" {
		for _, wipRecord := range wipRecords {
			partNumber := strings.Replace(wipRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s\t%.0f\t%.0f\n", partNumber, wipRecord.OrderQuantity, wipRecord.WIPNumber)
		}
	}
	if fileType == "csv" {
		soStr = "Material;Quantity;PO item\r\n"

		for _, wipRecord := range wipRecords {
			partNumber := strings.Replace(wipRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s;%.0f;%s\r\n", partNumber, wipRecord.OrderQuantity, Right(fmt.Sprintf("%.0f", wipRecord.WIPNumber), 6))
		}
	}

	if fileType == "xfr" {
		soStr = "Assist.XFR|AutolineAssist||\n"
		for _, wipRecord := range wipRecords {
			partNumber := strings.Replace(wipRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s||%03.0f|||\n", partNumber, wipRecord.OrderQuantity)
		}
	}

	return ctx.Blob(http.StatusOK, "text/csv", []byte(soStr))
}

func main() {
	var err error
	var cfg Config

	log.SetLevel(log.WarnLevel)

	log.Infof("Start %s\n", time.Now().Format(time.RFC3339))

	if _, err := toml.DecodeFile("autolineassist.toml", &cfg); err != nil {
		log.Errorf("Error Decode autolineassist.toml :", err.Error())
		os.Exit(2)
		return
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	autolineDS := cfg.AutolineDS
	port := cfg.Port
	if port == 0 {
		port = 8080
	}

	db.Db, err = sqlx.Open("odbc", autolineDS)
	if err != nil {
		log.Fatal("Open ODBC failed (%v).", err)
	}
	defer db.Db.Close()

	e := echo.New()
	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderOrigin, "Cache-Control", echo.HeaderXRequestedWith, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	api := e.Group("/api")
	api.GET("/einvoiceinfo/:DocNumber", einvoiceinfo)
	api.GET("/getporecords/:Type/:Loc/:PONumber", PORecords)
	api.GET("/getwiprecords/:Type/:Loc/:WIPNumber", WIPRecords)
	e.GET("/*", echo.WrapHandler(http.FileServer(statikFS)))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
