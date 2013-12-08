package coinsummary

import (
	"appengine"
	"appengine/urlfetch"
	"coinsummary/models"
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	coinbaseApiUrl string = "https://coinbase.com/api/v1/prices/spot_rate"
	mtGoxApiUrl    string = "http://data.mtgox.com/api/2/BTCUSD/money/ticker_fast"
)

type CoinbaseResponse struct {
	Amount   float64 `json:",string"`
	Currency string
}

func GetCoinbasePrice(context appengine.Context) (*models.Price, error) {
	client := urlfetch.Client(context)
	resp, err := client.Get(coinbaseApiUrl)

	if err != nil {
		return new(models.Price), err
	}

	body, err := ioutil.ReadAll(resp.Body)

	var cr CoinbaseResponse
	err = json.Unmarshal(body, &cr)

	if err != nil {
		return new(models.Price), err
	}

	p := models.Price{Exchange: models.COINBASE, Price: cr.Amount, Currency: cr.Currency, Date: time.Now()}
	return &p, nil
}

type MtGoxResponse struct {
	Result string
	Data   map[string]MtGoxPriceObject
}

type MtGoxPriceObject struct {
	Value        float64 `json:",string"`
	ValueInt     int     `json:"value_int,string"`
	Display      string
	DisplayShort string `json:"display_short"`
	Currency     string
}

func GetMtGoxPrice(context appengine.Context) (*models.Price, error) {
	client := urlfetch.Client(context)
	resp, err := client.Get(mtGoxApiUrl)

	if err != nil {
		return new(models.Price), err
	}

	body, err := ioutil.ReadAll(resp.Body)

	var response MtGoxResponse
	err = json.Unmarshal(body, &response)

	if err == nil {
		return new(models.Price), err
	}

	p := models.Price{Exchange: models.MT_GOX, Price: response.Data["last"].Value, Currency: response.Data["last"].Currency, Date: time.Now()}
	return &p, nil
}
