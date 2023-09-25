package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"vartalabh.com/m/agents"
	"vartalabh.com/m/model"
)

func main() {
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
				response := agents.MakeChatGptCall(messages)
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
			response := agents.MakeChatGptCall(messages)
			fmt.Println(response.Content)
			messages = append(messages, *response)
		}
	}
}
