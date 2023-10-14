package model

type GPT3Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPT3Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type User struct {
	UserId   string `json:"userID" validate:"required, gte=3"`
	Password string `json:"password"`
	ChatID   string `json:"chatID"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type CreateChatRequest struct {
	UserId string `json:"userId"`
	Prompt string `json:"Prompt"`
}
