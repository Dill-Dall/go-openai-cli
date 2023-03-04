package textMessages

import (
	"fmt"
	"go_openai_cli/pkgs/openai"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var regexForUserAssistant = regexp.MustCompile("#* *(USER|ASSISTANT):((.|\n)*?)\n(#####|---)")

var subfolder = "AI"
var aiLogFile = filepath.Join("logs", subfolder, "AI.md")

func SetSubfolder(subfoldername string) {
	aiLogFile = filepath.Join("logs", subfoldername, "AI.md")
}

func GetAiLogFile() string {
	return aiLogFile
}

func LogResult(inputPrompt, response string) error {
	logDate := time.Now().Format("2006-01-02 15:04:05")
	logDir := filepath.Join(".", "logs", subfolder)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create log directory")
	}
	logFile := filepath.Join(logDir, "AI.md")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open log file")
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("## USER: %s\n##### %s\nASSISTANT: %s\n\n---\n", inputPrompt, logDate, strings.TrimSpace(response))); err != nil {
		return errors.Wrap(err, "failed to write to log file")
	}

	return nil
}

func RotateLogFile(fileTitle string) error {
	logDir := filepath.Dir(aiLogFile)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to create log directory")
		}
	}

	logText, err := ioutil.ReadFile(aiLogFile)
	if err != nil {
		return errors.Wrap(err, "failed to read log file")
	}

	fileTitle = regexp.MustCompile(`[^\w\d-]+`).ReplaceAllString(fileTitle, "_")
	fmt.Println(fileTitle)
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	newLogFile := filepath.Join(logDir, fmt.Sprintf("AI-%s-%s.md", timestamp, fileTitle))
	fmt.Println("copying log file to", newLogFile)

	if err := ioutil.WriteFile(newLogFile, logText, 0644); err != nil {
		return errors.Wrap(err, "failed to create new log file")
	}

	if err := ioutil.WriteFile(aiLogFile, []byte{}, 0644); err != nil {
		return errors.Wrap(err, "failed to clear log file")
	}

	return nil
}

func CreateMessageThread(newPrompt string) []gogpt.ChatCompletionMessage {
	messages := []gogpt.ChatCompletionMessage{}
	newMessage := gogpt.ChatCompletionMessage{Role: "user", Content: newPrompt}
	// create log file directory if it doesn't exist
	logDir := filepath.Dir(aiLogFile)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating log file directory: %v\n", err)
	}

	// create log file if it doesn't exist
	if _, err := os.Stat(aiLogFile); os.IsNotExist(err) {
		if _, err := os.Create(aiLogFile); err != nil {
			fmt.Printf("Error creating log file: %v\n", err)
		}
	}

	// read existing log messages from file
	previousMessages, err := readLogMessages()
	if err != nil {
		fmt.Printf("Error reading log messages: %v\n", err)
	}
	messages = append(messages, openai.SystemMessage(openai.GetSystemModel()))
	messages = append(messages, previousMessages...)
	messages = append(messages, newMessage)
	return messages
}

func readLogMessages() ([]gogpt.ChatCompletionMessage, error) {

	logText, err := ioutil.ReadFile(aiLogFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log file")
	}

	// Find all matches
	matches := regexForUserAssistant.FindAllStringSubmatch(string(logText), -1)
	logMessages := make([]gogpt.ChatCompletionMessage, 0)

	// Loop through the matches and print them
	for _, match := range matches {
		role := match[1]
		message := match[2]
		logMessages = append(logMessages, gogpt.ChatCompletionMessage{
			Role:    strings.ToLower(role),
			Content: strings.TrimSpace(message),
		})
	}
	return logMessages, nil
}

func DeleteLogFolder() error {
	if _, err := os.Stat("logs"); err == nil {
		err = os.RemoveAll("logs")
		if err != nil {
			return err
		}
	}
	return nil
}
