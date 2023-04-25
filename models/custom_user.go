package models

type CustomUser struct {
	ID          int64
	Email       string
	Role        Role
	PhoneNumber string
	Name        string
	UID         string
}
