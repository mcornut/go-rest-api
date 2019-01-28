package repositories

import (
	"database/sql"

	"github.com/mcornut/go-rest-api/models"
)

// CreateDocument func
func CreateDocument(db *sql.DB, name, filePath, thumbPath string) (int, error) {
	const query = `
		insert into documents (
			name,
			file_path,
			thumb_path
		) values (
			$1,
			$2,
			$3
		) returning id
	`
	var id int
	err := db.QueryRow(query, name, filePath, thumbPath).Scan(&id)
	return id, err
}

// GetDocuments func
func GetDocuments(db *sql.DB, page, resultsPerPage int) ([]*models.Document, error) {
	const query = `
		select
			id,
			name,
			file_path,
			thumb_path
		from
			documents
		limit $1 offset $2
	`
	documents := make([]*models.Document, 0)
	offset := (page - 1) * resultsPerPage

	rows, err := db.Query(query, resultsPerPage, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var document models.Document
		err = rows.Scan(&document.ID, &document.Name, &document.FilePath, &document.ThumbPath)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}
	return documents, err
}
