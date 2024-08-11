package model

type Task struct {
	Id      int64  `json:"id,string,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Output struct {
	ID    int64 `json:"id,string,omitempty"`
	Error error `json:"error"`
}

const LimitShowTasks = "10"
