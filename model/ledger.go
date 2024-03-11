package model

import "time"

type Ledger struct {
	Id         string    `gorm:"id"`
	TransferId string    `gorm:"transaction_id"`
	MerchantId string    `gorm:"merchant_id"`
	Credit     int64     `gorm:"credit"`
	Debit      int64     `gorm:"debit"`
	Created    time.Time `gorm:"created"`
	Updated    time.Time `gorm:"updated"`
}
