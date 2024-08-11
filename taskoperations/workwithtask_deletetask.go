package taskoperations

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskOperations/deletetask] ")
	log.SetOutput(os.Stdout)
}

func DeleteTask(db *sql.DB, id string) ([]byte, int, error) {
	task, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		log.Printf("Ошибка при удалении задачи из базы данных: %v", err)
		return nil, 500, fmt.Errorf(`{"error":"%s"}`, err)
	}

	rowsAffected, err := task.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при получении количества удаленных строк: %v", err)
		return nil, 500, err
	}

	if rowsAffected == 0 {
		log.Println("Ошибка: не удается найти задачу")
		return nil, 400, errors.New(`{"error":"Не удается найти задачу"}`)
	}
	var empty Task
	response, err := json.Marshal(empty)
	if err != nil {
		log.Printf("Ошибка при маршализации пустого ответа: %v", err)
		return []byte{}, 400, err
	}

	return response, 200, nil
}
