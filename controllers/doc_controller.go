package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mcornut/go-rest-api/requests"
)

const uploadPath = "./uploads"

// DocumentController struct
type DocumentController struct {
	DB *sql.DB
}

// NewDocumentController func
func NewDocumentController(db *sql.DB) *DocumentController {
	return &DocumentController{
		DB: db,
	}
}

// Create func
func (doc *DocumentController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}

	// Unmarshal
	var params requests.CreateDocumentRequest
	err = json.Unmarshal(b, &params)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = params.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// downloadPDFFile
	err = downloadPDFFile(params.Name, params.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List func
func (doc *DocumentController) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// downloadPDFFile func
func downloadPDFFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/pdf/%s.pdf", uploadPath, filepath))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check content-type
	contentType := resp.Header.Get("Content-type")
	if contentType != "application/pdf" {
		return errors.New("Invalid file type")
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
