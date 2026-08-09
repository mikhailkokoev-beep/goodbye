package api

import (
	"net/http"
	"time"

	"todo/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.UpdateDate(nextDate, id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{})
}
