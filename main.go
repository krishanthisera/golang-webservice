package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ProductID    int    `json:"prductId"`
	Manufacturer string `json:"manufacturer"`
	PricePerUnit int    `json:"pricePerUnit"`
	ProductName  string `json:"productName"`
	Quantity     int    `json:"qty"`
}

var productList []Product

func init() {
	productsJSON := `[
		{
			"prductId": 1,
			"manufacturer": "Sony",
			"pricePerUnit": 45,
			"productName": "32inch TV",
			"qty": 5
		},
		{
			"prductId": 2,
			"manufacturer": "Sony",
			"pricePerUnit": 544,
			"productName": "PS4 PRO",
			"qty": 45
		},
		{
			"prductId": 3,
			"manufacturer": "Apple",
			"pricePerUnit": 999,
			"productName": "29inch iMac",
			"qty": 55
		}
	]`

	err := json.Unmarshal([]byte(productsJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1

}
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//Add new product
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusConflict)
			return
		}

		newProduct.ProductID = getNextID()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodGet:
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)
	}
}

// To search a prodcuct and update
func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductById(productID)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var updateProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updateProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if updateProduct.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product = &updateProduct
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodGet:
		productsJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJSON)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func findProductById(ProductID int) (*Product, int) {
	for i, product := range productList {
		if product.ProductID == ProductID {
			return &product, i
		}
	}
	return nil, 0
}

func middlewareHandler(handler http.Handler) http.Handler(){
	return http.HandlerFunc(func (w http.ResponseWriter, r http.Request)){
		fmt.Println("Before Handler; Middleware start")
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Println("Middleware executed; %s", time.Since(start) )
	}
}
func main() {

	prodcuctListHandler := http.HandlerFunc(productsHandler)
	productItemHandler 	:= http.HandlerFunc(productHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":8004", nil)

}
