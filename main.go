package main

import (
	"go_final_project/db"
	"go_final_project/handlers"
	"go_final_project/task_repo"
	"log"

	"net/http"
)

func main() {
	db, err := db.OpenCloseDb()
	if err != nil {
		log.Fatalf("Ошибка при создании базы: %v", err)
	}
	defer db.Close()
	tr := &task_repo.TaskRepo{DB: db}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", handlers.WorkWithTaskHandler(tr))
	mux.HandleFunc("/api/tasks", handlers.WorkWithTasksHandler(tr))
	mux.HandleFunc("/api/task/done", handlers.TaskDoneHandler(tr))

	err = http.ListenAndServe(":7540", mux)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return
	}
}
