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

func OpenCloseDb() {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	var dbFile string
	envdbFile := os.Getenv("TODO_DBFILE")

	switch {
	case len(envdbFile) > 0:
		dbFile = envdbFile
	default:
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	Db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Db.Close()

	if install {
		createTable := `CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date    CHAR(8) NOT NULL DEFAULT "",
				title   VARCHAR(128) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat TEXT NOT NULL DEFAULT ""
			);`
		_, err = Db.Exec(createTable)
		if err != nil {
			fmt.Println(err)
			return
		}
		createIndex := "CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date)"
		_, err = Db.Exec(createIndex)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
