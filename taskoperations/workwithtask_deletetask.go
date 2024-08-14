package taskoperations

import (
	"database/sql"
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

func (tr *TaskRepo) DeleteTask(id string) ([]byte, int, error) {
	task, err := tr.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
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

	return []byte("{}"), 200, nil
}
