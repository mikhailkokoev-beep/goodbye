package api

import (
	"net/http"

	"todo/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, "ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
