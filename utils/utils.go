package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

func NewProduct(id, name, price, quantity string) Product {
	return Product{id, name, price, quantity}
}

// Parses a json file and returns a list of Products
func JsonToArray(path string) []Product {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var products []Product

	err = json.Unmarshal(bytes, &products)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return products
}

func create(p Product) {

	url := "http://localhost:8888/products"

	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error marshalling product to JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("POST request failed with status: %s\n", resp.Status)
	}
}

func AddAllProductsToDB(path string) {
	products := JsonToArray(path)

	for _, p := range products {
		create(p)
	}
}
