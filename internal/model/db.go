package model

type News struct {
	Id         int64   `json:"Id" db:"id"`
	Title      string  `json:"Title" db:"title"`
	Content    string  `json:"Content" db:"content"`
	Categories []int64 `json:"Categories" db:"categories"`
}
