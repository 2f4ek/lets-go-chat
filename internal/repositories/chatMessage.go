package repositories

import (
	"time"

	"github.com/2f4ek/lets-go-chat/database"
	"github.com/2f4ek/lets-go-chat/internal/models"
)

var CMRInstanse *ChatMessageRepository

type ChatMessageRepository struct {
	db database.Database
}

func ProvideChatMessageRepository(db *database.Database) *ChatMessageRepository {
	once.Do(func() {
		CMRInstanse = &ChatMessageRepository{}
		CMRInstanse.db = *db
	})
	return CMRInstanse
}

func (cm *ChatMessageRepository) Save(m models.Message) (*models.Message, error) {
	m.SendAt = time.Now()
	err := cm.db.GetDatabase().Create(&m).Error
	if err != nil {
		return &models.Message{}, err
	}
	return &m, nil
}

func (cm *ChatMessageRepository) GetMissedMessages(user *models.User) ([]*models.Message, error) {
	var messages []*models.Message
	err := cm.db.GetDatabase().Where("send_at >= ?", user.LastActivity).Find(&messages).Error
	if err != nil {
		return messages, err
	}
	return messages, nil
}
