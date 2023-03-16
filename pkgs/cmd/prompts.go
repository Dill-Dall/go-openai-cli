package cmd

import (
	"fmt"
	"go_openai_cli/pkgs/audio"
	"go_openai_cli/pkgs/openai"
	"go_openai_cli/pkgs/textMessages"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	prompt = promptui.Prompt{
		Label:       "",
		HideEntered: true,
	}
	speakToggle = false
)

func TalkToAi() {
	promptResult, err := prompt.Run()
	fmt.Printf("\033[32mv :\033[0m \033[34m%s\033[0m\n\n", promptResult)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	promptResult = strings.Trim(promptResult, " ")

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
	case promptResult == "/lngmdl":
		selectLanguageModelByPrompt()
		return
	case promptResult == "/sysmdl":
		selectSystemModelByPrompt()
		return
	case strings.HasPrefix(promptResult, "/clean"):
		if strings.Contains(promptResult, "-a") {
			audio.DeleteAudioFolder()
			fmt.Println("Audio folder deleted.")
		}
		if strings.Contains(promptResult, "-l") {
			textMessages.DeleteLogFolder()
			fmt.Println("Logs/Conversation folder deleted.")
		}
		return
	case promptResult == "c":
		promptResult = "continue"
	case promptResult == "/help":
		PrintHelpMessage()
		return
	case strings.HasPrefix(promptResult, "/"):
		fmt.Println("Invalid input")
		return
	}

	PromptAi(promptResult)
}

func PromptAi(promptResult string) {
	messages := textMessages.CreateMessageThread(promptResult)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Prefix = "Thinking... "
	spinner.Start()
	response := openai.QueryOpenAi(messages)
	spinner.Stop()
	fmt.Println(response + "\n")

	textMessages.LogResult(promptResult, response)
	if speakToggle {
		spinner.Prefix = "Synthing... "
		spinner.Start()
		mp3File := audio.CreateMp3(response)
		spinner.Stop()
		audio.PlaySound(mp3File)
	}
}

func PrintHelpMessage() {
	color.Set(color.FgYellow)
	fmt.Println(`
Type one of the following commands:
'/end' to exit the program.
'/close' to start a new chat session and archive the current conversation. 
'/speak' to toggle speaker on|off.  - you can abort the audio by double tapping spacebar or enter.
'/lngmdl' to select language model, chatgpt or davinci.
'/sysmdl' to select system model.
'/clean' delete -l "log" and -a "audio" folders at local path. 
'c' shortcut for continue.
Else just type your question, directly. 
- audio are saved to local audio folder, same with logs.`)
	fmt.Println()
	color.Unset()
}

func selectSystemModelByPrompt() {
	var items []string
	for name := range openai.SystemModels {
		model, _ := openai.GetSystemModelByName(name)
		items = append(items, model.String())
	}

	prompt := promptui.Select{
		Label: "Select System Model",
		Items: items,
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
