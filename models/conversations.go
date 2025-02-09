package models

import "time"

type Conversation struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Messages  []Message
}

type Message struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Content        string
	Role           string
	ConversationID Conversation `gorm:"foreignKey:ID"`
}
