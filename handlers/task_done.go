package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go_final_project/nextdate"
	"go_final_project/task_repo"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[Handlers/TaskDoneHandler] ")
	log.SetOutput(os.Stdout)
}

var ResponseStatus int

func Task_Done(db *task_repo.TaskRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")
		if id == "" {
			log.Println("Идентификатор не найден")
			http.Error(w, "{\"error\":\"Идентификатор не найден\"}", http.StatusBadRequest)
			return
		}

		taskDone, err := db.TaskDone(id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				log.Println("Не удается обнаружить задачу")
				json.NewEncoder(w).Encode(map[string]string{"error": "Не удается обнаружить задачу"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			log.Printf("Ошибка при сканировании задачи: %v", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при сканировании задачи"})
			return
		}

		if taskDone.Repeat == "" {
			err := db.DeleteTask(id)
			if err != nil {
				log.Printf("Ошибка удаления задачи: %v", err)
				json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при сканировании задачи"})
				http.Error(w, "{\"error\":\"Ошибка удаления задачи: "+err.Error()+"\"}", ResponseStatus)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})

			return
		} else {
			now := time.Now()
			nextDateInStr, err := nextdate.Next_Date(now, taskDone.Date, taskDone.Repeat)
			if err != nil {
				log.Printf("Ошибка при получении следующей даты: %v", err)
				http.Error(w, "{\"error\":\"Ошибка при получении следующей даты: "+err.Error()+"\"}", http.StatusBadRequest)
				return
			}

			nextDate, err := time.Parse("20060102", nextDateInStr)
			if err != nil {
				log.Printf("Ошибка приведения даты к формату: %v", err)
				http.Error(w, "{\"error\":\"Ошибка приведения даты к формату: "+err.Error()+"\"}", http.StatusInternalServerError)
				return
			}

			taskDone.Date = nextDate.Format("20060102")

			err = db.UpdateTask(taskDone)
			if err != nil {
				log.Printf("Ошибка при обновлении задачи: %v", err)
				http.Error(w, "{\"error\":\"Ошибка при обновлении задачи: "+err.Error()+"\"}", ResponseStatus)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})

		}
	}
}
