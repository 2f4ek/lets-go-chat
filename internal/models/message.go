package models

import "time"

type Message struct {
	Id      int       `gorm:"primaryKey"`
	Message string    `gorm:"type:text" json:"content"`
	SendAt  time.Time `gorm:"index"`
}
