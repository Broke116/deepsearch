package main

import (
	"html/template"
	"io/ioutil"
	"fmt"
	"net/http"
)

// PageData stores the data of the landing page
type PageData struct {
	PageTitle string
	Data []result
}

type result struct {
	FileName string
	Tsheadline template.HTML
	PageTitle string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	MainTemplate.ExecuteTemplate(w, "main", nil)
}

// uploadFile endpoint handles file upload operations, validations and error checks are handled before file upload
// the name of the file refers to title column on db, the content of the file refers to content column on db
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
		if err := rows.Scan(&r.FileName, &r.Tsheadline); err != nil {
			logger.Println("Keyword not found", err)
			http.Error(w, "keyword not found", http.StatusNotFound)
			return
		}

		results = append(results, r)

		if err := rows.Err(); err != nil {
			logger.Println("Keyword not found", err)
			http.Error(w, "keyword not found", http.StatusNotFound)
			return
		}
	}

	pageData := PageData{ PageTitle: "Search Result", Data: results }

	logger.Println("Keyword search is completed")

	if len(results) > 0 {
		MainTemplate.ExecuteTemplate(w, "search", pageData)
		return
	}

	fmt.Fprintf(w, "nothing found")
}
