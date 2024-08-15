package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/task_repo"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[Handlers/WorkWithTaskHandler] ")
	log.SetOutput(os.Stdout)
}

var empty = task_repo.Task{
	ID:      "",
	Date:    "",
	Title:   "",
	Comment: "",
	Repeat:  "",
}

type IDTask struct {
	Id int64 `json:"id"`
}

func WorkWithTaskHandler(db *task_repo.TaskRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		param := req.URL.Query().Get("id")
		var response []byte
		var err error
		var RespStatus int

		switch req.Method {
		case http.MethodGet:
			if param == "" {
				log.Println("Ошибка: некорректный идентификатор")
				http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
				return
			}
			task, err := db.GetTaskByID(param)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("Ошибка при получении задачи: %v", err)
					http.Error(w, `{"error":"task not found"}`, ResponseStatus)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response, err = json.Marshal(task)
			if err != nil {
				log.Printf("Ошибка при маршализации задачи: %v", err)
				http.Error(w, `{"error":"Ошибка при маршализации задачи"}`, http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodPost:
			var idtask IDTask
			task, RespStatus, err := task_repo.Check(req)
			if err != nil {
				log.Printf("Ошибка при проверке задачи: %v", err)
				http.Error(w, err.Error(), RespStatus)
				return
			}

			idtask.Id, err = db.AddTask(task)
			if err != nil {
				log.Printf("Ошибка при добавлении задачи: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response, err = json.Marshal(idtask)
			if err != nil {
				log.Printf("Ошибка при маршализации идентификатора: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodPut:
			task, RespStatus, err := task_repo.Check(req)
			if err != nil {
				log.Printf("Ошибка при проверке задачи: %v", err)
				http.Error(w, err.Error(), RespStatus)
				return
			}

			err = db.UpdateTask(task)
			if err != nil {
				if err.Error() == `{"error":"Задача не найдена"}` {
					log.Printf("Ошибка при обновлении задачи - задача не найдена: %v", err)
					http.Error(w, err.Error(), http.StatusNotFound)
				}
				log.Printf("Ошибка при обновлении задачи: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response, err = json.Marshal(empty)
			if err != nil {
				log.Printf("Ошибка при маршализации пустого ответа: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodDelete:
			if param == "" {
				log.Println("Ошибка: некорректный идентификатор")
				http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
				return
			}
			err := db.DeleteTask(param)
			if err != nil {
				if err.Error() == `{"error":"not found the task"}` {
					log.Println("Ошибка: не удается найти задачу - методы delete")
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				log.Printf("Ошибка при удалении задачи: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accept := map[string]interface{}{}
			response, err = json.Marshal(accept)
			if err != nil {
				log.Printf("Ошибка при маршализации пустого ответа: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		default:
			log.Println("Ошибка: метод не разрешен")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(RespStatus)
		_, err = w.Write(response)
		if err != nil {
			log.Printf("Ошибка при записи ответа в http.MethodGet: %v", err)
		}
	}
}
