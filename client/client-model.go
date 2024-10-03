package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewClient(apiKey, apiBaseUrl string) *openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(apiBaseUrl),
	)
	return client
}
