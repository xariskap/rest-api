package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rest/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {

	baseURL := "http://localhost:8888/products"
	testCases := []struct {
		page           int
		limit          int
		name           string
		expectSuccess  bool
		expectedStatus int
	}{
		{page: 1, limit: 5, name: "Page 0, Limit 0", expectSuccess: true, expectedStatus: http.StatusOK},
		{page: 2, limit: 10, name: "Page 2, Limit 10", expectSuccess: true, expectedStatus: http.StatusOK},
		{page: 1, limit: 30, name: "Page 1, Limit 30", expectSuccess: true, expectedStatus: http.StatusOK},
		{page: 3, limit: 15, name: "Page 3, Limit 15 (Expected to Fail)", expectSuccess: false, expectedStatus: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			url := fmt.Sprintf("%s?page=%d&limit=%d", baseURL, tc.page, tc.limit)

			resp, err := http.Get(url)
			if err != nil {
				log.Fatalf("Error making GET request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected HTTP status to match")

			if tc.expectSuccess {
				var products []utils.Product
				if err := json.Unmarshal(body, &products); err != nil {
					t.Fatalf("Error unmarshalling JSON: %v", err)
				}

				assert.Equal(t, tc.limit, len(products))
			}
		})
	}
}

func TestGetProduct(t *testing.T) {

	baseURL := "http://localhost:8888/product"
	testCases := []struct {
		id             string
		name           string
		expectSuccess  bool
		expectedStatus int
	}{
		{id: "1", name: "id 1", expectSuccess: true, expectedStatus: http.StatusOK},
		{id: "ena", name: "id ena", expectSuccess: false, expectedStatus: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			url := fmt.Sprintf("%s?id=%s", baseURL, tc.id)
			expectedProduct := utils.NewProduct("1", "Apple iPhone 15", "999.99", "150")

			resp, err := http.Get(url)
			if err != nil {
				log.Fatalf("Error making GET request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected HTTP status to match")

			if tc.expectSuccess {
				var product utils.Product
				if err := json.Unmarshal(body, &product); err != nil {
					t.Fatalf("Error unmarshalling JSON: %v", err)
				}

				assert.Equal(t, expectedProduct, product)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	p1 := utils.NewProduct("1", "updated", "updated", "updated")
	p2 := utils.NewProduct("2", "updated", "", "")
	p3 := utils.NewProduct("3", "", "updated", "")
	p4 := utils.NewProduct("4", "", "", "updated")
	p5 := utils.NewProduct("", "updated", "updated", "updated")
	p6 := utils.NewProduct("DoesNotExist", "updated", "updated", "updated")

	expectedp1 := utils.NewProduct("1", "updated", "updated", "updated")
	expectedp2 := utils.NewProduct("2", "updated", "849.99", "200")
	expectedp3 := utils.NewProduct("3", "Sony PlayStation 5", "updated", "75")
	expectedp4 := utils.NewProduct("4", "Dell XPS 13 Laptop", "1199.99", "updated")

	var expectedProduct []utils.Product = []utils.Product{
		expectedp1,
		expectedp2,
		expectedp3,
		expectedp4,
	}

	baseURL := "http://localhost:8888/product"
	testCases := []struct {
		p              utils.Product
		name           string
		expectSuccess  bool
		expectedStatus int
	}{
		{p: p1, name: "Update product with id 1", expectSuccess: true, expectedStatus: http.StatusOK},
		{p: p2, name: "Update product with id 2", expectSuccess: true, expectedStatus: http.StatusOK},
		{p: p3, name: "Update product with id 3", expectSuccess: true, expectedStatus: http.StatusOK},
		{p: p4, name: "Update product with id 4", expectSuccess: true, expectedStatus: http.StatusOK},
		{p: p5, name: "Fail update id not given", expectSuccess: false, expectedStatus: http.StatusBadRequest},
		{p: p6, name: "Fail update id not exist", expectSuccess: false, expectedStatus: http.StatusInternalServerError},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			url := fmt.Sprintf("%s?id=%s", baseURL, tc.p.ID)

			jsonData, err := json.Marshal(tc.p)
			if err != nil {
				log.Fatalf("Error marshalling product to JSON: %v", err)
			}

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

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected HTTP status to match")

			if tc.expectSuccess && i < 4 {
				var product utils.Product
				if err := json.Unmarshal(body, &product); err != nil {
					t.Fatalf("Error unmarshalling JSON: %v", err)
				}

				assert.Equal(t, expectedProduct[i], product)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {

	baseURL := "http://localhost:8888/product"
	testCases := []struct {
		id             string
		name           string
		expectSuccess  bool
		expectedStatus int
	}{
		{id: "1", name: "Delete id 1", expectSuccess: true, expectedStatus: http.StatusNoContent},
		{id: "", name: "Delete, id not given", expectSuccess: false, expectedStatus: http.StatusBadRequest},
		{id: "DoesNotExist", name: "Delete, id does not exist", expectSuccess: false, expectedStatus: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			url := fmt.Sprintf("%s?id=%s", baseURL, tc.id)

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

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected HTTP status to match")

		})
	}
}

func TestAddProduct(t *testing.T) {

	p1 := utils.NewProduct("31", "NAME", "PRICE", "QUANTITY")
	p2 := utils.NewProduct("", "NAME", "PRICE", "QUANTITY")
	p3 := utils.NewProduct("31", "NAME", "PRICE", "QUANTITY")
	type EmptyStruct struct{}
	var p4 EmptyStruct

	baseURL := "http://localhost:8888/products"
	testCases := []struct {
		p              utils.Product
		name           string
		emptyProduct   EmptyStruct
		expectSuccess  bool
		expectedStatus int
	}{
		{p: p1, name: "Add product 1", expectSuccess: true, expectedStatus: http.StatusCreated},
		{p: p2, name: "Add product 2. empty field", expectSuccess: false, expectedStatus: http.StatusBadRequest},
		{p: p3, name: "Add product 3, id already in db", expectSuccess: false, expectedStatus: http.StatusConflict},
		{emptyProduct: p4, name: "Add product 4. empty struct", expectSuccess: false, expectedStatus: http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			jsonData, err := json.Marshal(tc.p)
			if err != nil {
				log.Fatalf("Error marshalling product to JSON: %v", err)
			}

			resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Fatalf("Error making POST request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "Expected HTTP status to match")

			if tc.expectSuccess {
				var product utils.Product
				if err := json.Unmarshal(body, &product); err != nil {
					t.Fatalf("Error unmarshalling JSON: %v", err)
				}

				assert.Equal(t, p1, product)
			}

		})
	}
}
