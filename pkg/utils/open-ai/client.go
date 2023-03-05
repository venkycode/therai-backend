package openai

import (
	"context"
	"fmt"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/venkycode/therai-backend/pkg/models"
)

const OpenAIGPT3_5Turbo = "gpt-3.5-turbo"

var ErrChatFinished = fmt.Errorf("ai decided to finish chat")

type OpenAIClient interface {
	ChatCompletion(conversation *models.Conversation) (string, error)
}

type OpenAIClientDeps interface {
	OpenAIApiKey() string
}

type openAIClient struct {
	deps   OpenAIClientDeps
	client gpt3.Client
}

func NewOpenAIClient(deps OpenAIClientDeps) OpenAIClient {
	client := gpt3.NewClient(deps.OpenAIApiKey(), gpt3.WithDefaultEngine(OpenAIGPT3_5Turbo))

	return &openAIClient{
		client: client,
		deps:   deps,
	}
}

func (o *openAIClient) ChatCompletion(conversation *models.Conversation) (string, error) {

	req := gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{},
	}

	for _, message := range conversation.ChronologicalMessages {
		req.Messages = append(req.Messages, gpt3.ChatCompletionRequestMessage{
			Role:    string(message.Role),
			Content: message.Text,
		})
	}

	res, err := o.client.ChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("cannot client.ChatCompletion: %w", err)
	}

	if res.Choices[0].Message.Role != string(models.Assistant) {
		panic("response cannot be of non assistant role")
	}

	return res.Choices[0].Message.Content, nil
}
