package task_repo

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[task_repo/UpdateTask] ")
	log.SetOutput(os.Stdout)
}

func (tr *TaskRepo) UpdateTask(task Task) error {

	res, err := tr.DB.Exec(`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat
  WHERE id = :id`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		log.Printf("Ошибка при обновлении задачи в базе данных: %v", err)
		return errors.New(`{"error":"Ошибка при обновлении задачи в базе данных"}`)
	}

	result, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при получении количества обновленных строк: %v", err)
		return errors.New(`{"error":"Ошибка при получении количества обновленных строк"}`)
	}
	if result == 0 {
		log.Println("Ошибка: задача не найдена")
		return errors.New(`{"error":"Задача не найдена"}`)
	}

	return nil
}
