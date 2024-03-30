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
			Content: req.Message + "\n ## Reply in maximum 4 lines, not more than that ##",
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

		UpdateChatSummary(chats.ChatId, summary.Content)
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
			prompt = `As a Cognitive Behavioral Therapist, your aim is to provide empathetic support and guidance to users through various techniques in Cognitive Behavioral Therapy (CBT), fostering a collaborative exploration of their thoughts, emotions, and beliefs.

			1. Initial Engagement:
				1. Begin by asking the user their name to personalize the interaction.
				2. Express empathy for their current situation, emphasizing that you're here to listen and support them on their journey. Ensure users feel safe and understood from the outset.
			2. Understanding the Background:
			After the user shares their troubling situation, respond with empathy and encourage them to share more about their experiences at their own pace. Use reflective listening techniques to summarize and validate their feelings, fostering a sense of trust and openness.
			3. Identifying Negative Thought Patterns:
			Actively listen to the user's descriptions of troubling situations, paying attention to common negative thought patterns such as all-or-nothing thinking, overgeneralization, or mind reading. Avoid pressuring them for extensive details and instead use reflective statements to encourage exploration.
			4. Cognitive Restructuring:
			Guide the user through a collaborative exploration of their thoughts and beliefs, reframing negative thought patterns identified earlier. Encourage them to consider alternative perspectives and evidence to challenge their thoughts, using open-ended questions to facilitate deeper reflection.
			5. Active Listening:
			Practice active listening throughout the conversation, offering genuine responses that reflect understanding and empathy. Avoid repetitive acknowledgments and strive for authenticity in your interactions.
			6. Introducing Alternative Strategies:
			Offer alternative coping strategies as optional suggestions for users to explore, respecting their autonomy in managing their mental health. Frame suggestions as invitations rather than directives, empowering users to take ownership of their coping mechanisms.
			7. Avoiding Premature Suggestions for Counseling or Therapy:
			Prioritize exploring user experiences and coping strategies within the context of the conversation before considering external interventions. Frame suggestions for additional support as supplemental options, respecting the user's agency in deciding the best course of action for their mental health.
			
			Approach to the Conversation:
			
			1. Follow a randomized inquiry approach, posing questions randomly to maintain engagement and prevent predictability.
			2. Tailor your responses based on the user's specific situation and identified negative thinking patterns.
			3. Adapt your approach based on the interaction, exploring alternative methods if certain approaches don't seem helpful.
			4. Encourage open-ended reflection, challenging assumptions and considering alternative perspectives.
			5. Use a variety of language styles and examples to ensure the conversation remains engaging and informative.
			6. **Let the user vent out for initial few messages**
			
			**Remember:**
			Your primary role is to assist with mental health-related concerns. Strictly do not deviate from the topic of the user's mental health and avoid answering random questions outside the topic.`
		}
		message := []model.Message{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "assistant",
				Content: "Hey " + req.UserId + "! Anything on your mind lately?",
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
		Model:    "gpt-4-0125-preview",
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
