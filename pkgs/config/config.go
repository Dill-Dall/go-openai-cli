package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey         string
	GoogleCredentials string
	DataPath          string
}

const (
	OpenAIKeyEnvVar         = "OPENAI_API_KEY"
	GoogleCredentialsEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"
)

var Cfg Config

func GetDataPath() string {
	return Cfg.DataPath
}

func SetConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	Cfg.DataPath = filepath.Join(homeDir, "Library/Application Support/go-openai-cli/data")
	if err := os.MkdirAll(Cfg.DataPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating data file directory: %v\n", err)
	}

	println("Data saves to: " + Cfg.DataPath)

	configPath := filepath.Join(homeDir, "Library/Application Support/go-openai-cli/.env")
	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf(`Error loading .env file from %s. Make sure your .env file exists at the path and is formatted like:
    %s=***
    #Absolute path to your data google credentials file
    %s=***
    `, configPath, OpenAIKeyEnvVar, GoogleCredentialsEnvVar)
	}

	Cfg.OpenAIKey = os.Getenv(OpenAIKeyEnvVar)
	if Cfg.OpenAIKey == "" {
		log.Fatalf("%s not found in .env file", OpenAIKeyEnvVar)
	}

	Cfg.GoogleCredentials = os.Getenv(GoogleCredentialsEnvVar)
	if Cfg.GoogleCredentials == "" {
		log.Fatalf("%s not found in .env file", GoogleCredentialsEnvVar)
	}
}
