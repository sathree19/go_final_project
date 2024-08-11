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

func Contains(x []int, y int) bool {
	for _, v := range x {
		if v == y {
			return true
		}
	}
	return false
}

func Next(x []int, y int) int {

	if len(x) == 1 {
		return x[0]
	}

	for _, v := range x {
		if y < v {
			return v
		}
	}

	return 0
}

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
		temp++
	}
	if temp == len(needDayOfWeek) {
		tNew := t.AddDate(0, 0, (7-dayOfWeek)+needDayOfWeek[0])
		return tNew, nil
	}
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local), errors.New("не указаны дни недели, удалить задачу посе выполнения")
}

func NextDayOfMounth(needDayOfMounth []int, numberOfMounth []int, t time.Time) (time.Time, error) {

	temp := 0
	var needDayOfMounthTemp []int
	// Переводим последний и предпоследний день в обычные для текущего месяца

	if needDayOfMounth[0] == -2 || needDayOfMounth[0] == -1 {
		needDayOfMounth[0]++
	}
	if len(needDayOfMounth) > 1 && needDayOfMounth[1] == -1 {
		needDayOfMounth[1]++
	}

	if needDayOfMounth[0] == -1 {
		needDayOfMounth[0] = int(time.Date(t.Year(), t.Month()+1, -1, 0, 0, 0, 0, time.Local).Day())
		temp = 1
	}

	if needDayOfMounth[0] == 0 {
		needDayOfMounth[0] = int(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day())
		temp = 2
	}

	if len(needDayOfMounth) > 1 && needDayOfMounth[1] == 0 {
		needDayOfMounth[1] = int(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day())
		temp = 3
	}

	// Сортируем с последним и предпоследним днем

	slices.Sort(needDayOfMounth)

	// Убираем дубли

	for j, x := range needDayOfMounth {
		for i := j + 1; i < len(needDayOfMounth); i++ {

			if x == needDayOfMounth[i] {

				copy(needDayOfMounth[i:], needDayOfMounth[i+1:])
				i--
				needDayOfMounth = needDayOfMounth[:len(needDayOfMounth)-1]

			}
		}

	}

	//////////////////////////////

	if Contains(numberOfMounth, int(t.Month())) && t.Day() < needDayOfMounth[len(needDayOfMounth)-1] {

		for i, d := range needDayOfMounth {
			if d > time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local).Day() {
				needDayOfMounthTemp = needDayOfMounth[:i+1]
				break
			}
		}

		if len(needDayOfMounthTemp) != 0 {
			if len(needDayOfMounthTemp) == 1 {
				nextMounth, year := NextMounth(numberOfMounth, t)
				return time.Date(year, nextMounth, Next(needDayOfMounthTemp, t.Day()), 0, 0, 0, 0, time.Local), nil
			}
			tNew := time.Date(t.Year(), t.Month(), Next(needDayOfMounthTemp, t.Day()), 0, 0, 0, 0, time.Local)
			return tNew, nil
		}
		// nextMounth, year := NextMounth(numberOfMounth, t)

		tNew := time.Date(t.Year(), t.Month(), Next(needDayOfMounth, t.Day()), 0, 0, 0, 0, time.Local)
		return tNew, nil

	}

	if temp != 0 {
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	}

	nextMounth, year := NextMounth(numberOfMounth, t)

	if temp == 1 {
		needDayOfMounth = append(needDayOfMounth, int(time.Date(t.Year(), nextMounth+1, -1, 0, 0, 0, 0, time.Local).Day()))
		if len(needDayOfMounth) == 2 {
			return time.Date(year, nextMounth, needDayOfMounth[1], 0, 0, 0, 0, time.Local), nil
		}
	}

	if temp == 2 {
		needDayOfMounth = append(needDayOfMounth, int(time.Date(t.Year(), nextMounth+1, 0, 0, 0, 0, 0, time.Local).Day()))
		if len(needDayOfMounth) == 2 {
			return time.Date(year, nextMounth, needDayOfMounth[1], 0, 0, 0, 0, time.Local), nil
		}
	}

	if temp == 3 {
		needDayOfMounth = append(needDayOfMounth, int(time.Date(t.Year(), nextMounth+1, -1, 0, 0, 0, 0, time.Local).Day()))
		needDayOfMounth = append(needDayOfMounth, int(time.Date(t.Year(), nextMounth+1, 0, 0, 0, 0, 0, time.Local).Day()))
		if len(needDayOfMounth) == 4 {
			needDayOfMounth1 := []int{int(time.Date(t.Year(), nextMounth+1, -1, 0, 0, 0, 0, time.Local).Day()), int(time.Date(t.Year(), nextMounth+1, 0, 0, 0, 0, 0, time.Local).Day())}
			tNew := time.Date(year, nextMounth, needDayOfMounth1[0], 0, 0, 0, 0, time.Local)
			return tNew, nil
		}
	}

	// Сортируем с последним и предпоследним днем

	slices.Sort(needDayOfMounth)

	// Убираем дубли

	for j, x := range needDayOfMounth {
		for i := j + 1; i < len(needDayOfMounth); i++ {

			if x == needDayOfMounth[i] {

				copy(needDayOfMounth[i:], needDayOfMounth[i+1:])
				i--
				needDayOfMounth = needDayOfMounth[:len(needDayOfMounth)-1]

			}
		}

	}

	tNew := time.Date(year, nextMounth, needDayOfMounth[0], 0, 0, 0, 0, time.Local)
	return tNew, nil

}

