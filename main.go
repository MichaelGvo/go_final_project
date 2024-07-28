package main

import (
	//	"bytes"
	//	"encoding/json"

	"go_final_project/db"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	//	"github.com/go-chi/chi"
)

//type Task struct {
//	ID    string           `json:"id"`        // id задачи
//	Deadline  string       `json:"deadline"`  // дата дедлайна
//	Name  string           `json:"name"`      // заголовок
//	Comment string         `json:"comment"`   // комментарий
//	Rule string            `json:"rule"`      // правило, по которому задачи будут повторяться
//}

// var tasks = map[string]Task{}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("TODO_PORT")
	//r := router.Router()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	db.OpenCloseDb()

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

}
