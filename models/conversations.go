package models

import "time"

type Conversation struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Messages  []Message `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE;"`
}

type Message struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Content        string
	Role           string
	ConversationID string       `gorm:"index"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;references:ID"`
}
