package gpt

import (
	"context"
	"fmt"

	"github.com/raspincel/es-ii-backend/utils"
	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func init() {
	apiKey := utils.GetEnv("OPENAI_API_KEY")
	client = openai.NewClient(apiKey)
}

func RequestGroups() string {
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

	return resp.Choices[0].Message.Content
}
