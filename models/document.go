package models

// Document struct
type Document struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FilePath  string `json:"file_path"`
	ThumbPath string `json:"thumb_path"`
}
