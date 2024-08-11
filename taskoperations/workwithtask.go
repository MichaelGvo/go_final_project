package taskoperations

import (
	"database/sql"
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

func WorkWithTaskHandler(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("id")

	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		log.Printf("Ошибка при подключении к базе данных: %v", err)
		http.Error(w, "error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	switch req.Method {
	case http.MethodGet:
		if param == "" {
			log.Println("Ошибка: некорректный идентификатор")
			http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err := GetTaskById(db, param)
		defer db.Close()
		if err != nil {
			log.Printf("Ошибка при получении задачи: %v", err)
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case http.MethodPost:
		response, ResponseStatus, err := AddTask(db, req)
		if err != nil {
			log.Printf("Ошибка при добавлении задачи: %v", err)
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case http.MethodPut:
		response, ResponseStatus, err := UpdateTask(db, req)
		if err != nil {
			log.Printf("Ошибка при обновлении задачи: %v", err)
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case http.MethodDelete:
		if param == "" {
			log.Println("Ошибка: некорректный идентификатор")
			http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err := DeleteTask(db, param)
		if err != nil {
			log.Printf("Ошибка при удалении задачи: %v", err)
			http.Error(w, err.Error(), ResponseStatus)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	default:
		log.Println("Ошибка: метод не разрешен")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
