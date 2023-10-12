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
	Email    string `json:"email" validate:"required, gte=3"`
	Password string `json:"password"`
	Prompt   string `json:"prompt"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}
