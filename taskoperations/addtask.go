package taskoperations

import (
	"bytes"
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
	// используем тип bytes.Buffer для работы с байтовыми данными, в данном случае с содержимым тела запроса (r.Body)
	var buf bytes.Buffer
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// err = nil - чтение прошло успешно
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		resp := map[string]string{"error": "Чтение не прошло успешно"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	// err = nil - успешно десериализуем JSON из буфера в структуру Task
	// err != nil - вызываем ошибку (http.Error) с кодом состояния http.StatusBadRequest (должен вернуть статус 400 Bad Request)
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
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

	//_, err = time.Parse("20060102", task.Date)
	//if err != nil {
	//	resp := map[string]string{"error": "дата представлена в формате, отличном от 20060102"}
	//	w.WriteHeader(http.StatusBadRequest)
	//	json.NewEncoder(w).Encode(resp)
	//	return
	//}
	now := time.Now()
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

	//switch {
	//case task.Date == "":
	//	task.Date = nowForm
	//
	//	case task.Date == nowForm:
	//		task.Date = nowForm
	//
	//	case task.Date < nowForm && task.Repeat == "":
	//		task.Date = nowForm
	//	case task.Date < nowForm && task.Repeat != "":
	//		nextdate, err := finddate.NextDate(now, task.Date, task.Repeat)
	//		if err != nil {
	//			resp := map[string]string{"error": "правило повторения указано в неправильном формате"}
	//			w.WriteHeader(http.StatusBadRequest)
	//			json.NewEncoder(w).Encode(resp)
	//			return
	//		}
	//		task.Date = nextdate
	//	}

	id, err := AddTask(db, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		resp := map[string]string{"error": "Попытка добавить задачу не удалась"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		resp := map[string]string{"id": fmt.Sprintf("%d", id)}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}
