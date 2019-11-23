package main

import (
	"database/sql"
	"fmt"
	"os"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	serverAddress = ":8080"
	host     = "localhost"
	port     = 5433
	user     = "admin"
	password = "admin"
	dbname   = "deepsearch"
)

var logger = log.New(os.Stdout, "http: ", log.LstdFlags)

func createRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadFile)
}

func initDB() {
	postgres := fmt.Sprintf("host=%s port=%d user=%s "+
    	"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
		
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		logger.Print("database connection error: ", err)
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Print("database ping error: ", err)
		panic(err)	
	}

	logger.Println("Database connection initialized")
}

func main() {
	initDB()
	createRoutes()

	logger.Println("Listening on port", serverAddress)
	logger.Fatal(http.ListenAndServe(serverAddress, nil))
}
