package main

import (
	"html/template"
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
	// HomeTemplate stores the html file for home page
	HomeTemplate *template.Template
	// SearchTemplate show the html structure of search page.
	SearchTemplate *template.Template
)

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

func prepareTemplates() {
	HomeTemplate = template.Must(template.ParseFiles("static/index.html"))
	SearchTemplate = template.Must(template.ParseFiles("static/search.html"))
}

func createRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/search", searchFile)
}

func main() {
	initDB()
	prepareTemplates()
	createRoutes()

	defer DBCon.Close()

	logger.Println("Listening on port", serverAddress)
	logger.Fatal(http.ListenAndServe(serverAddress, nil))
}
