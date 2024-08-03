package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

// GetDBInstance возвращает единственный экземпляр соединения с базой данных
func GetDBInstance() *sql.DB {
	once.Do(func() {
		appPath, err := os.Executable()
		if err != nil {
			panic(err)
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

		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			panic(err)
		}

		if install {
			createTable := `CREATE TABLE IF NOT EXISTS scheduler (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                date    CHAR(8) NOT NULL DEFAULT "",
                title   VARCHAR(128) NOT NULL DEFAULT "",
                comment TEXT NOT NULL DEFAULT "",
                repeat TEXT NOT NULL DEFAULT ""
            );`
			_, err = db.Exec(createTable)
			if err != nil {
				panic(err)
			}
			createIndex := "CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date)"
			_, err = db.Exec(createIndex)
			if err != nil {
				panic(err)
			}
		}

		dbInstance = db
	})
	return dbInstance
}
