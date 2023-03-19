package textMessages

import (
	"fmt"
	"go_openai_cli/pkgs/api"
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

func SetSubfolder(subFolderName string) {
	subfolder = subFolderName
	aiLogFile = filepath.Join("logs", subFolderName, "AI.md")
}

func GetAiLogFile() string {
	return aiLogFile
}

func LogResult(conversationID, inputPrompt, response string) error {
	logDate := time.Now().Format("2006-01-02 15:04:05")
	logDir := filepath.Join(".", "logs", subfolder)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create log directory")
	}
	logFile := filepath.Join(logDir, fmt.Sprintf("%s.md", conversationID))

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

func DeleteLogFile(conversationID string) error {
	logDir := filepath.Dir(aiLogFile)
	logFile := filepath.Join(logDir, fmt.Sprintf("%s.md", conversationID))

	// check if the log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return fmt.Errorf("log file does not exist for conversation ID %s", conversationID)
	}

	// delete the log file
	if err := os.Remove(logFile); err != nil {
		return errors.Wrap(err, "failed to delete log file")
	}

	return nil
}

func RotateLogFile(conversationID, fileTitle string) error {
	logDir := filepath.Dir(aiLogFile)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to create log directory")
		}
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("%s.md", conversationID))
	logText, err := ioutil.ReadFile(logFile)
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

func CreateMessageThread(promptModel api.PromptModel) []gogpt.ChatCompletionMessage {
	messages := []gogpt.ChatCompletionMessage{}
	newMessage := gogpt.ChatCompletionMessage{Role: "user", Content: promptModel.Content}
	// create log file directory if it doesn't exist
	logDir := filepath.Dir(aiLogFile)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating log file directory: %v\n", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("%s.md", promptModel.ID))
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			fmt.Printf("Error creating log file: %v\n", err)
		}
		defer file.Close()
		file.WriteString(fmt.Sprintf("# Conversation ID: %s\n", promptModel.ID))
	}

	// read existing log messages from file
	previousMessages, err := ReadLogMessages(promptModel.ID)
	if err != nil {
		fmt.Printf("Error reading log messages: %v\n", err)
	}

	sysmodel, found := openai.GetSystemModelByName(promptModel.SystemModel)
	if !found {
		sysmodel = openai.GetSystemModel()
	}

	messages = append(messages, openai.SystemMessage(sysmodel))
	messages = append(messages, previousMessages...)
	messages = append(messages, newMessage)
	return messages
}

func ReadLogMessages(conversationID string) ([]gogpt.ChatCompletionMessage, error) {

	logDir := filepath.Dir(aiLogFile)
	logFile := filepath.Join(logDir, fmt.Sprintf("%s.md", conversationID))
	logText, err := ioutil.ReadFile(logFile)

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

func LoadConversations() ([]api.Conversation, error) {
	conversations := []api.Conversation{}

	logDir := filepath.Dir(aiLogFile)
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log directory")
	}

	identityRegex := regexp.MustCompile(`# Conversation ID: (.*)`)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(logDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read conversation file")
			}

			identityMatch := identityRegex.FindStringSubmatch(string(content))
			if len(identityMatch) > 1 {
				conversationID := identityMatch[1]
				messages, err := ReadLogMessages(conversationID)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read log messages for conversation")
				}

				conversationName := getLastTitleLine(messages)

				if len(conversationName) == 0 {
					conversationName = filepath.Base(file.Name())
				}

				conversations = append(conversations, api.Conversation{
					ID:       conversationID,
					Name:     conversationName,
					Messages: messagesToStrings(messages),
				})
			}
		}
	}

	return conversations, nil
}

func LoadConversation(conversationID string) (*api.Conversation, error) {
	logDir := filepath.Dir(aiLogFile)
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log directory")
	}

	identityRegex := regexp.MustCompile(`# Conversation ID: (.*)`)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(logDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read conversation file")
			}

			identityMatch := identityRegex.FindStringSubmatch(string(content))
			if len(identityMatch) > 1 && identityMatch[1] == conversationID {
				messages, err := ReadLogMessages(conversationID)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read log messages for conversation")
				}

				conversationName := getLastTitleLine(messages)

				if len(conversationName) == 0 {
					conversationName = filepath.Base(file.Name())
				}

				return &api.Conversation{
					ID:       conversationID,
					Name:     conversationName,
					Messages: messagesToStrings(messages),
				}, nil
			}
		}
	}

	return nil, errors.New("conversation not found")
}

func messagesToStrings(messages []gogpt.ChatCompletionMessage) []string {
	messageStrings := make([]string, len(messages))
	for i, msg := range messages {
		messageStrings[i] = fmt.Sprintf("%s: %s", strings.ToUpper(msg.Role), msg.Content)
	}
	return messageStrings
}

func getLastTitleLine(messages []gogpt.ChatCompletionMessage) string {
	var lastTitleLine string
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		lines := strings.Split(msg.Content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Title:") {
				lastTitleLine = strings.Split(line, "Title:")[1]
			}
		}
		if lastTitleLine != "" {
			break
		}
	}
	if lastTitleLine == "" {
		// Uses the first line as the title
		if len(messages[0].Content) > 60 {
			lastTitleLine = messages[0].Content[0:60] + "..."
		}
	}
	return lastTitleLine
}
