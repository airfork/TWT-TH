package main

import (
	"encoding/json"
	"fmt"
)

type carSale struct {
	ID            int
	ImportCountry string `json:"import_country"`
	Model         string
	Make          string
	SoldBy        string `json:"sold_by"`
	Price         int    `json:"sale_price"`
}

type countryMap struct {
	Cars       map[string]int // Maps make and model to num occurences
	Makes      map[string]int // Maps makes to num occurences
	Sellers    map[string]int // Maps sellers to num occurences
	TotalSales *int           // Total of things sold
}

type countryMappings struct {
	Countries map[string]countryMap // Maps countries to their country maps
}

func main() {
	// response, err := http.Get("https://my.api.mockaroo.com/tunji.json?key=e6ac1da0")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer response.Body.Close()

	// // Read from response body
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = ioutil.WriteFile("output.txt", body, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	salesInput := []byte(`[{"id":1,"import_country":"Brazil","model":"I","make":"Infiniti","sold_by":"Ruby Brimblecombe","sale_price":19497},{"id":2,"import_country":"Mongolia","model":"370Z","make":"Nissan","sold_by":"Richard Lowndesbrough","sale_price":17489}, {"id":3,"import_country":"Brazil","model":"250","make":"Infiniti","sold_by":"Ruby Brimblecombe","sale_price":19497}]`)
	var sales []carSale
	json.Unmarshal(salesInput, &sales)
	fmt.Println(sales[0].Price)
	m := mapCountries(sales)
	fmt.Printf("%+v\n", m)
	fmt.Println(*m.Countries["Brazil"].TotalSales)
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
			tempCMap.Cars = make(map[string]int)
			tempCMap.Makes = make(map[string]int)
			tempCMap.Sellers = make(map[string]int)
			// Create int and pass reference to it to my map
			var total int
			tempCMap.TotalSales = &total
			// Map import country of sale to new country map and update to reflect this sale
			mapOfCountries.Countries[sale.ImportCountry] = tempCMap
			updateCountryMap(sale, tempCMap)
		}
	}
	return mapOfCountries
}

// Updates the paramters of the sale country's countryMap
func updateCountryMap(sale carSale, cMap countryMap) {
	*cMap.TotalSales = sale.Price + *cMap.TotalSales
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

// // Return list of sorted
// func countryList(sales []carSale) []string {
// 	countries := make([]string, len(sales))
// 	for i, sale := range sales {
// 		countries[i] = sale.ImportCountry
// 	}
// 	less := func(i, j int) bool {
// 		return countries[i] < countries[j]
// 	}
// 	sort.Slice(countries, less)
// 	return countries
// }
