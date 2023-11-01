package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/joho/godotenv"

	"go_openai_cli/pkgs/cmd"
	"go_openai_cli/pkgs/openai"
)

func main() {
	//audio.RecordPrompt()

	for true {
		cmd.TalkToAi()
	}

}

func init() {
	godotenv.Load()
	openai.Init()

	color.Set(color.FgHiCyan)
	fmt.Println(`
╔════════════════════════════════════════════════════╗
║           Welcome to the Go Openai Client!         ║
║           a client tool made by Dill-Dall          ║
║                                                    ║
║  https://github.com/Dill-Dall/go-openai-cli        ║
╚════════════════════════════════════════════════════╝`)
	fmt.Println()
	color.Unset()
	cmd.PrintHelpMessage()
}
