package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
)

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		WriteError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format(dateFormat)

	// Если дата не указана, ставим сегодня
	if req.Date == "" {
		req.Date = today
	}

	// Проверяем корректность даты
	t, err := time.Parse(dateFormat, req.Date)
	if err != nil {
		WriteError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	// Получаем сегодняшнюю дату для сравнения (без времени)
	todayTime, _ := time.Parse(dateFormat, today)

	// Если есть правило повторения
	if req.Repeat != "" {
		// Проверяем правило
		_, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Сравниваем только даты (без времени)
		if t.Before(todayTime) {
			nextDate, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			req.Date = nextDate
		}
		// Если дата равна сегодня или в будущем - оставляем как есть
	} else {
		// Если нет правила повторения и дата в прошлом, ставим сегодня
		if t.Before(todayTime) {
			req.Date = today
		}
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	id, err := db.AddTask(task)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, map[string]interface{}{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	WriteJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		WriteError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		WriteError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
	today := now.Format(dateFormat)

	if req.Date == "" {
		req.Date = today
	}

	t, err := time.Parse(dateFormat, req.Date)
	if err != nil {
		WriteError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	todayTime, _ := time.Parse(dateFormat, today)

	if req.Repeat != "" {
		_, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if t.Before(todayTime) {
			nextDate, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			req.Date = nextDate
		}
	} else {
		if t.Before(todayTime) {
			req.Date = today
		}
	}

	task := &db.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	err = db.UpdateTask(task)
	if err != nil {
		WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	WriteJSON(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		WriteError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	WriteJSON(w, map[string]interface{}{})
}
