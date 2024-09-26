package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Name string
	Conn *pgx.Conn
}

func NewDatabase(user, password, host, port, name string, conn *pgx.Conn) Database {
	return Database{name, conn}
}

func Create(user, password, host, port, name string, conn *pgx.Conn) *Database {
	db := NewDatabase(user, password, host, port, name, conn)
	db.createSchema()

	return &db
}

func (db Database) createSchema() {
	databaseInit := []string{
		"DROP DATABASE IF EXISTS " + db.Name,
		"CREATE DATABASE " + db.Name,
		"USE " + db.Name,
		"CREATE TABLE products (id SERIAL PRIMARY KEY, name STRING, price STRING, quantity STRING)",
	}

	db.ExecSQL(databaseInit)

}

func (db Database) ExecSQL(sql []string) {
	for _, stmt := range sql {
		if _, err := db.Conn.Exec(context.Background(), stmt); err != nil {
			log.Fatal(err)
		}
	}
}

func (db Database) Query(sql string, values ...any) (pgx.Rows, error) {
	rows, err := db.Conn.Query(context.Background(), sql, values...)
	return rows, err
}

func (db Database) ExecQuery(sql string, values ...any) error {
	_, err := db.Conn.Exec(context.Background(), sql, values...)
	return err
}

func (db Database) QueryRow(sql string, values ...any) pgx.Row {
	row := db.Conn.QueryRow(context.Background(), sql, values...)
	return row
}

func GetConnection(user, host, port, name string) (*pgx.Conn, error) {
	connectionString := fmt.Sprintf("postgresql://%s:@%s:%s/%s", user, host, port, name)
	conn, err := pgx.Connect(context.Background(), connectionString)

	return conn, err
}

func USE(user, password, host, port, name string, conn *pgx.Conn) *Database {
	db := NewDatabase(user, password, host, port, name, conn)
	db.ExecQuery("USE " + name)

	return &db
}
