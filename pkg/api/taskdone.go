package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		WriteError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		WriteJSON(w, map[string]interface{}{})
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateDate(id, nextDate)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, map[string]interface{}{})
}
