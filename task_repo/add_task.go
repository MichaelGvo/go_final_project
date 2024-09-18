package task_repo

import (
	"database/sql"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[task_repo/AddTask] ")
	log.SetOutput(os.Stdout)
}

type Id struct {
	Id int64 `json:"id"`
}

func (tr *TaskRepo) AddTask(task Task) (int64, error) {
	var id int64

	result, err := tr.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
  VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в базу данных: %v", err)
		return id, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка при получении идентификатора новой задачи: %v", err)
		return id, err
	}

	return id, nil
}
