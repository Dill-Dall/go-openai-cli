package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"

	"go_openai_client/pkgs/audio"
	"go_openai_client/pkgs/openai"
	"go_openai_client/pkgs/textMessages"
)

var (
	prompt = promptui.Prompt{
		Label:       "Question",
		HideEntered: true,
	}
	speakToggle = false
)

func main() {
	godotenv.Load()
	openai.Init()
	fmt.Println("Welcome to the AI chatbot!")
	fmt.Println("Enter your questions below.")

	printHelpMessage()

	//openai.SelectSystemModel("DnDm")
	for true {
		talkToAi()
	}

	testing()
}

func testing() {
	newPrompt := "result"

	messages := textMessages.CreateMessageThread(newPrompt)
	res := openai.QueryOpenAi(messages)
	textMessages.LogResult(newPrompt, res)
}

func talkToAi() {
	promptResult, err := prompt.Run()
	fmt.Printf("\033[32m%s\033[0m Question: %s\n\n", "v", promptResult)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case promptResult == "/end":
		fmt.Println("Exiting...")
		os.Exit(1)

	case promptResult == "/close":
		messages := textMessages.CreateMessageThread("Can you create a title for all my questions? Max 5 words.")
		response := openai.QueryOpenAi(messages)
		fmt.Println("Closing thread...")
		textMessages.RotateLogFile(response)
		return
	case promptResult == "/speak":
		speakToggle = !speakToggle
		if speakToggle {
			fmt.Println("Speaker enabled.")
		} else {
			fmt.Println("Speaker disabled.")
		}
		return
	case promptResult == "/lngm":
		selectLanguageModelByPrompt()
		return
	case promptResult == "/sysmodel":
		selectSystemModelByPrompt()
		return
	case promptResult == "/gpt":
		fmt.Println("Switching to ChatGPT model...")
		openai.SetModel(openai.GPTModel)
		return
	case promptResult == "c":
		promptResult = "continue"
	case promptResult == "/help":
		printHelpMessage()
		return
	}

	messages := textMessages.CreateMessageThread(promptResult)
	response := openai.QueryOpenAi(messages)
	fmt.Println(response + "\n")
	textMessages.LogResult(promptResult, response)
	if speakToggle {
		mp3File := audio.CreateMp3(response)
		audio.PlaySound(mp3File)
	}
}

func printHelpMessage() {
	color.Set(color.FgYellow)
	fmt.Println("Type '/end' to exit the program.")
	fmt.Println("Type '/close' to start a new chat session and archive the current conversation.")
	fmt.Println("Type '/speak' to toggle speaker on.")
	fmt.Println("Type '/lngm' to select language model, chatgpt or davinci.")
	fmt.Println("Type '/sysmodel' to select system model.")
	fmt.Println("Type 'c' shortcut for continue.")
	color.Unset()
}

func selectSystemModelByPrompt() {
	prompt := promptui.Select{
		Label: "Select System Model",
		Items: []string{openai.AI.String(), openai.Detective.String(), openai.DnDm.String(), openai.Editor.String()},
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	err = openai.SelectSystemModel(result)
	if err != nil {
		fmt.Printf("Could not set model %v\n", err)
		return
	}
	fmt.Println("SystemModel set to: " + result)
	textMessages.SetSubfolder(result)
}

func selectLanguageModelByPrompt() {
	prompt := promptui.Select{
		Label: "Select Language Model",
		Items: []string{"chatgpt", "davinci"},
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	switch result {
	case "chatgpt":
		openai.SetModel(openai.GPTModel)
	case "davinci":
		openai.SetModel(openai.DavinciModel)
	default:
		fmt.Println("Could not set model " + result)
		return
	}
}
