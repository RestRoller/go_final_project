package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

const defaultTasksLimit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")
	tasks, err := db.Tasks(defaultTasksLimit, search)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	WriteJSON(w, TasksResponse{Tasks: tasks})
}
