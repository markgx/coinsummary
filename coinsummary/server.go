package coinsummary

import (
	"appengine"
	"coinsummary/models"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"net/http"
)

func init() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
	}))

	m.Get("/", rootHandler)
	m.Get("/e/:e", rootHandler)
	m.Get("/update", updateHandler)

	http.Handle("/", m)
}

func rootHandler(r render.Render, params martini.Params, req *http.Request) {
	c := appengine.NewContext(req)

	priceRepository := models.PriceRepository{Context: c}

	var price *models.Price

	switch params["e"] {
	case "mtgox":
		price = priceRepository.GetLastPrice(models.MT_GOX)
	case "coinbase":
		price = priceRepository.GetLastPrice(models.COINBASE)
	case "btc-e":
		price = priceRepository.GetLastPrice(models.BTCE)
	default:
		price = priceRepository.GetLastPrice(models.MT_GOX)
	}

	r.HTML(200, "home", price)
}

func updateHandler(params martini.Params, req *http.Request) string {
	c := appengine.NewContext(req)

	priceRepository := models.PriceRepository{Context: c}

	if p, err := GetPrice(c, models.COINBASE); err != nil {
		c.Errorf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		c.Errorf("Error %+v", rErr)
	}

	if p, err := GetPrice(c, models.MT_GOX); err != nil {
		c.Errorf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		c.Errorf("Error %+v", rErr)
	}

	if p, err := GetPrice(c, models.BTCE); err != nil {
		c.Errorf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		c.Errorf("Error %+v", rErr)
	}

	return "ok"
}
