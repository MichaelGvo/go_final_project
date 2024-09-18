package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go_final_project/task_repo"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[Handlers/WorkWithTasksHandler] ")
	log.SetOutput(os.Stdout)
}

func Get_Tasks(tr *task_repo.TaskRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response []byte
		var responseStatus int

		if req.Method != http.MethodGet {
			log.Println("Ошибка: метод не разрешен")
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
		id := req.URL.Query().Get("id")
		if id != "" {
			taskSome, err := tr.GetTaskByID(id)
			if err != nil {
				if err.Error() == "{\"error\":\"Задача не найдена\"}" {
					log.Println("Ошибка: задача не найдена. хэндлер")
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response, err = json.Marshal(taskSome)
			if err != nil {
				log.Printf("Ошибка при сериализации задачи: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			responseStatus = http.StatusOK
		} else {
			tasks, err := tr.GetTasks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response, err = json.Marshal(map[string][]task_repo.Task{"tasks": tasks})
			if err != nil {
				log.Printf("Ошибка при сериализации списка задач: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			responseStatus = http.StatusOK
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseStatus)
		_, err := w.Write(response)
		if err != nil {
			log.Printf("Ошибка при записи ответа: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
