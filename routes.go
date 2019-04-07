package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// index renders a page with only basic content
func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		out := fmt.Sprintln("Something went wrong, please try again")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(out))
		return
	}
}

// sendCountryJSON makes API request and sends JSON to client
func sendCountryJSON(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://my.api.mockaroo.com/tunji.json?key=e6ac1da0")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	// Read from response body
	salesInput, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Marshal response into a slice of car sales
	var sales []carSale
	_ = json.Unmarshal(salesInput, &sales)
	// Create mapping so that information about sales is mapped to a country
	m := mapCountries(sales)
	// Create slice of countryMaps ordered by total revenue
	out := m.countriesByRevenue()
	// Turn slice into JSON and send to user
	j, err := json.Marshal(out)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
