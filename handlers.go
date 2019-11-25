package main

import (
	"html/template"
	"io/ioutil"
	"fmt"
	"net/http"
)

// HomePageData stores the data of the landing page
type HomePageData struct {
    PageTitle string
}

type result struct {
	Title string
	Ts_headline template.HTML
}

func homePage(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		PageTitle: "Home",
	}

	err := HomeTemplate.Execute(w, data)
	if err != nil {
		logger.Print("HomeTemplate cant be loaded: ", err)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	logger.Print("Starting to upload the file")

	// max upload of 10 MB
	r.ParseMultipartForm(10 << 20)

	// retrieving the file from the request. this will give us the details of the given file
	file, handler, err := r.FormFile("inputFile")

	if err != nil {
		logger.Print("Error while retrieving the file ", err)
		fmt.Fprintf(w, "Error while retrieving the file: %v", err)
		return
	}

	defer file.Close()
	logger.Printf("Uploaded file: %+v\n", handler.Filename)
	logger.Printf("File size: %+v\n", handler.Size)
	logger.Printf("MIME header: %+v\n", handler.Header)

	// creating a temporary file
	tempFile, err := ioutil.TempFile("temp", "file-*.txt")

	if err != nil {
		logger.Println("tempfile", err)
	}
	defer tempFile.Close()

	// reading the contents of the upload file into a byte array
	fileByteFormat, err := ioutil.ReadAll(file)

	if err != nil {
		logger.Println("can't read the temp file", err)
	}

	_, err = DBCon.Exec(insertText, handler.Filename, fileByteFormat)

	if err != nil {
		logger.Println("database insertion error ", err)
	}

	tempFile.Write(fileByteFormat)
	fmt.Fprintf(w, string(fileByteFormat))
}

func searchFile(w http.ResponseWriter, r *http.Request) {
	logger.Println("Searching for keywords")

	if err := r.ParseForm(); err != nil {
		logger.Printf("ParseForm err: %v", err)
		return
	}

	keywords := r.FormValue("searchKeyword")
	
	rows, err := DBCon.Query(fullTextSearch, keywords)
	if err != nil {
		logger.Println("Keyword not found")
		http.Error(w, "keyword not found", http.StatusNotFound)
		return
	}
	defer rows.Close()

	results := make([]result, 0, 10)

	for rows.Next() {
		var r result
		if err := rows.Scan(&r.Title, &r.Ts_headline); err != nil {
			logger.Println("Keyword not found")
			http.Error(w, "keyword not found", http.StatusNotFound)
			return
		}

		results = append(results, r)

		if err := rows.Err(); err != nil {
			logger.Println("Keyword not found")
			http.Error(w, "keyword not found", http.StatusNotFound)
			return
		}
	}

	logger.Println("Keyword search is completed")

	if len(results) > 0 {
		err := SearchTemplate.Execute(w, results)
		if err != nil {
			logger.Print("SearchTemplate cant be loaded: ", err)
		}
		return
	}

	fmt.Fprintf(w, "nothing found")
}
