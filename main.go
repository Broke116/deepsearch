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

var (
    // DBCon is the connection handle
    // for the database
	DBCon *sql.DB
	logger = log.New(os.Stdout, "http: ", log.LstdFlags)
	err error
)

func createRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/upload", uploadFile)
}

func initDB() {
	postgres := fmt.Sprintf("host=%s port=%d user=%s "+
    	"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
		
	DBCon, err = sql.Open("postgres", postgres)
	if err != nil {
		logger.Print("database connection error: ", err)
	}	

	err = DBCon.Ping()
	if err != nil {
		logger.Print("database ping error: ", err)
	}

	logger.Println("Database connection initialized")
}

func main() {
	initDB()
	createRoutes()

	defer DBCon.Close()

	logger.Println("Listening on port", serverAddress)
	logger.Fatal(http.ListenAndServe(serverAddress, nil))
}
