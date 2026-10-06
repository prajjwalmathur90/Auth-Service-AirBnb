package models

type User struct {
	Id        int64
	Username  string
	Email     string
	Password  string `json:"-"`
	CreatedAt string
	UpdatedAt string
}