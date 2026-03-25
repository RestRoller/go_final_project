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
		nextDate := date
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format(dateFormat), nil
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("неверный интервал для d: должен быть от 1 до 400")
		}

		// Начинаем с начальной даты
		nextDate := date

		// Добавляем интервал пока дата не станет больше now
		for {
			nextDate = nextDate.AddDate(0, 0, interval)
			if nextDate.After(now) {
				return nextDate.Format(dateFormat), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w")
		}
		weekdays := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, d := range weekdays {
			num, err := strconv.Atoi(d)
			if err != nil || num < 1 || num > 7 {
				return "", fmt.Errorf("неверный день недели: %s", d)
			}
			allowedDays[num] = true
		}

		nextDate := date
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if nextDate.After(now) {
				weekday := int(nextDate.Weekday())
				if weekday == 0 {
					weekday = 7
				}
				if allowedDays[weekday] {
					return nextDate.Format(dateFormat), nil
				}
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("неверный формат правила m")
		}

		daysPart := parts[1]
		days := strings.Split(daysPart, ",")
		allowedDays := make(map[int]bool)
		for _, d := range days {
			num, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("неверный день месяца: %s", d)
			}
			if num < -31 || num == 0 || num > 31 {
				return "", fmt.Errorf("неверный день месяца: %d", num)
			}
			allowedDays[num] = true
		}

		allowedMonths := make(map[int]bool)
		if len(parts) > 2 {
			monthsPart := parts[2]
			months := strings.Split(monthsPart, ",")
			for _, m := range months {
				num, err := strconv.Atoi(m)
				if err != nil || num < 1 || num > 12 {
					return "", fmt.Errorf("неверный месяц: %s", m)
				}
				allowedMonths[num] = true
			}
		}

		nextDate := date
		for {
			nextDate = nextDate.AddDate(0, 0, 1)
			if nextDate.After(now) {
				month := int(nextDate.Month())
				if len(allowedMonths) > 0 && !allowedMonths[month] {
					nextDate = time.Date(nextDate.Year(), nextDate.Month()+1, 1, 0, 0, 0, 0, nextDate.Location())
					continue
				}

				day := nextDate.Day()
				lastDay := nextDate.AddDate(0, 1, -day).Day()

				for allowedDay := range allowedDays {
					targetDay := allowedDay
					if allowedDay < 0 {
						targetDay = lastDay + allowedDay + 1
					}
					if targetDay >= 1 && targetDay <= lastDay && day == targetDay {
						return nextDate.Format(dateFormat), nil
					}
				}
			}
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
