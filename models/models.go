package models

import "github.com/google/uuid"

type JsonPlaceholder struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Type  string  `json:"type"`
	Owner string  `json:"owner"`
}

type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type User struct {
	Uid      uuid.UUID `json:"userid"`
	UserData `json:"userdata"`
}

type RegisterData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Status struct {
	CurStatus string `json:"pg_status"`
}
