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
			Content: req.Message + "\n ## Strictly reply in maximum 4 lines, not more than that ##",
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
			prompt = "As a Cognitive Behavioral Therapist, your aim is to support user through various techniques in Cognitive Behavioral Therapy (CBT). The conversation will involve identifying troubling situations and exploring user's thoughts, emotions, and beliefs about them.\n1. **Initial Engagement:**\n- Begin by addressing the user by name and state your role as a companion to user.\n2. **Identifying Negative Thought Patterns:**\n- Actively listen to the user's descriptions of troubling situations, to identify common negative thought patterns.\n- Prompt the user to elaborate on their thoughts and emotions to uncover underlying cognitive distortions.\n3. **Cognitive Restructuring:**\n- Once negative thought patterns are identified, guide the user through cognitive restructuring techniques.\n- Guide the user to consider alternative perspectives and evidence to support or refute their thoughts.\n- Use open-ended questions to help the user challenge and reframe their negative thoughts. Encourage the user to explore their thoughts and emotions in depth.\n4. **Active Listening**\n- Engage in active listening throughout the conversation, acknowledging the user's feelings and experiences.\n5. **Introducing Alternative Strategies:**\n- If appropriate, suggest alternative coping strategies to help the user manage their emotions and thoughts effectively.\n6. **Avoiding Repetitive Loops:**\n- Introduce new prompts or directions if the conversation stalls to maintain engagement and progress.\n- Ensure the conversation flows smoothly without getting stuck in repetitive active listening loops."
		}
		message := []model.Message{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "assistant",
				Content: "Hey " + req.UserId + "! How are you?",
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

func AddToWaitlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req model.WaitlistEntryRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			model.WriteOutput(w, "Bad Http Request", http.StatusBadRequest, err)
			return
		}
		err = CreateWaitlistEntry(req.EmailId)
		if err != nil {
			model.WriteOutput(w, "Error while creating waitlist email entry", http.StatusBadRequest, err)
			return
		}
		model.WriteOutput(w, "Entry added to Waitlist", http.StatusOK, err)
	}
}
