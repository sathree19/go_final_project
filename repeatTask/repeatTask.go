package repeatTask

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDayOfWeek(needDayOfWeek []int, dayOfWeek int, t time.Time) (time.Time, error) {
	temp := 0

	if dayOfWeek == 0 {
		tNew := t.AddDate(0, 0, needDayOfWeek[0])
		return tNew, nil
	}
	for i := 0; i < len(needDayOfWeek); i++ {
		needDayWeek := needDayOfWeek[i]

		if (needDayWeek - dayOfWeek) > 0 {
			tNew := t.AddDate(0, 0, (needDayWeek - dayOfWeek))
			return tNew, nil
		}
		if (needDayWeek - dayOfWeek) <= 0 {
			temp++
		}
	}
	if temp == len(needDayOfWeek) {
		tNew := t.AddDate(0, 0, (7-dayOfWeek)+needDayOfWeek[len(needDayOfWeek)-1])
		return tNew, nil
	}
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), errors.New("не указаны дни недели, удалить задачу посе выполнения")
}

func NextDayOfMounth(needDayOfMounth []int, numberOfMounth []int, t time.Time) (time.Time, error) {
	temp := 0

	if needDayOfMounth[0] == -2 || needDayOfMounth[0] == -1 {
		needDayOfMounth[0]++
	} else if len(needDayOfMounth) > 1 && needDayOfMounth[1] == -1 {
		needDayOfMounth[1]++
	}

	slices.Sort(needDayOfMounth)

	for i := 0; i < len(needDayOfMounth); i++ {
		needDayMounth := needDayOfMounth[i]

		if (needDayMounth-t.Day()) > 0 && needDayMounth <= time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day() {
			tNew := time.Date(t.Year(), t.Month(), needDayMounth, 0, 0, 0, 0, time.Local)
			return tNew, nil
		}
		if (needDayMounth - t.Day()) <= 0 {
			temp++
		}
	}
	if temp == len(needDayOfMounth) {

		temp1 := 1

		for _, i := range numberOfMounth {
			mounth := t.AddDate(0, i, 0)

			if needDayOfMounth[len(needDayOfMounth)-1] <= time.Date(t.Year(), mounth.Month()+1, 0, 0, 0, 0, 0, time.Local).Day() {

				tNew := time.Date(t.Year(), mounth.Month(), needDayOfMounth[0], 0, 0, 0, 0, time.Local)
				return tNew, nil
			}
			temp1++

		}
		mounth := t.AddDate(0, temp1, 0)
		tNew := time.Date(t.Year(), mounth.Month(), needDayOfMounth[0], 0, 0, 0, 0, time.Local)
		return tNew, nil

	}
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), errors.New("не указаны дни недели, удалить задачу посе выполнения")
}

