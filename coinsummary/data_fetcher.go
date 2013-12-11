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
	mtGoxApiUrl    string = "http://data.mtgox.com/api/2/BTCUSD/money/ticker_fast"
	coinbaseApiUrl string = "https://coinbase.com/api/v1/prices/spot_rate"
	btcApiUrl      string = "https://btc-e.com/api/2/btc_usd/ticker"
)

func GetPrice(context appengine.Context, exchange string) (*models.Price, error) {
	switch exchange {
	case models.MT_GOX:
		body, err := fetchUrl(context, mtGoxApiUrl)

		if err != nil {
			return nil, err
		}

		var response MtGoxResponse
		err = json.Unmarshal(body, &response)

		if err != nil {
			return new(models.Price), err
		}

		p := models.Price{
			Exchange: models.MT_GOX,
			Price:    response.Data.Last.Value,
			Currency: response.Data.Last.Currency,
			Date:     time.Now(),
		}

		return &p, nil
	case models.COINBASE:
		body, err := fetchUrl(context, coinbaseApiUrl)

		if err != nil {
			return nil, err
		}

		var cr CoinbaseResponse
		err = json.Unmarshal(body, &cr)

		if err != nil {
			return nil, err
		}

		p := models.Price{
			Exchange: models.COINBASE,
			Price:    cr.Amount,
			Currency: cr.Currency,
			Date:     time.Now(),
		}

		return &p, nil
	case models.BTCE:
		body, err := fetchUrl(context, btcApiUrl)

		if err != nil {
			return nil, err
		}

		var response BtcEResponse
		err = json.Unmarshal(body, &response)

		if err != nil {
			return new(models.Price), err
		}

		p := models.Price{
			Exchange: models.BTCE,
			Price:    response.Ticker.Last,
			Currency: "USD",
			Date:     time.Unix(response.Ticker.Updated, 0),
		}

		return &p, nil
	}

	return nil, nil
}

func fetchUrl(context appengine.Context, url string) ([]byte, error) {
	client := urlfetch.Client(context)
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

type MtGoxResponse struct {
	Result string
	Data   MtGoxDataObject
}

type MtGoxDataObject struct {
	LastLocal MtGoxPriceObject `json:"last_local"`
	Last      MtGoxPriceObject
	LastOrig  MtGoxPriceObject `json:"last_orig"`
	LastAll   MtGoxPriceObject `json:"last_all"`
	Buy       MtGoxPriceObject
	Sell      MtGoxPriceObject
	Now       string
}

type MtGoxPriceObject struct {
	Value        float64 `json:",string"`
	ValueInt     int     `json:"value_int,string"`
	Display      string
	DisplayShort string `json:"display_short"`
	Currency     string
}

type CoinbaseResponse struct {
	Amount   float64 `json:",string"`
	Currency string
}

type BtcEResponse struct {
	Ticker BtcETickerObject
}

type BtcETickerObject struct {
	High       float64
	Low        float64
	Avg        float64
	Vol        float64
	VolCur     float64 `json:"vol_cur"`
	Last       float64
	Buy        float64
	Sell       float64
	Updated    int64
	ServerTime int64 `json:"server_time"`
}
