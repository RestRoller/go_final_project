package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила")
	}

	switch parts[0] {
	case "y":
		// Ежегодно
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "d":
		// Через N дней
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("неверный интервал для d: должен быть от 1 до 400")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "w":
		// По дням недели
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w")
		}
		weekdays := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, d := range weekdays {
			num, err := strconv.Atoi(d)
			if err != nil || num < 1 || num > 7 {
				return "", fmt.Errorf("неверный день недели: %s (должен быть от 1 до 7)", d)
			}
			allowedDays[num] = true
		}

		// Начинаем со следующего дня
		date = date.AddDate(0, 0, 1)
		for {
			if date.After(now) {
				weekday := int(date.Weekday())
				if weekday == 0 {
					weekday = 7 // воскресенье
				}
				if allowedDays[weekday] {
					return date.Format(dateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		// По дням месяца
		if len(parts) < 2 {
			return "", fmt.Errorf("неверный формат правила m")
		}

		// Парсим дни
		daysPart := parts[1]
		days := strings.Split(daysPart, ",")
		allowedDays := make(map[int]bool)
		for _, d := range days {
			num, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("неверный день месяца: %s", d)
			}
			if num < -31 || num == 0 || num > 31 {
				return "", fmt.Errorf("неверный день месяца: %d (должен быть от -31 до -1 или от 1 до 31)", num)
			}
			allowedDays[num] = true
		}

		// Парсим месяцы (опционально)
		allowedMonths := make(map[int]bool)
		if len(parts) > 2 {
			monthsPart := parts[2]
			months := strings.Split(monthsPart, ",")
			for _, m := range months {
				num, err := strconv.Atoi(m)
				if err != nil || num < 1 || num > 12 {
					return "", fmt.Errorf("неверный месяц: %s (должен быть от 1 до 12)", m)
				}
				allowedMonths[num] = true
			}
		}

		// Ищем следующую дату
		date = date.AddDate(0, 0, 1)
		for {
			if date.After(now) {
				month := int(date.Month())
				// Проверяем месяц, если указаны конкретные месяцы
				if len(allowedMonths) > 0 && !allowedMonths[month] {
					// Переходим на первый день следующего месяца
					date = date.AddDate(0, 1, -date.Day()+1)
					continue
				}

				day := date.Day()
				lastDay := date.AddDate(0, 1, -day).Day()

				// Проверяем день
				for allowedDay := range allowedDays {
					targetDay := allowedDay
					if allowedDay < 0 {
						targetDay = lastDay + allowedDay + 1
					}
					if targetDay >= 1 && targetDay <= lastDay && day == targetDay {
						return date.Format(dateFormat), nil
					}
				}
			}
			date = date.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", parts[0])
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr != "" {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid now date", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}
