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
	log.SetPrefix("[TaskOperations/updatetask] ")
	log.SetOutput(os.Stdout)
}

func UpdateTask(db *sql.DB, req *http.Request) ([]byte, int, error) {
	task, ResponseStatus, err := Check(req)
	if err != nil {
		log.Printf("Ошибка при проверке задачи: %v", err)
		return []byte{}, ResponseStatus, err
	}

	res, err := db.Exec(`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat
  WHERE id = :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		log.Printf("Ошибка при обновлении задачи в базе данных: %v", err)
		return []byte{}, 500, errors.New(`{"error":"Ошибка при обновлении задачи в базе данных"}`)
	}

	result, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при получении количества обновленных строк: %v", err)
		return []byte{}, 500, errors.New(`{"error":"Ошибка при получении количества обновленных строк"}`)
	}
	if result == 0 {
		log.Println("Ошибка: задача не найдена")
		return []byte{}, 400, errors.New(`{"error":"Задача не найдена"}`)
	}
	var empty Task
	response, err := json.Marshal(empty)
	if err != nil {
		log.Printf("Ошибка при маршализации пустого ответа: %v", err)
		return []byte{}, 400, err
	}

	return response, 200, nil
}
