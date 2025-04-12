package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID   string    `json:"id" db:"id"`
	Name string `json:"name" db:"username"`
	Password string `json:"password" db:"password"`
	Email string `json:"email" db:"email"`
}
