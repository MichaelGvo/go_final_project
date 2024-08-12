package taskoperations

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[TaskOperations/workwithtasks] ")
	log.SetOutput(os.Stdout)
}

func WorkWithTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response []byte
		var responseStatus int
		var err error

		if req.Method != http.MethodGet {
			log.Println("Ошибка: метод не разрешен")
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		response, responseStatus, err = GetTasks(db)

		if err != nil {
			log.Printf("Ошибка при получении задач: %v", err)
			http.Error(w, err.Error(), responseStatus)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func GetTasks(db *sql.DB) ([]byte, int, error) {
	var tasks []Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Printf("Ошибка при чтении данных из базы данных: %v", err)
			return nil, http.StatusInternalServerError, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Ошибка при чтении данных из базы данных: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	if tasks == nil {
		tasks = []Task{}
	}

	response, err := json.Marshal(map[string][]Task{"tasks": tasks})
	if err != nil {
		log.Printf("Ошибка при маршализации ответа: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusOK, nil
}
