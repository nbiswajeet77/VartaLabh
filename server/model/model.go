package model

import "time"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type User struct {
	UserId    string    `json:"userID" validate:"required, gte=3"`
	Password  string    `json:"password"`
	ChatID    string    `json:"chatID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Chat struct {
	ChatId    string    `json:"chatId"`
	Messages  []Message `json:"messages"`
	Prompt    string    `json:"prompt"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
