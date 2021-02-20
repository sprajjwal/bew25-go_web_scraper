package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Crypto struct {
	Name   string
	Symbol string
	Price  string
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("["); err != nil {
		panic(err)
	}

	// On every a element which has href attribute call callback
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		name, symbol, price := "", "", ""

		e.ForEach("td", func(i int, elem *colly.HTMLElement) {

			if i == 3 {
				price = elem.Text
			} else if i == 2 {
				elem.ForEach("p", func(i int, p *colly.HTMLElement) {
					if i == 0 {
						name = p.Text
					}
					if i == 1 {
						symbol = p.Text
					}
				})
			}
		})

		if name != "" {
			coin := Crypto{
				Name:   name,
				Symbol: symbol,
				Price:  price,
			}
			coinJSON, _ := json.MarshalIndent(coin, "", " ")
			fmt.Println(string(coinJSON))

			if _, err = f.WriteString(string(coinJSON)); err != nil {
				panic(err)
			}
			if _, err = f.WriteString(","); err != nil {
				panic(err)
			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://coinmarketcap.com/")
	if _, err = f.WriteString("]"); err != nil {
		panic(err)
	}
}
