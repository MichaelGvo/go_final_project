package taskoperations

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"go_final_project/nextdate"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskDoneHandler] ")
	log.SetOutput(os.Stdout)
}

func TaskDoneHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		log.Println("Идентификатор не найден")
		http.Error(w, `{"error":"Идентификатор не найден"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		log.Printf("Ошибка при открытии базы данных: %v", err)
		http.Error(w, `{"error":"Ошибка при открытии базы данных: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var taskDone Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err = row.Scan(&taskDone.ID, &taskDone.Date, &taskDone.Title, &taskDone.Comment, &taskDone.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не удается обнаружить задачу")
			http.Error(w, `{"error":"Не удается обнаружить задачу"}`, http.StatusNotFound)
		} else {
			log.Printf("Ошибка при сканировании задачи: %v", err)
			http.Error(w, `{"error":"Ошибка при сканировании задачи: `+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	if taskDone.Repeat == "" {
		_, ResponseStatus, err := DeleteTask(db, id)
		if err != nil {
			log.Printf("Ошибка удаления задачи: %v", err)
			http.Error(w, `{"error":"Ошибка удаления задачи: `+err.Error()+`"}`, ResponseStatus)
			return
		}
	} else {
		now := time.Now()
		nextDateInStr, err := nextdate.NextDate(now, taskDone.Date, taskDone.Repeat)
		if err != nil {
			log.Printf("Ошибка при получении следующей даты: %v", err)
			http.Error(w, `{"error":"Ошибка при получении следующей даты: `+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		nextDate, err := time.Parse("20060102", nextDateInStr)
		if err != nil {
			log.Printf("Ошибка приведения даты к формату: %v", err)
			http.Error(w, `{"error":"Ошибка приведения даты к формату: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		taskDone.Date = nextDate.Format("20060102")
		taskJson, err := json.Marshal(taskDone)
		if err != nil {
			log.Printf("Ошибка преобразования задачи: %v", err)
			http.Error(w, `{"error":"Ошибка преобразования задачи: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		newRequest, err := http.NewRequest(http.MethodPut, "", bytes.NewBuffer(taskJson))
		if err != nil {
			log.Printf("Ошибка при создании нового запроса: %v", err)
			http.Error(w, `{"error":"Ошибка при создании нового запроса: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		newRequest.Header.Set("Content-Type", "application/json")

		q := newRequest.URL.Query()
		q.Add("id", id)
		newRequest.URL.RawQuery = q.Encode()

		_, ResponseStatus, err := UpdateTask(db, newRequest)
		if err != nil {
			log.Printf("Ошибка при обновлении задачи: %v", err)
			http.Error(w, `{"error":"Ошибка при обновлении задачи: `+err.Error()+`"}`, ResponseStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
