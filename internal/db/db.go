package db

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/encoding/charmap"
	"strings"
)

type DB = sqlx.DB

var Db *DB

func DecodeWindows1254(str string) string {
	sr := strings.NewReader(str)
	tr := charmap.Windows1254.NewDecoder().Reader(sr)
	buf := new(bytes.Buffer)
	buf.ReadFrom(tr)
	return buf.String()
}

func GetInvoiceInfo(DocNumber string) (EInvoiceInfo, error) {

	ii := EInvoiceInfo{}

	err := Db.QueryRowx(`select top 1 UniqueDocRef, InvoiceID, CustCompanyName, PostingDate, UniqueDocKey, WIPNumber, InvoiceTotal, ReceiverEmail 
								from CS_TR_EINVOICEHEADER 
								where InvoiceID<>'' and UniqueDocKey=? and PostingDate>='2019/01/01'`, DocNumber).StructScan(&ii)

	if err != nil {

		fmt.Printf("Read failed (%v).\n", err)
		return ii, err
	}

	ii.CustCompanyName = DecodeWindows1254(ii.CustCompanyName)

	return ii, err
}

func GetPORecords(Loc string, PONumber string) ([]*PORecord, error) {
	// select WIPLineOrPONo ,QuantityAdvised,QuantityReceived,QuantityRequired,CustomerOrderNo,CustomerCode ,PartNumber,AccountNumber,BINLocation
	var PORecords []*PORecord

	sqlQuery := fmt.Sprintf(`select WIPLineOrPONo, QuantityRequired, PartNumber
								from PC_%s_PurchaseTransactions
								where PurchaseOrderNo= ?
								order by PartNumberPacked`, Loc)

	err := Db.Select(&PORecords, sqlQuery, PONumber)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PORecords from DB, Error: %v", err)
	}

	return PORecords, err
}

func GetWIPRecords(Loc string, WIPNumber string) ([]*WIPRecord, error) {
	var WIPRecords []*WIPRecord
	// select WIPNumber, OrderQuantity, PartNumber, OrderStatus, PartsDiscountCode
	sqlQuery := fmt.Sprintf(`select WIPNumber, OrderQuantity, PartNumber
								from SO_%s_PartsLines
								where WIPNumber= ?
								order by LineNumber`, Loc)

	err := Db.Select(&WIPRecords, sqlQuery, WIPNumber)

	if err != nil {
		return nil, fmt.Errorf("error retrieving WIPRecords from DB, Error: %v", err)
	}

	return WIPRecords, err
}
