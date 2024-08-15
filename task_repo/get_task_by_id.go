package task_repo

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[task_repo/GetTaskById] ")
	log.SetOutput(os.Stdout)
}

func (tr *TaskRepo) GetTaskByID(id string) (Task, error) {
	var t Task

	row := tr.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Ошибка: задача не найдена")
			return t, errors.New(`{"error":"Задача не найдена"}`)
		}
		log.Printf("Ошибка при записи данных: %v", err)
		return t, errors.New(`{"error":"Ошибка записи данных"}`)
	}

	return t, nil
}

//result, err := json.Marshal(t)
//if err != nil {
//	log.Printf("Ошибка при маршализации задачи: %v", err)
//	return nil, http.StatusInternalServerError, errors.New(`{"error":"Ошибка при маршализации задачи"}`)
//}
