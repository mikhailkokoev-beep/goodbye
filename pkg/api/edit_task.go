package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"todo/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		writeError(w, "некорректный идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkAndAdjustDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}
