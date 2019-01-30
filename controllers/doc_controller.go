package controllers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/mcornut/go-rest-api/repositories"
	"github.com/mcornut/go-rest-api/requests"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const uploadPath = "./uploads"

// DocumentController struct
type DocumentController interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type documentController struct {
	DB *sql.DB
}

// NewDocumentController func
func NewDocumentController(db *sql.DB) *documentController {
	return &documentController{
		DB: db,
	}
}

// Create func
func (doc *documentController) Create(w http.ResponseWriter, r *http.Request) {
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
	pdfPath, err := downloadPDFFile(params.Name, params.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repositories.CreateDocument(doc.DB, params.Name, pdfPath, "")
	if err != nil {
		log.Fatalf("Creating document: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// List func
func (doc *documentController) List(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != "GET" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	page := 1
	pageStr, ok := r.URL.Query()["page"]
	if ok {
		page, err = strconv.Atoi(pageStr[0])
		if err != nil {
			page = 1
		}
	}

	perPage := 10
	perPageStr, ok := r.URL.Query()["per_page"]
	if ok {
		perPage, err = strconv.Atoi(perPageStr[0])
		if err != nil {
			perPage = 1
		}
	}

	documents, err := repositories.GetDocuments(doc.DB, page, perPage)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(documents)
}

// downloadPDFFile func
func downloadPDFFile(filepath string, url string) (string, error) {

	pdfPath := fmt.Sprintf("%s/pdf/%s.pdf", uploadPath, filepath)
	tempPath := fmt.Sprintf("tmp/%s.pdf", filepath)

	if _, err := os.Stat(pdfPath); !os.IsNotExist(err) {
		return "", errors.New("File name already exists")
	}

	// Create the file
	out, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check content-type
	contentType := resp.Header.Get("Content-type")
	if contentType != "application/pdf" {
		return "", errors.New("Invalid file type")
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Get hash from new temporary file
	pdfHash, err := hashFileMD5(tempPath)
	if err != nil {
		log.Error(err)
		return "", errors.New("Internal server error")
	}

	// Check duplicate
	pdfDirectory := fmt.Sprintf("%s/pdf", uploadPath)
	existingFiles, err := getFileStringFromDirectory(pdfDirectory)
	if err != nil {
		log.Error(err)
		return "", errors.New("Internal server error")
	}

	for _, path := range existingFiles {
		currentFileHash, err := hashFileMD5(fmt.Sprintf("%s/%s", pdfDirectory, path))
		if err != nil {
			log.Error(err)
			return "", errors.New("Internal server error")
		}
		if pdfHash == currentFileHash {
			os.Remove(tempPath)
			return "", errors.New("Duplicate file")
		}
	}

	// Copy file to pdf directory (remove tempory file)
	err = cutAndPaste(tempPath, pdfPath)
	if err != nil {
		log.Error(err)
		os.Remove(tempPath)
		return "", errors.New("Internal server error")
	}

	// TODO
	err = convertPdfToJpg(pdfPath, "test.jpg")
	if err != nil {
		log.Error(err)
		os.Remove(pdfPath)
		return "", errors.New("Internal server error")
	}

	os.Remove(tempPath)

	return pdfPath, nil
}

func getFileStringFromDirectory(directory string) ([]string, error) {
	var filesString []string

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
		return filesString, err
	}

	for _, f := range files {
		filesString = append(filesString, f.Name())
	}

	return filesString, nil
}

func hashFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		log.Error(err)
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Error(err)
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func cutAndPaste(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		log.Error(err)
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		log.Error(err)
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Error(err)
		return err
	}

	return os.Remove(src)
}

func convertPdfToJpg(pdfName string, imageName string) error {

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(300, 300); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(pdfName); err != nil {
		return err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return err
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}

	// Select only first page of pdf
	mw.SetIteratorIndex(0)

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	// Resize image
	mw.ResizeImage(100, 150, imagick.FILTER_SINC, 1)

	// Save File
	return mw.WriteImage(imageName)
}
