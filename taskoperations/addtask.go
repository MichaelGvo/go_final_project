package taskoperations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/db"
	finddate "go_final_project/nextdate"
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

func AddTask(db *sql.DB, date, title, comment, repeat string) (int64, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))
	if err != nil {
		return 0, err
	}
	// возвращаем идентификатор последней добавленной записи через LastInsertId
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	db := db.GetDBInstance()

	var task Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		resp := map[string]string{"error": "ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if task.Title == "" {
		resp := map[string]string{"error": "Не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	now := time.Now()
	//now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 1, time.Local)
	nowForm := now.Format("20060102")

	if task.Date == "" {
		task.Date = nowForm
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			resp := map[string]string{"error": "Дата указана в неверном формате"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		if task.Date == nowForm {
			task.Date = nowForm
		} else if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = nowForm
			} else {
				nextDate, err := finddate.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					resp := map[string]string{"error": err.Error()}
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(resp)
					return
				}
				task.Date = nextDate
			}
		}
	}
	//fmt.Print(task.Date)
	//fmt.Print(task.Title)
	//fmt.Print(task.Comment)
	//fmt.Print(task.Repeat)
	id, err := AddTask(db, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		resp := map[string]string{"error": "Попытка добавить задачу не удалась"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	fmt.Print(task.Date)
	fmt.Print(task.Title)
	fmt.Print(task.Comment)
	fmt.Print(task.Repeat)

	fmt.Print(task.Date > now.Format("20060102"))
	fmt.Print(task.Date == now.Format("20060102"))
	fmt.Print(task.Date < now.Format("20060102"))
}
