package storage

import (
	"database/sql"
	"errors"
	"go_final_project/str"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(task str.Task, out str.Output) str.Output {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		out.ID = 0
		out.Error = err
		return out
	}

	out.ID, out.Error = res.LastInsertId()

	if err != nil {
		out.ID = 0
		out.Error = err
		return out
	}
	return out

}

func (s ParcelStore) SelectAll(limit string) ([]int, map[string][]str.Task, error) {

	var ids []int
	//var rows *sql.Rows
	var task str.Task
	var task1 []str.Task
	tasks := make(map[string][]str.Task)

	if limit == "ALL" {
		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
		if err != nil {
			return []int{}, map[string][]str.Task{}, err
		}
		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return []int{}, map[string][]str.Task{}, err
			}

			ids = append(ids, int(task.Id))

			task1 = append(task1, task)
		}

		tasks["tasks"] = task1

		if err := rows.Err(); err != nil {
			return []int{}, map[string][]str.Task{}, err
		}

		return ids, tasks, nil
	} else {
		lim, err := strconv.Atoi(limit)
		if err != nil {
			return []int{}, map[string][]str.Task{}, err
		}

		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", lim))
		if err != nil {
			return []int{}, map[string][]str.Task{}, err
		}
		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return []int{}, map[string][]str.Task{}, err
			}

			ids = append(ids, int(task.Id))

			task1 = append(task1, task)
		}

		tasks["tasks"] = task1

		if err := rows.Err(); err != nil {
			return []int{}, map[string][]str.Task{}, err
		}

		return ids, tasks, nil
	}

}

func (s ParcelStore) Update(task str.Task, out str.Output) (str.Task, str.Output) {

	_, err := s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		out.Error = errors.New("Задача не найдена")
		task = str.Task{}
		return task, out

	}
	task = str.Task{Id: task.Id, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat}
	return task, out

}

func (s ParcelStore) SelectId(param1 int) (str.Task, error) {
	var task str.Task
	var out str.Output

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", param1))

	out.Error = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, out.Error
}

func (s ParcelStore) Delete(param1 int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", param1))
	return err
}

func (s ParcelStore) Search(param string, limit int) (error, map[string][]str.Task) {
	tasks := make(map[string][]str.Task)

	param1, err := time.Parse("02.01.2006", param)
	t := param1.Format("20060102")

	if err != nil {

		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%' ORDER BY date LIMIT :limit ", sql.Named("search", param), sql.Named("search", param), sql.Named("limit", limit))
		if err != nil {

			return err, map[string][]str.Task{}
		}
		defer rows.Close()
		var task str.Task
		var task1 []str.Task
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return err, map[string][]str.Task{}
			}
			task1 = append(task1, task)

		}
		if task1 == nil {
			task1 = []str.Task{}
		}
		tasks["tasks"] = task1

		if err := rows.Err(); err != nil {
			return err, map[string][]str.Task{}
		}

		return nil, tasks
	} else {
		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit ", sql.Named("date", t), sql.Named("limit", limit))
		if err != nil {
			return err, map[string][]str.Task{}
		}
		defer rows.Close()
		var task str.Task
		var task1 []str.Task
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return err, map[string][]str.Task{}
			}
			task1 = append(task1, task)

		}
		if task1 == nil {
			task1 = []str.Task{}
		}
		tasks["tasks"] = task1

		if err := rows.Err(); err != nil {
			return err, map[string][]str.Task{}
		}

		return nil, tasks
	}

}
