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
var (
	// DBCon is a database connection name
	DBCon *sql.DB
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
		
	DBCon, err := sql.Open("postgres", postgres)
	if err != nil {
		logger.Print("database connection error: ", err)
	}
	defer DBCon.Close()

	err = DBCon.Ping()
	if err != nil {
		logger.Print("database ping error: ", err)
	}

	logger.Println("Database connection initialized")
}

func main() {
	initDB()
	createRoutes()

	logger.Println("Listening on port", serverAddress)
	logger.Fatal(http.ListenAndServe(serverAddress, nil))
}
