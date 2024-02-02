package controller

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	//"ftms/models"
	//"strconv"
	//"time"

	//"net/http"

	//"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
)

// var host = "http://localhost:5000"

var db *sql.DB

// connecting to mysql database
func CreateDbConnection() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/samurai")
	// port 3306 is the default port for mysql in xampp
	// here ftms is the database name

	if err != nil {
		fmt.Println("Error connecting databse!")
		panic(err.Error())
	}

	// Ping the database to ensure the connection is valid.
	if err := db.Ping(); err != nil {
		fmt.Printf("Could not connect to the database: %v", err)
	}

	//defer db.Close()
	fmt.Println("Successfully connected to mysql database")
}