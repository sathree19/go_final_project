package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(addr string) (*sql.DB, error) {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	dbFile := filepath.Join(appPath, addr)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	fmt.Println(dbFile)

	if install {
		_, err = os.Create(dbFile)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL DEFAULT """", 
		title VARCHAR(256) NOT NULL DEFAULT "", comment TEXT NOT NULL DEFAULT "", repeat VARCHAR(128) NOT NULL DEFAULT "")`)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		_, err = db.Exec(`CREATE INDEX task_date ON scheduler (date)`)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return db, nil

	}
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil

}
