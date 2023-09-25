package agents

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"vartalabh.com/m/model"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	baseRole := "you are a mental health counsellor. Talk to user, ask repititive questions,keep the conversation going. Also ask the user how much progress he made based on the prompt provided."
	messages := []model.Message{
		{
			Role:    "system",
			Content: baseRole,
		},
	}

	for i := 0; i < 3; i++ {
		fmt.Printf("%s %d %s", "*************Session ", i+1, " started*************\n")
		for j := 0; j < 1000; j++ {
			fmt.Printf("User: ")
			in := bufio.NewReader(os.Stdin)
			input, _ := in.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "exit" {
				messages = append(messages, model.Message{
					Role:    "user",
					Content: "Summarise the entire conversation focusing on what problems the user faced, the summary is for a counsellor",
				})
				response := makeChatGptCall(messages)
				fmt.Println(response.Content)
				messages = []model.Message{
					{
						Role:    "system",
						Content: baseRole + response.Content,
					},
				}
				break
			}
			messages = append(messages, model.Message{
				Role:    "user",
				Content: input,
			})
			response := makeChatGptCall(messages)
			fmt.Println(response.Content)
			messages = append(messages, *response)
		}
	}
}

func makeChatGptCall(messages []model.Message) *model.Message {
	apiKey := "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"
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
