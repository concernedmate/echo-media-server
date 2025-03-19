package models

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
)

type FileMetadata struct {
	FileId    string    `json:"file_id"`
	Filename  string    `json:"filename"`
	Directory string    `json:"directory"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UploadFile(file *multipart.FileHeader, directory string, username string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	file_id := uuid.New()
	dst, err := os.Create(path.Join("./uploads", file_id.String()))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT INTO files (file_id, filename, directory, username) VALUES (?, ?, ?, ?)`,
		file_id, file.Filename, directory, username,
	)
	if err != nil {
		_ = os.Remove(path.Join("./uploads", file_id.String()))
		return err
	}

	return nil
}

func GetFileMetadata(file_id string) (FileMetadata, error) {
	var row FileMetadata
	err := db.QueryRow(
		`SELECT file_id, filename, directory, created_at, updated_at WHERE file_id = ?`,
		file_id,
	).Scan(&row.FileId, &row.Filename, &row.Directory, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		return FileMetadata{}, err
	}

	return row, nil
}

func ListFiles(basedir string) ([]FileMetadata, error) {
	rows, err := db.Query(
		`SELECT file_id, filename, directory, created_at, updated_at WHERE directory LIKE ?`,
		basedir+"%",
	)
	if err != nil {
		return []FileMetadata{}, err
	}
	defer rows.Close()

	files := []FileMetadata{}
	for rows.Next() {
		var row FileMetadata
		err = rows.Scan(&row.FileId, &row.Filename, &row.Directory, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return []FileMetadata{}, err
		}

		files = append(files, row)
	}

	return files, nil
}
