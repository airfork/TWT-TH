package main

// carSale is for mapping the api json into a struct
type carSale struct {
	ID            int
	ImportCountry string `json:"import_country"`
	Model         string
	Make          string
	SoldBy        string `json:"sold_by"`
	Price         int    `json:"sale_price"`
}

type countryMap struct {
	Country            string         // Name of country
	Cars               map[string]int // Maps make and model to num occurrences
	Makes              map[string]int // Maps makes to num occurrences
	Sellers            map[string]int // Maps sellers to num occurrences
	TotalSales         *int           // Total of things sold
	TotalSalesString   string         // String version of TotalSales
	QuantitySold       *int           // Number of cars sold
	QuantitySoldString string         // String version of QuantitySold
	BestSelling        string         // best selling car
}

type countryMappings struct {
	Countries map[string]countryMap // Maps countries to their country maps
}
