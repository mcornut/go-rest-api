package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcornut/go-rest-api/models"
	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/**
* GET /documents
**/
func Test_ListDocumentsEmpty(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`{}`)
	req, err := http.NewRequest(http.MethodGet, "http://test.com/documents", &buf)
	if err != nil {
		t.Fatal(err)
	}

	// mocks
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("[select id, name, file_path, thumb_path from documents limit $1 offset $2]").WillReturnRows(sqlmock.NewRows([]string{}))

	rr := httptest.NewRecorder()
	documentController := NewDocumentController(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/documents", documentController.List)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %v, actual %v", http.StatusOK, rr.Code)
		return
	}
}

func Test_ListDocuments(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`{}`)
	req, err := http.NewRequest(http.MethodGet, "http://test.com/documents", &buf)
	if err != nil {
		t.Fatal(err)
	}

	// mocks
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"id", "name", "file_path", "thumb_path"}
	rs := sqlmock.NewRows(columns)
	rs.AddRow(1, "name1", "file_path", "thumb_path")
	rs.AddRow(2, "name2", "file_path2", "thumb_path2")
	mock.ExpectQuery("[select id, name, file_path, thumb_path from documents limit $1 offset $2]").WillReturnRows(rs)

	rr := httptest.NewRecorder()
	documentController := NewDocumentController(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/documents", documentController.List)
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %v, actual %v", http.StatusOK, rr.Code)
		return
	}
	d := json.NewDecoder(rr.Body)
	resp := []*models.Document{}
	assert.NoError(t, d.Decode(&resp))
	assert.Equal(t, 2, len(resp))
}

/**
* POST /documents
**/

// TODO
