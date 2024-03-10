package models

type User struct {
	ID   	       int    `json:"id"`
	Username       string `json:"username"`
	Password_hash  bool   `json:"password_hash"`
}