package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(id, date string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	var query string
	var args []interface{}

	if search != "" {
		// Проверяем, является ли search датой в формате DD.MM.YYYY
		if len(search) == 10 && search[2] == '.' && search[5] == '.' {
			// Преобразуем DD.MM.YYYY в YYYYMMDD
			date := search[6:10] + search[3:5] + search[0:2]
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			args = []interface{}{date, limit}
		} else {
			// Поиск по заголовку или комментарию
			searchPattern := "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			args = []interface{}{searchPattern, searchPattern, limit}
		}
	} else {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		args = []interface{}{limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
