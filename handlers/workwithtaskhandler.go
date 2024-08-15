package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go_final_project/nextdate"
	"go_final_project/task_repo"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[Handlers/WorkWithTaskHandler] ")
	log.SetOutput(os.Stdout)
}

type IDTaskResponse struct {
	ID int64 `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WorkWithTaskHandler(db *task_repo.TaskRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handleGetTask(w, req, db)
		case http.MethodPost:
			handleCreateTask(w, req, db)
		case http.MethodPut:
			handleUpdateTask(w, req, db)
		case http.MethodDelete:
			handleDeleteTask(w, req, db)
		default:
			sendErrorResponse(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetTask(w http.ResponseWriter, req *http.Request, tr *task_repo.TaskRepo) {
	id := req.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}
	task, err := tr.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Задача не найдена", http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, task, http.StatusOK)
}

func handleCreateTask(w http.ResponseWriter, req *http.Request, tr *task_repo.TaskRepo) {
	task, respStatus, err := Check(req)
	if err != nil {
		sendErrorResponse(w, err.Error(), respStatus)
		return
	}
	id, err := tr.AddTask(task)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, IDTaskResponse{ID: id}, http.StatusCreated)
}

func handleUpdateTask(w http.ResponseWriter, req *http.Request, tr *task_repo.TaskRepo) {
	task, respStatus, err := Check(req)
	if err != nil {
		sendErrorResponse(w, err.Error(), respStatus)
		return
	}
	if err := tr.UpdateTask(task); err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Задача не найдена", http.StatusNotFound)
		} else {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	sendJSONResponse(w, struct{}{}, http.StatusOK)
}

func handleDeleteTask(w http.ResponseWriter, req *http.Request, tr *task_repo.TaskRepo) {
	id := req.URL.Query().Get("id")
	if id == "" {
		sendErrorResponse(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}
	if err := tr.DeleteTask(id); err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Задача не найдена", http.StatusNotFound)
		} else {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	sendJSONResponse(w, struct{}{}, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	sendJSONResponse(w, ErrorResponse{Error: errorMsg}, statusCode)
}

func Check(req *http.Request) (task_repo.Task, int, error) {
	var task task_repo.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		log.Printf("Ошибка при чтении тела запроса: %v", err)
		return task, http.StatusInternalServerError, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Printf("Ошибка при разборе JSON-данных: %v", err)
		return task, http.StatusInternalServerError, err
	}

	if task.Title == "" {
		log.Println("Ошибка: не указано название задачи")
		return task, http.StatusBadRequest, errors.New(`{"error":"task title is not specified"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Printf("Ошибка при преобразовании даты: %v", err)
		return task, http.StatusBadRequest, errors.New(`{"error":"incorrect date"}`)
	}
	var dateNew string
	if task.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка при получении следующей даты: %v", err)
			return task, http.StatusBadRequest, err
		}
	}

	if task.Date == now.Format("20060102") {
		task.Date = now.Format("20060102")
	}

	if dateParse.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = dateNew
		}
	}

	return task, http.StatusOK, nil
}
