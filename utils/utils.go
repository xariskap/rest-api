package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

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
