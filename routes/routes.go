package routes

import (
	"net/http"

	"github.com/mcornut/go-rest-api/controllers"
)

// CreateRoutes func
func CreateRoutes(mux *http.ServeMux, doc *controllers.DocumentController) {
	mux.HandleFunc("/document", doc.Create)
	mux.HandleFunc("/documents", doc.List)
}
