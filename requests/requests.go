package requests

import (
	"fmt"
)

// CreateDocumentRequest struct
type CreateDocumentRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Validate func
func (r CreateDocumentRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	if r.URL == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}
