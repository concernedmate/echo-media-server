package models

import (
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/google/uuid"
)

type FileMetadata struct {
	FileId    string `json:"file_id"`
	Filename  string `json:"filename"`
	Filesize  int64  `json:"filesize"`
	Directory string `json:"directory"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

func UploadMultipleFiles(files []*multipart.FileHeader, directory string, username string) error {
	for _, file := range files {
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
	}
	return nil
}

func DeleteFile(file_id string) error {
	_, err := db.Exec(`DELETE FROM files WHERE file_id = ?`, file_id)
	if err != nil {
		return err
	}

	_ = os.Remove(path.Join("./uploads", file_id))
	return nil
}

func GetFileMetadata(file_id string) (FileMetadata, error) {
	var row FileMetadata
	err := db.QueryRow(
		`SELECT file_id, filename, directory, created_at, updated_at FROM files WHERE file_id = ?`,
		file_id,
	).Scan(&row.FileId, &row.Filename, &row.Directory, &row.CreatedAt, &row.UpdatedAt)
	if err != nil {
		return FileMetadata{}, err
	}

	return row, nil
}

func ListFiles(username string, basedir string) ([]FileMetadata, error) {
	rows, err := db.Query(
		`SELECT file_id, filename, directory, created_at, updated_at 
		FROM files 
		WHERE username = ? AND directory LIKE ?`,
		username, basedir+"%",
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

		file, err := os.Open(path.Join("./uploads", row.FileId))
		if err != nil {
			row.Filesize = -1
		} else {
			stat, err := file.Stat()
			if err != nil {
				row.Filesize = -1
			} else {
				row.Filesize = stat.Size()
			}
		}
		defer file.Close()

		files = append(files, row)
	}

	return files, nil
}
