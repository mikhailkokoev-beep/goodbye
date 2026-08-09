package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "ошибка сериализации JSON", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func writeError(w http.ResponseWriter, text string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	b, _ := json.Marshal(map[string]string{"error": text})
	w.Write(b)
}
