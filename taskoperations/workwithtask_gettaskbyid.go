package taskoperations

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskOperations/gettaskbyid] ")
	log.SetOutput(os.Stdout)
}

func GetTaskById(db *sql.DB, id string) ([]byte, int, error) {
	var t Task

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Ошибка: задача не найдена")
			return []byte{}, http.StatusNotFound, errors.New(`{"error":"Задача не найдена"}`)
		}
		log.Printf("Ошибка при записи данных: %v", err)
		return []byte{}, http.StatusInternalServerError, errors.New(`{"error":"Ошибка записи данных"}`)
	}

	result, err := json.Marshal(t)
	if err != nil {
		log.Printf("Ошибка при маршализации задачи: %v", err)
		return nil, http.StatusInternalServerError, errors.New(`{"error":"Ошибка при маршализации задачи"}`)
	}

	return result, http.StatusOK, nil
}
