package main

import (
	// "errors"
	"fmt"
	"golang.org/x/text/message"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

var tpl = template.Must(template.ParseGlob("views/*.html"))

func main() {

	// Setup server to run on port 8000
	r := mux.NewRouter()
	srv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":8000",
	}
	// Two routes
	r.HandleFunc("/", index)
	r.HandleFunc("/data", sendCountryJSON)
	// Allow assets and views directory to be served
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	r.PathPrefix("/views/").Handler(http.StripPrefix("/views/", http.FileServer(http.Dir("views/"))))
	// Use gorilla mux to handle requests and start up server
	http.Handle("/", r)
	fmt.Println("Server started on port 8000")
	log.Fatal(srv.ListenAndServe())
}

// mapCountries makes the country mappings data structure
func mapCountries(sales []carSale) countryMappings {
	// Create mappings struct
	var mapOfCountries countryMappings
	// Assign map to struct
	mapOfCountries.Countries = make(map[string]countryMap)
	for _, sale := range sales {
		// See if our import country exists in our map
		cMap, ok := mapOfCountries.Countries[sale.ImportCountry]
		if ok {
			updateCountryMap(sale, cMap)
		} else {
			// Create countryMap and initialize underlying maps
			var tempCMap countryMap
			tempCMap.Country = sale.ImportCountry
			tempCMap.Cars = make(map[string]int)
			tempCMap.Makes = make(map[string]int)
			tempCMap.Sellers = make(map[string]int)
			// Create int and pass reference to it to my map
			var total, quantity int
			tempCMap.TotalSales = &total
			tempCMap.QuantitySold = &quantity
			// Map import country of sale to new country map and update to reflect this sale
			mapOfCountries.Countries[sale.ImportCountry] = tempCMap
			updateCountryMap(sale, tempCMap)
		}
	}
	return mapOfCountries
}

// Updates the parameters of the sale country's countryMap
func updateCountryMap(sale carSale, cMap countryMap) {
	// Update total sales and the number of sold vehicles
	*cMap.TotalSales = sale.Price + *cMap.TotalSales
	*cMap.QuantitySold++
	// Check each map to see if elements in sale are unique
	// If so, create entry for them, otherwise increment existing entry
	_, ok := cMap.Cars[sale.Make+" "+sale.Model]
	if !ok {
		cMap.Cars[sale.Make+" "+sale.Model] = 1
	} else {
		cMap.Cars[sale.Make+" "+sale.Model]++
	}
	_, ok = cMap.Makes[sale.Make]
	if !ok {
		cMap.Makes[sale.Make] = 1
	} else {
		cMap.Makes[sale.Make]++
	}
	_, ok = cMap.Sellers[sale.SoldBy]
	if !ok {
		cMap.Sellers[sale.SoldBy] = 1
	} else {
		cMap.Sellers[sale.SoldBy]++
	}
}

// bestSellingCar returns the best selling car in a country
func (cm countryMap) bestSellingCar() (string, int) {
	var (
		max int
		car string
	)
	for key, value := range cm.Cars {
		if value > max {
			max = value
			car = key
		}
	}
	return car, max
}

// countriesByRevenue returns a slice of countryMaps sorted by total revenue
func (cms countryMappings) countriesByRevenue() []countryMap {
	countries := make([]countryMap, len(cms.Countries))
	var i int
	for _, value := range cms.Countries {
		countries[i] = value
		best, _ := countries[i].bestSellingCar()
		countries[i].BestSelling = best
		p := message.NewPrinter(message.MatchLanguage("en"))
		countries[i].TotalSalesString = p.Sprint(*countries[i].TotalSales)
		countries[i].QuantitySoldString = p.Sprint(*countries[i].QuantitySold)
		i++
	}
	less := func(i, j int) bool {
		return *countries[i].TotalSales > *countries[j].TotalSales
	}
	sort.Slice(countries, less)
	return countries
}
