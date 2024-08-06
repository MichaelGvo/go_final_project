package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func OpenCloseDb() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	Db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if install {
		createTable := `CREATE TABLE IF NOT EXISTS scheduler (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
			date    CHAR(8) NOT NULL DEFAULT "",
			title   VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat  VARCHAR(128) NOT NULL DEFAULT ""
		);`
		_, err = Db.Exec(createTable)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		createIndex := "CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date)"
		_, err = Db.Exec(createIndex)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return Db, nil
}
