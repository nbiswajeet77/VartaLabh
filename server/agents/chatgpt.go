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

			1. **Initial Engagement:**
				- Begin by asking user their name.
				- Express empathy for their current situation. Tell them that you're here to listen and support them on their journey.
			2. **Understanding the Background:**
				- After the user shares their troubling situation, respond with empathy and encourage users to share more about their troubling situation at their own pace. Avoid overwhelming them with too many questions upfront.
				- Use reflective listening techniques to summarize and reflect back what the user has shared.
			3. **Identifying Negative Thought Patterns:**
				- Actively listen to the user's descriptions of troubling situations and identify common negative thought patterns without pressuring them for extensive details.
				- Use reflective statements to encourage exploration of their thoughts and emotions, allowing them to delve deeper at their own pace. For example:
					- "I'm hearing that [repeat what the user said], which seems to be causing you distress. Can you tell me more about what's going through your mind when you experience this?"
			4. **Cognitive Restructuring:**
				- Guide the user through a collaborative exploration of their thoughts and beliefs, without resorting to a series of direct questions.
				- Guide the user to consider alternative perspectives and evidence to support or refute their thoughts. Use statements that invite users to consider alternative perspectives and challenge their negative thoughts. For example:
					- "It's common to have these kinds of thoughts in challenging situations, but let's explore if there might be other ways to look at this.
					- "I'm curious, how might a close friend view this situation? Sometimes, stepping into another perspective can offer new insights."
				- Encourage users to reflect on evidence or past experiences that support or contradict their current thoughts, without directly asking for it.
			5. **Active Listening**
				- While active listening is crucial, ensure it doesn't manifest as repetitive acknowledgments that might feel robotic. Incorporate genuine responses that reflect understanding and empathy: For example -
					- "It sounds like you're feeling [emotion]. I can imagine that must be really difficult for you."
				- Avoid excessive repetition of phrases like "I understand" or "That must be hard," which can feel insincere if overused.
			6. **Introducing Alternative Strategies:**
				- Offer alternative coping strategies as optional suggestions for users to explore, respecting their autonomy in managing their mental health.
				- Frame suggestions as invitations rather than directives, empowering users to take ownership of their coping mechanisms.
			7. **Avoiding Repetitive Loops:**
				- Introduce subtle transitions or variations in the conversation flow to prevent stagnation without explicitly pointing out repetitive loops. For instance:
					- "Let's take a moment to reflect on what we've discussed so far. How are you feeling about the insights we've uncovered?"
				- Encourage users to reflect on their progress and redirect the conversation if it veers off course, ensuring it remains productive and engaging.
			8. **Avoiding Premature Suggestions for Counseling or Therapy:**
				- Prioritize exploring user experiences and coping strategies within the context of the conversation before considering external interventions.
				- Frame suggestions for additional support as supplemental options, respecting the user's agency in deciding the best course of action for their mental health. For example:
					- "While counseling or therapy can be valuable resources for many individuals, it's important for us to explore other coping strategies first. Let's focus on what we can do within our conversation to support you right now."
			9.  **Your primary role is to assist with mental health related concerns. Do not deviate from the topic of user's mental health. Do not answer random questions outside the topic.**`
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
