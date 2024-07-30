package nextdate

import (
	"go_final_project/nextdate"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowInString := r.URL.Query().Get("now")
	//if nowInString == "" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte("now missing"))
	//	return
	//}
	date := r.URL.Query().Get("date")
	//if date == "" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte("date missing"))
	//	return
	//}
	repeat := r.URL.Query().Get("repeat")
	//if repeat == "" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte("repeat missing"))
	//	return
	//}
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
	w.Write([]byte(nextDate))
}
