package models

import (
	"appengine"
	"appengine/datastore"
)

type PriceRepository struct {
	Context appengine.Context
}

func (repository *PriceRepository) GetLastPrice(exchange string) *Price {
	results := datastore.NewQuery("prices").Filter("Exchange=", exchange).Order("-Date").Limit(1)

	var prices []*Price
	results.GetAll(repository.Context, &prices)

	if len(prices) == 0 {
		return nil
	}

	return prices[0]
}

func (repository *PriceRepository) Add(price *Price) error {
	_, err := datastore.Put(repository.Context,
		datastore.NewIncompleteKey(repository.Context, "prices", nil),
		price)

	return err
}
