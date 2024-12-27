package db

import (
	"time"
)

type EInvoiceInfo struct {
	UniqueDocRef    string    `db:"UniqueDocRef"`
	InvoiceID       string    `db:"InvoiceID"`
	CustCompanyName string    `db:"CustCompanyName"`
	PostingDate     time.Time `db:"PostingDate"`
	UniqueDocKey    string    `db:"UniqueDocKey"`
	WIPNumber       int       `db:"WIPNumber"`
	InvoiceTotal    float32   `db:"InvoiceTotal"`
	ReceiverEmail   string    `db:"ReceiverEmail"`
}

type PORecord struct {
	WIPLineOrPONo    float64 `db:"WIPLineOrPONo"`
	QuantityRequired float64 `db:"QuantityRequired"`
	PartNumber       string  `db:"PartNumber"`
	HeaderLogMagic   int32   `db:"HeaderLogMagic001"`
}

type WIPRecord struct {
	WIPNumber     float64 `db:"WIPNumber"`
	OrderQuantity float64 `db:"OrderQuantity"`
	PartNumber    string  `db:"PartNumber"`
}
