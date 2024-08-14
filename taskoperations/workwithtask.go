package taskoperations

import (
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskOperations/workwithtask] ")
	log.SetOutput(os.Stdout)
}

var ResponseStatus int

func WorkWithTaskHandler(db *TaskRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		param := req.URL.Query().Get("id")
		switch req.Method {
		case http.MethodGet:
			if param == "" {
				log.Println("Ошибка: некорректный идентификатор")
				http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
				return
			}
			response, ResponseStatus, err := db.GetTaskById(param)

			if err != nil {
				log.Printf("Ошибка при получении задачи: %v", err)
				http.Error(w, err.Error(), ResponseStatus)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		case http.MethodPost:
			response, ResponseStatus, err := db.AddTask(req)
			if err != nil {
				log.Printf("Ошибка при добавлении задачи: %v", err)
				http.Error(w, err.Error(), ResponseStatus)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		case http.MethodPut:
			response, ResponseStatus, err := db.UpdateTask(req)
			if err != nil {
				log.Printf("Ошибка при обновлении задачи: %v", err)
				http.Error(w, err.Error(), ResponseStatus)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		case http.MethodDelete:
			if param == "" {
				log.Println("Ошибка: некорректный идентификатор")
				http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
				return
			}
			response, ResponseStatus, err := db.DeleteTask(param)
			if err != nil {
				log.Printf("Ошибка при удалении задачи: %v", err)
				http.Error(w, err.Error(), ResponseStatus)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		default:
			log.Println("Ошибка: метод не разрешен")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
