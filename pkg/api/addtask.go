package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"todo/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "некорректный JSON", http.StatusBadRequest)
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "не удалось сохранить задачу", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkAndAdjustDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.ParseInLocation(DateFormat, task.Date, time.Local)
	if err != nil {
		return fmt.Errorf("некорректный формат даты")
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if isBefore(t, now) {
			task.Date = next
		}
	} else if isBefore(t, now) {
		task.Date = now.Format(DateFormat)
	}

	return nil
}

func isBefore(d1, d2 time.Time) bool {
	y1, m1, day1 := d1.Date()
	y2, m2, day2 := d2.Date()
	if y1 != y2 {
		return y1 < y2
	}
	if m1 != m2 {
		return m1 < m2
	}
	return day1 < day2
}
