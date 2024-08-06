package taskoperations

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go_final_project/nextdate"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id"`      // id задачи
	Date    string `json:"date"`    // дата дедлайна
	Title   string `json:"title"`   // заголовок
	Comment string `json:"comment"` // комментарий
	Repeat  string `json:"repeat"`  // правило, по которому задачи будут повторяться
}

func Check(req *http.Request) (Task, int, error) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, 500, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, 500, err
	}

	if task.Title == "" {
		return task, 400, errors.New(`{"error":"task title is not specified"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		return task, 400, errors.New(`{"error":"incorrect date"}`)
	}
	var dateNew string
	if task.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
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
	fmt.Println(task.Date)
	return task, 200, nil

}

type Id struct {
	Id int64 `json:"id"`
}

func AddTask(db *sql.DB, req *http.Request) ([]byte, int, error) {
	var idResp Id

	task, ResponseStatus, err := Check(req)
	if err != nil {
		return []byte{}, ResponseStatus, err
	}

	result, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return []byte{}, 500, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return []byte{}, 500, err
	}

	idResp.Id = id

	idResult, err := json.Marshal(idResp)
	if err != nil {
		return []byte{}, 500, err
	}
	return idResult, 200, nil

}

var ResponseStatus int

// WorkWithTaskHandler возвращает обработчик для создания и обновления задач
func WorkWithTaskHandler(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("id")
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(w, "error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch req.Method {
	case http.MethodGet:
		if param == "" {
			http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err := GetTaskID(db, param)
		defer db.Close()
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case http.MethodPost:
		response, ResponseStatus, err := AddTask(db, req)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func GetTaskID(db *sql.DB, id string) ([]byte, int, error) {

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10")
	if err != nil {
		return []byte{}, ResponseStatus, fmt.Errorf(`{"error":"ошибка чтения: %v"}`, err)
	}
	defer rows.Close()
	var res []Task
	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				return []byte{}, ResponseStatus, fmt.Errorf(`{"error":"отсутствует задача"}`)
			}
			return []byte{}, ResponseStatus, fmt.Errorf(`{"error":"ошибка чтения: %v"}`, err)
		}
		res = append(res, task)
	}
	taskResult, err := json.Marshal(res)
	if err != nil {
		return []byte{}, 500, err
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return []byte{}, ResponseStatus, err
	}
	return taskResult, 200, nil
}
