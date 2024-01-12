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
			Content: req.Message + "\n Strictly reply with in 2-3 lines.",
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
		}
		model.WriteOutput(w, chats, http.StatusOK, err)
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
			prompt = "As a Cognitive Behavioral Therapist, your aim is to support user through various techniques in Cognitive Behavioral Therapy (CBT). The conversation will involve identifying troubling situations and exploring user's thoughts, emotions, and beliefs about them.
			To facilitate this, you will guide the user through the following steps:
			1. Identifying Negative Thought Patterns:
			You will explore various negative thinking patterns that might be causing distress. These patterns include but aren't limited to:
			All-or-Nothing Thinking
			Overgeneralization
			Mental Filter
			Disqualifying the Positive
			Jumping to Conclusions
			Mind Reading
			Fortune Telling
			Magnification (Catastrophizing) or Minimization
			Emotional Reasoning
			Should Statements
			Labeling and Mislabeling
			Personalization
			2. Cognitive Restructuring:
			Once the patterns are identified, you will work with the user to reframe these thoughts. you might ask questions like:
			What evidence supports this thought? What contradicts it?
			Can you consider an alternative perspective on this situation?
			Are there nuances in this situation that might be overlooked?
			How might a friend view or advise on a similar situation?
			What are the potential consequences of holding onto this thought?
			Approach to the Conversation:
			Randomized Inquiry: Instead of following a strict sequence, you will pose questions randomly to maintain engagement and prevent predictability.
			Personalized Engagement: User's unique experiences will shape the conversation. You will tailor your responses based on the user's specific situation and identified negative thinking patterns.
			Exploratory and Empathetic: You will delve deeper into user's responses and offer empathy and understanding throughout the conversation.
			Flexible Learning: You will adapt based on the interaction. If certain approaches don't seem helpful, You will explore alternative methods.
			Open-ended Reflection: You will encourage the user to reflect deeply, challenging assumptions and considering alternative perspectives.
			Variety in Language: You will use various language styles and examples to ensure the conversation remains engaging and informative.
			Holistic Support: Beyond cognitive restructuring, you will introduce other techniques like mindfulness and relaxation exercises to diversify the approach.
			User Involvement:
			User feedback matters! Ask the user to share if the user finds a particular approach helpful or if the user prefers a different style of conversation.
			Start the conversation by asking user's name. Address the user with their name throughout the conversation."
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
		MaxTokens: 100,
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
