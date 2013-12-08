package models

import "time"

const (
	COINBASE = "Coinbase"
	MT_GOX   = "Mt.Gox"
)

type Price struct {
	Exchange string
	Price    float64
	Currency string
	Date     time.Time
}
