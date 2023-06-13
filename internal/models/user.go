package models

import "time"

type UserId int

type User struct {
	Id           UserId
	Name         string
	Password     string
	Token        string
	LastActivity time.Time
}
