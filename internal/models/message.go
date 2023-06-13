package models

import (
	"github.com/2f4ek/lets-go-chat/database"
	"time"
)

type Message struct {
	Id      int       `gorm:"primaryKey"`
	Message string    `gorm:"type:text" json:"content"`
	SendAt  time.Time `gorm:"index"`
}

func (m *Message) Save() (*Message, error) {
	m.SendAt = time.Now()
	err := database.Database.Create(&m).Error
	if err != nil {
		return &Message{}, err
	}
	return m, nil
}

func GetMissedMessages(user *User) ([]*Message, error) {
	var messages []*Message
	err := database.Database.Where("send_at >= ?", user.LastActivity).Find(&messages).Error
	if err != nil {
		return messages, err
	}
	return messages, nil
}