func NextMounth(needMounth []int, t time.Time) (time.Month, int) {

	slices.Sort(needMounth)

	if int(t.Month()) >= needMounth[len(needMounth)-1] {
		mounth := t.AddDate(0, 12-int(t.Month())+needMounth[0], 0)
		return mounth.Month(), t.Year() + 1
	}

	// mounth := t.AddDate(0, Next(needMounth, int(t.Month()))-int(t.Month()), 0)

	// if Next(needMounth, int(t.Month())) == int(mounth.Month()) {

	// 	return mounth.Month(), t.Year()

	// }

	tNew := time.Date(t.Year(), time.Month(Next(needMounth, int(t.Month()))), t.Day(), 0, 0, 0, 0, time.Local)

	if (tNew.Day() == 1 || tNew.Day() == 2 || tNew.Day() == 3) && (t.Day() == 29 || t.Day() == 30 || t.Day() == 31) {
		return time.Month(Next(needMounth, Next(needMounth, int(t.Month())))), t.Year()
	}

	return time.Month(Next(needMounth, int(t.Month()))), t.Year()

}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	symbolForRepeat := strings.Split(repeat, " ")
	t, err := time.Parse("20060102", date)
	if err != nil {
		return " ", err
	}

	switch {
	case len(symbolForRepeat) == 1 && symbolForRepeat[0] == " ":
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
			fmt.Println(err)
			return "", err
		}
		if day > 400 || day <= 0 {
			err = errors.New("перенос на больше чем 400 дней не возможен, удалить задачу посе выполнения")
			fmt.Println(err)
			return "", err
		}

		tNew := t.AddDate(0, 0, day)

		if tNew.Sub(now) < 0 {
			diff := now.Sub(tNew).Hours() / 24
			DurationDay := (float64(diff))

			tNew = tNew.AddDate(0, 0, int((math.Ceil(DurationDay/float64(day)))*float64(day)))

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

		if intDayMounthRepeat[0] < -2 || intDayMounthRepeat[len(intDayMounthRepeat)-1] > 31 {
			return "", errors.New("указаны неверные дни месяца, удалить задачу посе выполнения")
		}

		tNew, err := NextDayOfMounth(intDayMounthRepeat, intMounthRepeat, t)

		if err != nil {
			return "", err
		}

		if tNew.Sub(now) < 0 {
			tNew, err = NextDayOfMounth(intDayMounthRepeat, intMounthRepeat, now)
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

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "%s", nextDateTask)
}
