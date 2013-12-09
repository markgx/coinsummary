package coinsummary

import (
	"appengine"
	"coinsummary/models"
	"fmt"
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
	m.Get("/update", updateHandler)

	http.Handle("/", m)
}

func rootHandler(r render.Render, params martini.Params, req *http.Request) {
	c := appengine.NewContext(req)

	priceRepository := models.PriceRepository{Context: c}
	price := priceRepository.GetLastPrice(models.MT_GOX)

	r.HTML(200, "home", price)
}

func updateHandler(params martini.Params, req *http.Request) string {
	c := appengine.NewContext(req)

	priceRepository := models.PriceRepository{Context: c}

	if p, err := GetCoinbasePrice(c); err != nil {
		return fmt.Sprintf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		return fmt.Sprintf("Error %+v", rErr)
	}

	if p, err := GetMtGoxPrice(c); err != nil {
		return fmt.Sprintf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		return fmt.Sprintf("Error %+v", rErr)
	}

	if p, err := GetBtcEPrice(c); err != nil {
		return fmt.Sprintf("Error %+v", err)
	} else if rErr := priceRepository.Add(p); rErr != nil {
		return fmt.Sprintf("Error %+v", rErr)
	}

	return "ok"
}
