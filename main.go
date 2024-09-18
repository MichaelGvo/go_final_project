package main

import (
	"log"
	"net/http"

	"go_final_project/db"
	"go_final_project/handlers"
	"go_final_project/task_repo"
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
	mux.HandleFunc("/api/nextdate", handlers.Next_Date)
	mux.HandleFunc("/api/task", handlers.Get_Task(tr))
	mux.HandleFunc("/api/tasks", handlers.Get_Tasks(tr))
	mux.HandleFunc("/api/task/done", handlers.Task_Done(tr))

	err = http.ListenAndServe(":7540", mux)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		return
	}
}
