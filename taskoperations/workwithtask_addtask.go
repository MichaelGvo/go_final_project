package taskoperations

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go_final_project/nextdate"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskOperations/addtask] ")
	log.SetOutput(os.Stdout)
}

type Id struct {
	Id int64 `json:"id"`
}

func (tr *TaskRepo) AddTask(req *http.Request) ([]byte, int, error) {
	var idResp Id

	task, ResponseStatus, err := Check(req)
	if err != nil {
		log.Printf("Ошибка при проверке задачи: %v", err)
		return []byte{}, ResponseStatus, err
	}

	result, err := tr.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
  VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи в базу данных: %v", err)
		return []byte{}, 500, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка при получении идентификатора новой задачи: %v", err)
		return []byte{}, 500, err
	}

	idResp.Id = id

	idResult, err := json.Marshal(idResp)
	if err != nil {
		log.Printf("Ошибка при маршализации идентификатора: %v", err)
		return []byte{}, 500, err
	}
	return idResult, 200, nil
}

func Check(req *http.Request) (Task, int, error) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Ошибка при чтении тела запроса: %v", err)
		return task, 500, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Printf("Ошибка при разборе JSON-данных: %v", err)
		return task, 500, err
	}

	if task.Title == "" {
		log.Println("Ошибка: не указано название задачи")
		return task, 400, errors.New(`{"error":"task title is not specified"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Ошибка при преобразовании даты: %v", err)
		return task, 400, errors.New(`{"error":"incorrect date"}`)
	}
	var dateNew string
	if task.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка при получении следующей даты: %v", err)
			return task, 400, err
		}
	}

	if task.Date == now.Format("20060102") {
		task.Date = now.Format("20060102")
	}

	if dateParse.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = dateNew
		}
	}

	return task, 200, nil
}
