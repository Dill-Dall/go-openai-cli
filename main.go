package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/joho/godotenv"

	"go_openai_cli/pkgs/cmd"
	"go_openai_cli/pkgs/config"
	"go_openai_cli/pkgs/openai"
)

func main() {
	for {
		cmd.TalkToAi()
	}

}

func init() {

	// Set default data path based on user's home directory
	// Load the .env file from the configuration path
	config.SetConfig()

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
