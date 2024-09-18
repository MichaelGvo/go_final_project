package task_repo

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[task_repo/TaskDone] ")
	log.SetOutput(os.Stdout)
}

func (tr *TaskRepo) TaskDone(id string) (Task, error) {
	var taskDone Task
	row := tr.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&taskDone.ID, &taskDone.Date, &taskDone.Title, &taskDone.Comment, &taskDone.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не удается обнаружить задачу")
			return taskDone, errors.New(`{"error":"not find the task"}`)
		} else {
			log.Printf("Ошибка при сканировании задачи: %v", err)
			return taskDone, err
		}
	}

	if err := row.Err(); err != nil {
		return taskDone, err
	}

	return taskDone, nil
}
