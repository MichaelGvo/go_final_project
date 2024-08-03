package main

import (
	//	"bytes"
	//	"encoding/json"

	"go_final_project/db"
	handlers "go_final_project/handlers"
	"go_final_project/taskoperations"
	"log"

	"net/http"
)

//type Task struct {
//	ID    string           `json:"id"`        // id задачи
//	Deadline  string       `json:"deadline"`  // дата дедлайна
//	Name  string           `json:"name"`      // заголовок
//	Comment string         `json:"comment"`   // комментарий
//	Rule string            `json:"rule"`      // правило, по которому задачи будут повторяться
//}

//var tasks = map[string]Task{}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", taskoperations.PostTaskHandler)
	db.OpenCloseDb()
	err := http.ListenAndServe(":7540", mux)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		panic(err)
	}
}

//godotenv.Load(".env")
//port := ":7540"
//mux := http.NewServeMux()

//mux.HandleFunc("/api/task", auth(handlers.TaskHandle))
//mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
//mux.Handle("/", http.FileServer(http.Dir("web/")))
//http.Handle("/", http.FileServer(http.Dir("./web")))

//db.OpenCloseDb()

//err := http.ListenAndServe(":7540", mux)
//if err != nil {
//	panic(err)
////}

//}

//mux = http.NewServeMux()

//mux.HandleFunc("/api/nextdate", handlers.NextDateHandle)
//mux.HandleFunc("/api/task", auth(handlers.TaskHandle))
//...
//mux.Handle("/", http.FileServer(http.Dir("web/")))

////strPort := defStrPort()
//err := http.ListenAndServe(strPort, mux)
//if err != nil {
//	panic(err)
//}
