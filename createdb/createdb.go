package main

import (
	"fmt"
	"rest/db"
)

var dbConn = db.GetConnection()

func CreateDB(name string) {
	database := db.Create(name, dbConn)
	fmt.Println(database.Name)
}

func main() {
	CreateDB("test")
}
