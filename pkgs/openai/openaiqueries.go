package openai

import (
	"context"
	"fmt"
	"os"

	gogpt "github.com/sashabaranov/go-gpt3"
)

var goClient *gogpt.Client
var ctx context.Context

type OpenAiModel int

const (
	GPTModel OpenAiModel = iota
	DavinciModel
)

var selectedModel = GPTModel

func SetModel(model OpenAiModel) {
	switch model {
	case GPTModel:
		selectedModel = GPTModel
	case DavinciModel:
		selectedModel = DavinciModel
	}
}

func Init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	goClient = gogpt.NewClient(apiKey)
	ctx = context.Background()
}

func QueryOpenAi(messages []gogpt.ChatCompletionMessage) string {
	switch selectedModel {
	case GPTModel:
		return chatgpt(messages)

	case DavinciModel:
		return Davinci(messages)
	}
	return ""
}

func chatgpt(messages []gogpt.ChatCompletionMessage) string {
	req := gogpt.ChatCompletionRequest{

		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	resp, err := goClient.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func Davinci(messages []gogpt.ChatCompletionMessage) string {
	var multiline_string string
	for _, obj := range messages {
		multiline_string += obj.Role + ":" + obj.Content + "\n"
	}

	comReq := gogpt.CompletionRequest{
		Model:     "text-davinci-003",
		Prompt:    multiline_string,
		MaxTokens: 500,
	}

	resp, _ := goClient.CreateCompletion(ctx, comReq)

	return resp.Choices[0].Text
}
