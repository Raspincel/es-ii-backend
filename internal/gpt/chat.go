package gpt

import (
	"context"
	"fmt"
	"strings"

	"github.com/raspincel/es-ii-backend/internal/utils"
	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func init() {
	apiKey := utils.GetEnv("OPENAI_API_KEY")
	client = openai.NewClient(apiKey)
}

func RequestGroups() string {
	fmt.Println("Requesting groups", attempts)

	if attempts == uint16(1000) {
		fmt.Println("Uh-oh")
		return ""
	}

	jobId := utils.GetEnv("JOB_ID")

	ctx := context.Background()

	fineTuningJob, err := client.RetrieveFineTuningJob(ctx, jobId)

	if err != nil {
		fmt.Printf("Getting fine tune model error: %v\n", err)
		return ""
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: fineTuningJob.FineTunedModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: utils.GetEnv("PROMPT"),
			},
		},
	})

	if err != nil {
		fmt.Printf("Create completion error %v\n", err)
		return ""
	}

	answer := resp.Choices[0].Message.Content
	verification := verifyGroups(client, answer)

	if !verification {
		attempts++
		answer = RequestGroups()
	}

	if answer == "" {
		fmt.Println("Failed to get groups")
		return ""
	}

	attempts = 0
	return answer
}

var attempts uint16 = 0

func verifyGroups(client *openai.Client, content string) bool {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`Responda apenas true ou false. A seguir, irei lhe enviar uma coleção de quatro grupos temáticos, e palavras relacionadas a cada grupo. Você deve analisar cada grupo individualmente e retornar true se: 
					1- Todas as palavras seguem a gramática correta da língua portuguesa
					2- Todas as quatro palavras de um determinado grupo estão relacionadas com o tema do grupo
					Caso contrário, retorne false
					
					%s`, content),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Create completion error %v\n", err)
		return false
	}

	lowerCaseResponse := strings.ToLower(resp.Choices[0].Message.Content)

	return strings.Contains(lowerCaseResponse, "true")
}
