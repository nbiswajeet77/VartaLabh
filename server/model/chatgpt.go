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
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
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

type WaitlistEntryRequest struct {
	EmailId string `json:"emailId"`
}

type GetChatHistoryRequest struct {
	UserId string `json:"userId"`
}

type ChatHistoryResponse struct {
	ChatId  string `json:"chatId"`
	Prompt  string `json:"prompt"`
	Summary string `json:"summary"`
}

type GetChatRequest struct {
	ChatId string `json:"chatId"`
}

type ExitChatRequest struct {
	ChatId string `json:"chatId"`
}

type EditPromptRequest struct {
	ChatId string `json:"chatId"`
	Prompt string `json:"prompt"`
}

type DeleteChatRequest struct {
	ChatId string `json:"chatId"`
}

type SendMessageRequest struct {
	ChatId  string `json:"chatId"`
	Message string `json:"message"`
}

type GetChatResponse struct {
	ChatId   string    `json:"chatId"`
	Messages []Message `json:"messages"`
	Prompt   string    `json:"prompt"`
	Summary  string    `json:"summary"`
}
