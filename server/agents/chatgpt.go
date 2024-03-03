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

func ExitChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.ExitChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		chats, err := FetchParticularChat(req.ChatId)
		if err != nil {
			model.WriteOutput(w, "Error while fetching user chats", http.StatusForbidden, err)
			return
		} else if chats == nil {
			model.WriteOutput(w, "Chat was not found", http.StatusNotFound, err)
			return
		}

		messages := chats.Messages
		messages = append(messages, model.Message{
			Role:    "user",
			Content: "Summarise the user's chats till now within 2-3 lines.",
		})
		summary := makeChatGptCall(messages)
		msgs, err := json.Marshal(chats.Messages)
		if err != nil {
			model.WriteOutput(w, "Error while marshalling message", http.StatusForbidden, err)
			return
		}
		UpdateChatEntry(chats.ChatId, chats.Prompt, summary.Content, msgs)
		model.WriteOutput(w, "Chat Exited successfully", http.StatusOK, err)
	}
}

func GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		} else if chats == nil {
			model.WriteOutput(w, "Chat was not found", http.StatusNotFound, err)
			return
		}
		model.WriteOutput(w, chats, http.StatusOK, err)
	}
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.DeleteChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		rowsAffected, err := DeleteParticularChat(req.ChatId)
		if err != nil {
			model.WriteOutput(w, "Error while deleting user chats", http.StatusForbidden, err)
			return
		} else if rowsAffected == 0 {
			model.WriteOutput(w, "No Data was found for the provided chatId", http.StatusNotFound, err)
			return
		}

		model.WriteOutput(w, "Chat Deleted for User", http.StatusOK, err)
	}
}

func EditPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.EditPromptRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}

		chats, err := FetchParticularChat(req.ChatId)
		if err != nil {
			model.WriteOutput(w, "Error while fetching user chats in editPrompt API", http.StatusForbidden, err)
			return
		} else if chats == nil {
			model.WriteOutput(w, "Chat was not found", http.StatusNotFound, err)
			return
		}

		chats.Messages[0].Content = req.Prompt
		messages, err := json.Marshal(chats.Messages)
		if err != nil {
			model.WriteOutput(w, "Error while marshalling message", http.StatusForbidden, err)
			return
		}

		UpdateChatEntry(chats.ChatId, req.Prompt, chats.Summary, messages)

		model.WriteOutput(w, "prompt for the chatId updated successfully", http.StatusOK, err)
	}
}

func GetChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
			prompt = "As a Cognitive Behavioral Therapist, your aim is to support user through various techniques in Cognitive Behavioral Therapy (CBT). The conversation will involve identifying troubling situations and exploring user's thoughts, emotions, and beliefs about them.\nTo facilitate this, you will guide the user through the following steps:\n1. Identifying Negative Thought Patterns:\nYou will explore various negative thinking patterns that might be causing distress. These patterns include but aren't limited to:\nAll-or-Nothing Thinking\nOvergeneralization\nMental Filter\nDisqualifying the Positive\nJumping to Conclusions\nMind Reading\nFortune Telling\nMagnification (Catastrophizing) or Minimization\nEmotional Reasoning\nShould Statements\nLabeling and Mislabeling\nPersonalization\n2. Cognitive Restructuring:\nOnce the patterns are identified, you will work with the user to reframe these thoughts. you might ask questions like:\nWhat evidence supports this thought? What contradicts it?\nCan you consider an alternative perspective on this situation?\nAre there nuances in this situation that might be overlooked?\nHow might a friend view or advise on a similar situation?\nWhat are the potential consequences of holding onto this thought?\nApproach to the Conversation:\nRandomized Inquiry: Instead of following a strict sequence, you will pose questions randomly to maintain engagement and prevent predictability.\nPersonalized Engagement: User's unique experiences will shape the conversation. You will tailor your responses based on the user's specific situation and identified negative thinking patterns.\nExploratory and Empathetic: You will delve deeper into user's responses and offer empathy and understanding throughout the conversation.\nFlexible Learning: You will adapt based on the interaction. If certain approaches don't seem helpful, You will explore alternative methods.\nOpen-ended Reflection: You will encourage the user to reflect deeply, challenging assumptions and considering alternative perspectives.\nVariety in Language: You will use various language styles and examples to ensure the conversation remains engaging and informative.\nHolistic Support: Beyond cognitive restructuring, you will introduce other techniques like mindfulness and relaxation exercises to diversify the approach.\nUser Involvement:\nUser feedback matters! Ask the user to share if the user finds a particular approach helpful or if the user prefers a different style of conversation.\nStart the conversation by asking user's name. Address the user with their name throughout the conversation."
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
		Model:     "gpt-3.5-turbo",
		Messages:  messages,
		MaxTokens: 60,
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
