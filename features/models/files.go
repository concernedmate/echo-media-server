package models

import (
	"io"
	"media-server/configs"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

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

type DirectoryMetadata struct {
	Dirname   string `json:"dirname"`
	Directory string `json:"directory"`
}

func UploadFile(file *multipart.FileHeader, directory string, username string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	file_id := uuid.New()
	dst, err := os.Create(path.Join(configs.UPLOAD_BASEDIR(), file_id.String()))
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
		_ = os.Remove(path.Join(configs.UPLOAD_BASEDIR(), file_id.String()))
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
		dst, err := os.Create(path.Join(configs.UPLOAD_BASEDIR(), file_id.String()))
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
			_ = os.Remove(path.Join(configs.UPLOAD_BASEDIR(), file_id.String()))
			return err
		}
	}
	return nil
}

func SaveFileWebsocket(file_id string, filename string, directory string, username string) error {
	_, err := db.Exec(
		`INSERT INTO files (file_id, filename, directory, username) VALUES (?, ?, ?, ?)`,
		file_id, filename, directory, username,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(file_id string) error {
	_, err := db.Exec(`DELETE FROM files WHERE file_id = ?`, file_id)
	if err != nil {
		return err
	}

	_ = os.Remove(path.Join(configs.UPLOAD_BASEDIR(), file_id))
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
		WHERE username = ? AND directory = ?`,
		username, basedir,
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

		file, err := os.Open(path.Join(configs.UPLOAD_BASEDIR(), row.FileId))
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

func ListDirectory(username string, basedir string) ([]DirectoryMetadata, error) {
	rows, err := db.Query(
		`SELECT DISTINCT directory 
		FROM files 
		WHERE username = ? AND directory != ? AND directory LIKE ?`,
		username, basedir, basedir+"%",
	)
	if err != nil {
		return []DirectoryMetadata{}, err
	}
	defer rows.Close()

	dirs := []DirectoryMetadata{}
	cache := map[string]bool{}
	for rows.Next() {
		var row DirectoryMetadata
		err = rows.Scan(&row.Directory)
		if err != nil {
			return []DirectoryMetadata{}, err
		}

		if basedir == "/" {
			split := strings.Split(row.Directory, "/")
			if len(split) > 1 {
				row.Dirname = split[1]
				row.Directory = "/" + split[1]
				if !cache[row.Dirname] {
					dirs = append(dirs, row)
					cache[row.Dirname] = true
				}
			}
		} else {
			split1 := strings.Split(row.Directory, basedir)[1]
			split2 := strings.Split(split1, "/")
			if len(split2) > 1 {
				row.Dirname = split2[1]
				row.Directory = path.Join(basedir, split2[1])
				if !cache[row.Dirname] {
					dirs = append(dirs, row)
					cache[row.Dirname] = true
				}
			}
		}
	}

	return dirs, nil
}

func GetTotalSize(username string, basedir string) (stored string, max_storage string, err error) {
	rows, err := db.Query(
		`SELECT file_id FROM files WHERE username = ? AND directory = ?`,
		username, basedir,
	)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var total int64 = 0
	for rows.Next() {
		var file_id string
		err = rows.Scan(&file_id)
		if err != nil {
			return "", "", err
		}

		file, err := os.Open(path.Join(configs.UPLOAD_BASEDIR(), file_id))
		if err != nil {
			continue
		}
		defer func() {
			_ = file.Close()
		}()

		stat, err := file.Stat()
		if err != nil {
			continue
		}

		total += stat.Size()

		_ = file.Close()
	}

	if total < 1000*1000 {
		stored = strconv.Itoa(int(total/1000)) + " KB"
	} else if total < 1000*1000*1000 {
		stored = strconv.Itoa(int(total/1000/1000)) + " MB"
	} else {
		stored = strconv.Itoa(int(total/1000/1000/1000)) + " GB"
	}

	var max_stored_int int64
	err = db.QueryRow(`SELECT max_storage FROM users WHERE username = ?`, username).Scan(&max_stored_int)
	if err != nil {
		return "", "", err
	}

	if max_stored_int < 1000*1000 {
		max_storage = strconv.Itoa(int(max_stored_int/1000)) + " KB"
	} else if max_stored_int < 1000*1000*1000 {
		max_storage = strconv.Itoa(int(max_stored_int/1000/1000)) + " MB"
	} else {
		max_storage = strconv.Itoa(int(max_stored_int/1000/1000/1000)) + " GB"
	}
	if max_stored_int == -1 {
		max_storage = "unlimited"
	}

	return stored, max_storage, nil
}
