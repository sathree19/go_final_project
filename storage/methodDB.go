package storage

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"go_final_project/model"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(task model.Task, out model.Output) model.Output {
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

func (s ParcelStore) SelectAll(limit string) ([]int, []model.Task, error) {

	var ids []int
	var task model.Task
	var task1 []model.Task

	if limit == "ALL" {
		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
		if err != nil {
			return []int{}, []model.Task{}, err
		}
		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return []int{}, []model.Task{}, err
			}

			ids = append(ids, int(task.Id))

			task1 = append(task1, task)
		}

		if err := rows.Err(); err != nil {
			return []int{}, []model.Task{}, err
		}

		return ids, task1, nil
	} else {
		lim, err := strconv.Atoi(limit)
		if err != nil {
			return []int{}, []model.Task{}, err
		}

		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit", sql.Named("limit", lim))
		if err != nil {
			return []int{}, []model.Task{}, err
		}
		defer rows.Close()
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return []int{}, []model.Task{}, err
			}

			ids = append(ids, int(task.Id))

			task1 = append(task1, task)
		}

		if err := rows.Err(); err != nil {
			return []int{}, []model.Task{}, err
		}

		return ids, task1, nil
	}

}

func (s ParcelStore) Update(task model.Task, out model.Output) (model.Task, model.Output) {

	_, err := s.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		out.Error = errors.New("Задача не найдена")
		task = model.Task{}
		return task, out

	}
	task = model.Task{Id: task.Id, Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat}
	return task, out

}

func (s ParcelStore) SelectId(param1 int) (model.Task, error) {
	var task model.Task
	var out model.Output

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", param1))

	out.Error = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, out.Error
}

func (s ParcelStore) Delete(param1 int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", param1))
	return err
}

func (s ParcelStore) Search(param string, limit string) (error, []model.Task) {

	lim, err := strconv.Atoi(limit)
	if err != nil {
		return err, []model.Task{}
	}

	param1, err := time.Parse("02.01.2006", param)
	t := param1.Format("20060102")

	if err != nil {

		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%' ORDER BY date LIMIT :limit ", sql.Named("search", param), sql.Named("search", param), sql.Named("limit", lim))
		if err != nil {

			return err, []model.Task{}
		}
		defer rows.Close()
		var task model.Task
		var task1 []model.Task
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return err, []model.Task{}
			}
			task1 = append(task1, task)

		}
		if task1 == nil {
			task1 = []model.Task{}
		}

		if err := rows.Err(); err != nil {
			return err, []model.Task{}
		}

		return nil, task1
	} else {
		rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit ", sql.Named("date", t), sql.Named("limit", lim))
		if err != nil {
			return err, []model.Task{}
		}
		defer rows.Close()
		var task model.Task
		var task1 []model.Task
		for rows.Next() {

			err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return err, []model.Task{}
			}
			task1 = append(task1, task)

		}
		if task1 == nil {
			task1 = []model.Task{}
		}

		if err := rows.Err(); err != nil {
			return err, []model.Task{}
		}

		return nil, task1
	}

}
