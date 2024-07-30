package router

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	handlers "go_final_project/handlers"
)

func NewRouter() *mux.Router {
	godotenv.Load(".env")
	//port := os.Getenv("TODO_PORT")
	r := mux.NewRouter()

	r.HandleFunc("/api/nextdate", handlers.NextDateHandler).Methods("GET")

	return r
}
