package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"vartalabh.com/m/model"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.SendMessageRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		chats, err := FetchParticularChat(req.ChatId)
		if err != nil {
			model.WriteOutput(w, "Error while fetching user chats", http.StatusForbidden, err)
			return
		}
		messages := chats.Messages
		messages = append(messages, model.Message{
			Role:    "user",
			Content: req.Message,
		})
		response := makeChatGptCall(messages)
		messages = append(messages, model.Message{
			Role:    "system",
			Content: response.Content,
		})
		marshalledMessages, err := json.Marshal(messages)
		if err != nil {
			model.WriteOutput(w, "Error while marshalling message", http.StatusForbidden, err)
			return
		}
		err = UpdateChatMessages(req.ChatId, marshalledMessages)
		if err != nil {
			model.WriteOutput(w, "Error while updating users's chat messages", http.StatusForbidden, err)
			return
		}
		model.WriteOutput(w, response.Content, http.StatusOK, err)
	}
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var req model.GetChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		chats, err := FetchParticularChat(req.ChatId)
		if err != nil {
			model.WriteOutput(w, "Error while fetching user chats", http.StatusForbidden, err)
			return
		}
		model.WriteOutput(w, chats, http.StatusOK, err)
	}
}

func GetChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var req model.GetChatHistoryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		chats, err := FetchUserChats(req.UserId)
		if err != nil {
			model.WriteOutput(w, "Error while fetching user chats", http.StatusForbidden, err)
			return
		}
		model.WriteOutput(w, chats, http.StatusOK, err)
	}
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.CreateChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		chatId := uuid.New().String()
		prompt := req.Prompt
		if prompt == "" {
			prompt = "you are a mental health counsellor. Talk to user, ask repititive questions,keep the conversation going. Also ask the user how much progress he made based on the prompt provided."
		}
		message := []model.Message{
			{
				Role:    "system",
				Content: prompt,
			},
		}
		messages, err := json.Marshal(message)
		if err != nil {
			model.WriteOutput(w, "Error while marshalling message", http.StatusForbidden, err)
			return
		}
		err = UpdateUserCurrentChat(req.UserId, chatId)
		if err != nil {
			model.WriteOutput(w, "Error while updating users's current chat", http.StatusForbidden, err)
			return
		}
		err = CreateChatEntry(req.UserId, chatId, prompt, messages)
		if err != nil {
			model.WriteOutput(w, "Error while creating chat entry", http.StatusForbidden, err)
			return
		}
		data := map[string]interface{}{
			"chatId":  chatId,
			"message": "Successfully created new chat for user",
		}
		model.WriteOutput(w, data, http.StatusOK, err)
	}
}

func makeChatGptCall(messages []model.Message) *model.Message {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	endpoint := "https://api.openai.com/v1/chat/completions"

	reqBody := &model.GPT3Request{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	requestBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling JSON request:", err)
		return nil
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	// Parse the JSON response
	var gpt3Response model.GPT3Response
	if err := json.Unmarshal(body, &gpt3Response); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil
	}

	return &gpt3Response.Choices[0].Message
}
