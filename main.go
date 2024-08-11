package main

import (
	"go_final_project/db"
	handlers "go_final_project/handlers"
	"go_final_project/taskoperations"
	"log"

	"net/http"
)

func main() {
	db, err := db.OpenCloseDb()
	if err != nil {
		log.Fatalf("Ошибка при создании базы: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("web/")))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", taskoperations.WorkWithTaskHandler)
	mux.HandleFunc("/api/tasks", taskoperations.WorkWithTasksHandler(db))
	mux.HandleFunc("/api/task?id=", taskoperations.WorkWithTaskHandler)
	mux.HandleFunc("/api/task/done?id=", taskoperations.TaskDoneHandler)

	err = http.ListenAndServe(":7540", mux)
	if err != nil {
		log.Printf("Error occurred: %v", err)
		panic(err)
	}
}
