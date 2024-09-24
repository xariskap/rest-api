package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest/db"
	"rest/utils"
)

var dbConn = db.GetConnection()
var dbName = "Simpler"

func createDB(name string) {
	database := db.Create(name, dbConn)
	fmt.Println(database.Name)
}

func addAllProductsToDB(path string) {
	products := utils.JsonToArray(path)

	for _, p := range products {
		jsonData, err := json.Marshal(p)
		if err != nil {
			log.Fatalf("Error marshalling product to JSON: %v", err)
		}

		url := "http://localhost:8888/products"

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
}

func main() {
	createDB(dbName)
	addAllProductsToDB("../data.json")
}
