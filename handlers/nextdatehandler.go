package nextdate

import (
	"go_final_project/nextdate"
	"log"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowInString := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	now, err := time.Parse("20060102", nowInString)
	if err != nil {
		http.Error(w, "время не может быть преобразовано в корректную дату", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Printf("Ошибка при записи ответа в NextDateHandler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
