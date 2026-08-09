package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.ParseInLocation(DateFormat, nowStr, time.Local)
		if err != nil {
			http.Error(w, "некорректный формат now", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.ParseInLocation(DateFormat, dstart, time.Local)
	if err != nil {
		return "", errors.New("некорректная дата: " + dstart)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		return nextDays(date, now, parts)
	case "y":
		return nextYears(date, now)
	case "w":
		return nextWeeks(date, now, parts)
	case "m":
		return nextMonths(date, now, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

func isAfter(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	if y1 != y2 {
		return y1 > y2
	}
	if m1 != m2 {
		return m1 > m2
	}
	return d1 > d2
}

func nextDays(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила d")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days < 1 || days > 400 {
		return "", errors.New("интервал d должен быть от 1 до 400")
	}
	for {
		date = date.AddDate(0, 0, days)
		if isAfter(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextYears(date, now time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if isAfter(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextWeeks(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила w")
	}

	allowed := make(map[time.Weekday]bool)
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(s)
		if err != nil || d < 1 || d > 7 {
			return "", errors.New("день недели должен быть от 1 до 7")
		}
		allowed[time.Weekday(d%7)] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		if isAfter(date, now) && allowed[date.Weekday()] {
			return date.Format(DateFormat), nil
		}
	}
}

func nextMonths(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("неверный формат правила m")
	}

	allowedDays := make(map[int]bool)
	lastDay, prevLastDay := false, false
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(s)
		if err != nil || d == 0 || d > 31 || d < -2 {
			return "", errors.New("недопустимый день месяца")
		}
		switch d {
		case -1:
			lastDay = true
		case -2:
			prevLastDay = true
		default:
			allowedDays[d] = true
		}
	}

	allowedMonths := make(map[int]bool)
	if len(parts) == 3 {
		for _, s := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(s)
			if err != nil || m < 1 || m > 12 {
				return "", errors.New("недопустимый месяц")
			}
			allowedMonths[m] = true
		}
	} else {
		for m := 1; m <= 12; m++ {
			allowedMonths[m] = true
		}
	}

	for i := 0; i < 366*5; i++ {
		date = date.AddDate(0, 0, 1)

		if !isAfter(date, now) || !allowedMonths[int(date.Month())] {
			continue
		}

		day := date.Day()
		last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()

		if allowedDays[day] || (lastDay && day == last) || (prevLastDay && day == last-1) {
			return date.Format(DateFormat), nil
		}
	}
	return "", errors.New("не найдена подходящая дата (проверьте правило)")
}
