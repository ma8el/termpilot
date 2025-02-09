package db

import (
	"termpilot/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("termpilot.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB.AutoMigrate(&models.Conversation{}, &models.Message{})
	return nil
}

func GetConversation(id string) (*models.Conversation, error) {
	var conversation models.Conversation
	if err := DB.Preload("Messages").Where("id = ?", id).First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func GetAllConversations() ([]models.Conversation, error) {
	var conversations []models.Conversation
	if err := DB.Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func CreateConversation(conversation models.Conversation) (*models.Conversation, error) {
	if err := DB.Create(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func UpdateConversation(conversation models.Conversation) (*models.Conversation, error) {
	if err := DB.Save(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func DeleteConversation(id string) error {
	if err := DB.Delete(&models.Conversation{}, id).Error; err != nil {
		return err
	}
	return nil
}

func GetLastConversation() (*models.Conversation, error) {
	var conversation models.Conversation
	if err := DB.Order("created_at DESC").Preload("Messages").First(&conversation).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}
