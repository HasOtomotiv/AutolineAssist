package main

import (
	"AutolineAssist/internal/db"
	_ "AutolineAssist/statik"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/alexbrainman/odbc"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	AutolineDS       string
	Port         uint32

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
	soStr := ""
	if fileType == "so" {
		for _, poRecord := range poRecords {
			partNumber := strings.Replace(poRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s\t%.0f\t%.0f\n", partNumber, poRecord.QuantityRequired, poRecord.WIPLineOrPONo)
		}
	}
	if fileType == "csv" {
		for _, poRecord := range poRecords {
			partNumber := strings.Replace(poRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s;%.0f;%.0f\n", partNumber, poRecord.QuantityRequired, poRecord.WIPLineOrPONo)
		}
	}

	if fileType == "xfr" {
		soStr = soStr + "Assist.XFR|AutolineAssist||\n"
		for _, poRecord := range poRecords {
			partNumber := strings.Replace(poRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s||%03.0f|||\n", partNumber, poRecord.QuantityRequired)
		}
	}

	return ctx.Blob(http.StatusOK, "text/csv", []byte (soStr))
}

func WIPRecords(ctx echo.Context) (err error) {
	fileType := ctx.Param("Type")
	loc := ctx.Param("Loc")
	wipNumber := ctx.Param("WIPNumber")

	fmt.Printf("Start %s\n",fileType)

	wipRecords, err := db.GetWIPRecords(loc, wipNumber)


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
		for _, wipRecord := range wipRecords {
			partNumber := strings.Replace(wipRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s;%.0f;%.0f\n", partNumber, wipRecord.OrderQuantity, wipRecord.WIPNumber)
		}
	}

	if fileType == "xfr" {
		soStr = soStr + "Assist.XFR|AutolineAssist||\n"
		for _, wipRecord := range wipRecords {
			partNumber := strings.Replace(wipRecord.PartNumber, "M", "", -1)
			partNumber = strings.Replace(partNumber, " ", "", -1)
			partNumber = strings.Replace(partNumber, "/", "", -1)
			soStr = soStr + fmt.Sprintf("%s||%03.0f|||\n", partNumber,  wipRecord.OrderQuantity)
		}
	}

	return ctx.Blob(http.StatusOK, "text/csv", []byte (soStr))
}

func main() {
	var err error
	var cfg Config

	fmt.Printf("Start %s\n", time.Now().Format(time.RFC3339))

	if _, err := toml.DecodeFile("autolineassist.toml", &cfg); err != nil {
		log.Println(err)
		os.Exit(2)
		return
	}


	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	autolineDS := cfg.AutolineDS
	port:= cfg.Port
	if port==0 {
		port=8080
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

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d",port)))
}

