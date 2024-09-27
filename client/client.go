package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rest/utils"
)

// Makes a POST request to the server
func create(p utils.Product) {

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

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("POST request successful!")
	} else {
		fmt.Printf("POST request failed with status: %s\n", resp.Status)
	}
}

// Makes a GET request to the server
func read(page int, limit int) {

	baseURL := "http://localhost:8888/products"
	url := fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page, limit)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making Get request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {

		var products []utils.Product
		if err := json.Unmarshal(body, &products); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		for _, product := range products {
			fmt.Printf("ID: %s, Name: %s, Price: %s, Quantity: %s\n", product.ID, product.Name, product.Price, product.Quantity)
		}
	} else {
		fmt.Printf("Get request failed with status: %s\n", resp.Status)
	}
}

// Makes a PUT request to the server
func update(p utils.Product) {

	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error marshalling product to JSON: %v", err)
	}

	baseURL := "http://localhost:8888/product"
	url := fmt.Sprintf("%s?id=%s", baseURL, p.ID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making PUT request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {

		var product utils.Product
		if err := json.Unmarshal(body, &product); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		fmt.Printf("ID: %s, Name: %s, Price: %s, Quantity: %s\n", product.ID, product.Name, product.Price, product.Quantity)

	}

}

// Makes a DELETE request to the server
func delete(id string) {

	baseURL := "http://localhost:8888/product"
	url := fmt.Sprintf("%s?id=%s", baseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making DELETE request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Product deleted successfully")
	}

}

func addAllProductsToDB(path string) {
	products := utils.JsonToArray(path)

	for _, p := range products {
		create(p)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide an argument.")
		return
	}

	arg := os.Args[1]

	switch arg {
	case "db":
		fmt.Println("Adding data to DB")
		addAllProductsToDB("../data.json")
	case "c":
		fmt.Println("Making a POST request")
		create(utils.NewProduct("0", "TEST", "POLLA LEFTA", "FULL"))
	case "r":
		fmt.Println("Making a GET request")
		read(1, 30) // RETURNS ALL PRODUCTS
	case "u":
		fmt.Println("Making a PUT request")
		update(utils.NewProduct("0", "TEST", "TEST", "TEST")) // !!! MAKE SURE TO USE A VALID ID !!!
	case "d":
		fmt.Println("Making a DELETE request")
		delete("1") // !!! MAKE SURE TO USE A VALID ID !!!
	default:
		fmt.Println("Unknown request. Please use 'db', 'c', 'r', 'u', or 'd'.")
	}
}