func NextMounth(needMounth []int, t time.Time) (time.Time, error) {
	temp := 0

	slices.Sort(needMounth)

	for i := 0; i < len(needMounth); i++ {
		needNumberMounth := needMounth[i]

		if (needNumberMounth - int(t.Month())) > 0 {
			mounth := t.AddDate(0, needNumberMounth-int(t.Month()), 0)

			tNew := time.Date(t.Year(), mounth.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			return tNew, nil
		}
		if (needNumberMounth - int(t.Month())) <= 0 {
			temp++
		}
	}
	if temp == len(needMounth) {
		mounth := t.AddDate(0, needMounth[len(needMounth)-1]-int(t.Month()), 0)

		tNew := time.Date(t.Year(), mounth.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		return tNew, nil
	}

	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), errors.New("ошибка, удалить задачу посе выполнения")
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	symbolForRepeat := strings.Split(repeat, " ")
	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch {
	case len(symbolForRepeat) == 1 && symbolForRepeat[0] == "":
		err := errors.New("пустая строка, удалить задачу посе выполнения")
		return "", err
	case len(symbolForRepeat) == 1 && symbolForRepeat[0] == "y":
		tNew := t.AddDate(1, 0, 0)
		if tNew.Sub(now) < 0 {
			tNew = tNew.AddDate(now.Year()-tNew.Year(), 0, 0)
		}
		return tNew.Format("20060102"), nil
	case len(symbolForRepeat) == 2 && symbolForRepeat[0] == "d":
		day, err := strconv.Atoi(symbolForRepeat[1])
		if err != nil {
			return "", err
		}
		if day > 400 {
			err = errors.New("перенос на больше чем 400 дней не возможен, удалить задачу посе выполнения")
			return "", err
		}
		tNew := t.AddDate(0, 0, day)
		if tNew.Sub(now) < 0 {
			DurationDay := (float64(now.Day()) - float64(tNew.Day()))

			tNew = tNew.AddDate(0, 0, int((math.Ceil(DurationDay/float64(day)))*float64(day)))
			//fmt.Println((math.Ceil(DurationDay / 12)))
		}
		return tNew.Format("20060102"), nil

	case len(symbolForRepeat) == 2 && symbolForRepeat[0] == "w":
		dayWeekRepeat := strings.Split(symbolForRepeat[1], ",")
		var intDayWeekRepeat []int
		for _, needDayWeek := range dayWeekRepeat {
			nD, err := strconv.Atoi(needDayWeek)
			if err != nil {
				return "", err
			}
			intDayWeekRepeat = append(intDayWeekRepeat, nD)
		}

		slices.Sort(intDayWeekRepeat)

		if intDayWeekRepeat[0] <= 0 || intDayWeekRepeat[len(intDayWeekRepeat)-1] > 7 {
			return "", errors.New("указаны неверные дни недели, удалить задачу посе выполнения")
		}

		day := int(t.Weekday())

		tNew, err := NextDayOfWeek(intDayWeekRepeat, day, t)

		if tNew.Sub(now) < 0 {
			tNew, err = NextDayOfWeek(intDayWeekRepeat, int(now.Weekday()), now)
		}
		return tNew.Format("20060102"), err

	case len(symbolForRepeat) == 2 && symbolForRepeat[0] == "m":
		dayMounthRepeat := strings.Split(symbolForRepeat[1], ",")
		var intDayMounthRepeat []int
		for _, needDayMounth := range dayMounthRepeat {
			mD, err := strconv.Atoi(needDayMounth)
			if err != nil {
				return "", err
			}
			intDayMounthRepeat = append(intDayMounthRepeat, mD)
		}

		slices.Sort(intDayMounthRepeat)

		numberOfMonth := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		if intDayMounthRepeat[0] < -2 || intDayMounthRepeat[len(intDayMounthRepeat)-1] > 31 {
			return "", errors.New("указаны неверные дни месяца, удалить задачу посе выполнения")
		}
		tNew, err := NextDayOfMounth(intDayMounthRepeat, numberOfMonth, t)

		if tNew.Sub(now) < 0 {
			tNew, err = NextDayOfMounth(intDayMounthRepeat, numberOfMonth, now)
		}
		return tNew.Format("20060102"), err

	case len(symbolForRepeat) == 3 && symbolForRepeat[0] == "m":
		mounthRepeat := strings.Split(symbolForRepeat[2], ",")
		var intMounthRepeat []int
		for _, needMounth := range mounthRepeat {
			mD, err := strconv.Atoi(needMounth)
			if err != nil {
				return "", err
			}
			intMounthRepeat = append(intMounthRepeat, mD)
		}

		slices.Sort(intMounthRepeat)

		if intMounthRepeat[0] < 1 || intMounthRepeat[len(intMounthRepeat)-1] > 12 {
			return "", errors.New("указаны неверные номера месяцев, удалить задачу посе выполнения")
		}

		tNewofMounth, err := NextMounth(intMounthRepeat, t)
		if err != nil {
			return "", err
		}

		dayMounthRepeat := strings.Split(symbolForRepeat[1], ",")
		var intDayMounthRepeat []int
		for _, needDayMounth := range dayMounthRepeat {
			mD, err := strconv.Atoi(needDayMounth)
			if err != nil {
				return "", err
			}
			intDayMounthRepeat = append(intDayMounthRepeat, mD)
		}

		slices.Sort(intDayMounthRepeat)

		if intDayMounthRepeat[0] <= -2 || intDayMounthRepeat[len(intDayMounthRepeat)-1] > 31 {
			return "", errors.New("указаны неверные дни месяца, удалить задачу посе выполнения")
		}

		tNew, err := NextDayOfMounth(intDayMounthRepeat, intMounthRepeat, tNewofMounth)

		if err != nil {
			return "", err
		}

		if tNew.Sub(now) < 0 {
			tNewofMounth, err := NextMounth(intMounthRepeat, t)
			if err != nil {
				return "", err
			}
			tNew, err = NextDayOfMounth(intDayMounthRepeat, intMounthRepeat, tNewofMounth)
		}
		return tNew.Format("20060102"), err

	default:
		err := errors.New("указан неверный формат, удалить задачу посе выполнения")
		return "", err
	}
}

func MainHandle(w http.ResponseWriter, r *http.Request) {

	now1, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		fmt.Println(err)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDateTask, err := NextDate(now1, date, repeat)

	//fmt.Fprintf(w, "Новая дата: %s ошибка: %s", nextDateTask, err)
	fmt.Fprintf(w, "%s", nextDateTask)
}
