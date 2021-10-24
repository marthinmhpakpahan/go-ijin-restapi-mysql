package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"os"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var _ = godotenv.Load(".env")
var (
	ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("user"),
		os.Getenv("pass"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("db_name"),
	)
)

func getDB() (*sql.DB, error) {
	return sql.Open("mysql", ConnectionString)
}