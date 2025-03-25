package models

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitSQLite() error {
	_, err := os.Open(path.Join("./databases", "database.sqlite3"))
	if err != nil {
		schema, err := os.Open(path.Join("./migrations", "schema.sqlite3"))
		if err != nil {
			return err
		}
		defer schema.Close()

		dst, err := os.Create(path.Join("./databases", "database.sqlite3"))
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, schema)
		if err != nil {
			return err
		}
	}

	return initDBPool()
}

func initDBPool() error {
	dbPool, err := sql.Open("sqlite3", `file:./databases/database.sqlite3?_foreign_keys=true&cache=private`)
	if err != nil {
		return fmt.Errorf("error initalizing database: %s", err)
	}

	pingErr := dbPool.Ping()
	if pingErr != nil {
		return fmt.Errorf("error connecting to database: %s", pingErr)
	}

	dbPool.SetConnMaxIdleTime(time.Minute * 60)
	dbPool.SetMaxOpenConns(10)
	dbPool.SetMaxIdleConns(5)

	db = dbPool

	return nil
}
