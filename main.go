package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"rest/db"
	"rest/utils"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// Handles the get requests
func getProducts(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var totalProducts int
		err = db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Check if the offset is out of range
		if offset >= totalProducts {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page out of range"})
			return
		}

		rows, err := db.Query("SELECT id, name, price, quantity FROM products LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)

	}
}

// Get a specific product
func getProduct(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		var p Product

		err := db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

// Handles POST requests
func addProduct(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if product.ID == "" || product.Name == "" || product.Price == "" || product.Quantity == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty field"})
			return
		}

		err := db.ExecQuery("INSERT INTO products (id, name, price, quantity) VALUES ($1, $2, $3, $4)", product.ID, product.Name, product.Price, product.Quantity)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

// Handles PUT requests
func updatePruduct(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var dbProduct Product
		var newProduct Product
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is not provided!"})
			return
		}

		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&dbProduct.ID, &dbProduct.Name, &dbProduct.Price, &dbProduct.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if newProduct.Name == "" {
			newProduct.Name = dbProduct.Name
		}

		if newProduct.Price == "" {
			newProduct.Price = dbProduct.Price
		}

		if newProduct.Quantity == "" {
			newProduct.Quantity = dbProduct.Quantity
		}

		err = db.ExecQuery("UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4", newProduct.Name, newProduct.Price, newProduct.Quantity, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, newProduct)
	}
}

// Handles DELETE requests
func deleteProduct(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is not provided!"})
			return
		}

		// Product does not exist
		var p Product
		err := db.QueryRow("SELECT * FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = db.ExecQuery("DELETE FROM products WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}

}

func main() {
	dbUser := "root"
	dbPassword := "root"
	dbHost := os.Getenv("DB_HOST")
	dbPort := "26257"
	dbName := "restdb"

	if dbHost == "" {
		dbHost = "localhost"
	}

	for {
		_, err := db.GetConnection(dbUser, dbHost, dbPort, dbName)
		if err == nil {
			break
		}
		log.Println("Failed to connect to database. Retrying...")
		time.Sleep(1 * time.Second)
	}

	dbConn, _ := db.GetConnection(dbUser, dbHost, dbPort, dbName)
	db.Create(dbUser, dbPassword, dbHost, dbPort, dbName, dbConn)
	database := db.USE(dbUser, dbPassword, dbHost, dbPort, dbName, dbConn)
	defer dbConn.Close(context.Background())

	r := gin.Default()

	r.GET("/products", getProducts(database))
	r.GET("/product", getProduct(database))
	r.POST("/products", addProduct(database))
	r.PUT("/product", updatePruduct(database))
	r.DELETE("/product", deleteProduct(database))

	// add data to the database after server starts running
	go func() {
		time.Sleep(3 * time.Second)
		utils.AddAllProductsToDB("data.json")
	}()

	fmt.Println("Server is running on http://localhost:8888")
	if err := r.Run(":8888"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
