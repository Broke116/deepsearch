package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
)

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

	sqlStatement := `
		INSERT INTO public.textual (title, content)
		VALUES ($1, $2)
	`	
	_, err = DBCon.Exec(sqlStatement, handler.Filename, fileByteFormat)

	if err != nil {
		logger.Println("database insertion error ", err)
	}

	tempFile.Write(fileByteFormat)
	fmt.Fprintf(w, string(fileByteFormat))
}
