package openai

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var goClient *openai.Client
var ctx context.Context

type OpenAiModel int

const (
	GPTModel OpenAiModel = iota
	DavinciModel
	GPT4Model
)

var selectedModel = GPTModel

func SetModel(model OpenAiModel) {
	switch model {
	case GPTModel:
		selectedModel = GPTModel
	case DavinciModel:
		selectedModel = DavinciModel
	case GPT4Model:
		print("selected GPT4Model")
		selectedModel = GPT4Model
	}
}

func Init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	goClient = openai.NewClient(apiKey)
	ctx = context.Background()
}

func QueryOpenAi(messages []openai.ChatCompletionMessage) string {
	switch selectedModel {
	case GPTModel:
		return chatgpt(messages)

	case DavinciModel:
		return Davinci(messages)

	case GPT4Model:
		return chatgpt4(messages)
	}
	return ""
}

func chatgpt(messages []openai.ChatCompletionMessage) string {
	req := openai.ChatCompletionRequest{

		Model:    openai.GPT3Dot5Turbo1106,
		Messages: messages,
	}

	resp, err := goClient.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func chatgpt4(messages []openai.ChatCompletionMessage) string {
	req := openai.ChatCompletionRequest{

		Model:    openai.GPT4TurboPreview,
		Messages: messages,
	}

	resp, err := goClient.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func Davinci(messages []openai.ChatCompletionMessage) string {
	var multiline_string string
	for _, obj := range messages {
		multiline_string += obj.Role + ":" + obj.Content + "\n"
	}

	comReq := openai.CompletionRequest{
		Model:     openai.GPT3Davinci,
		Prompt:    multiline_string,
		MaxTokens: 500,
	}

	resp, _ := goClient.CreateCompletion(ctx, comReq)

	return resp.Choices[0].Text
}

func Dalle(prompt string) string {

	imageRequest := openai.ImageRequest{
		Prompt: prompt,
		N:      1,
		Size:   "512x512",
	}
	resp, _ := goClient.CreateImage(ctx, imageRequest)

	return resp.Data[0].URL

}

func Dalle3(prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         prompt,
			Size:           openai.CreateImageSize1792x1024,
			Model:          openai.CreateImageModelDallE3,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
		},
	)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return ""
	}

	b, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return ""
	}
	fileName := writeFile(prompt, b)
	return fmt.Sprintf("file saved to %s", fileName)
}

func getNextFileName(basename string, dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	max := 0
	for _, file := range files {
		name := file.Name()
		ext := filepath.Ext(name)
		base := name[:len(name)-len(ext)]

		parts := strings.SplitN(base, "-", 2)
		if len(parts) > 1 {
			num, err := strconv.Atoi(parts[0])
			if err == nil && num > max {
				max = num
			}
		}
	}

	return fmt.Sprintf("%s/%d-%s.png", dir, max+1, basename), nil
}

func writeFile(filename string, b []byte) string {
	fileName, err := getNextFileName(filename, "dalle")
	if err != nil {
		fmt.Printf("Error getting next file name: %v\n", err)
		return ""
	}

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return ""
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		fmt.Printf("File write error: %v\n", err)
		return ""
	}

	return f.Name()
}
