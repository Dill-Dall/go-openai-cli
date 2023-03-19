package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/joho/godotenv"

	"go_openai_cli/pkgs/api"
	"go_openai_cli/pkgs/cmd"
	"go_openai_cli/pkgs/openai"
	"go_openai_cli/pkgs/textMessages"

	"encoding/json"
	"net/http"
	"sync"
)

var conversations []api.Conversation
var conversationMutex sync.Mutex

func main() {
	http.Handle("/api/conversations", CORS(http.HandlerFunc(handleConversations)))
	http.Handle("/api/systemodels", CORS(http.HandlerFunc(handleSystemModels)))

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)

	//for true {
	//	cmd.TalkToAi()
	//}
}

func handleSystemModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openai.GetSystemModels())
}

func handleConversations(w http.ResponseWriter, r *http.Request) {
	conversationMutex.Lock()
	defer conversationMutex.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		conversations, err := textMessages.LoadConversations()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(conversations)

	case http.MethodPost:
		var messageData api.PromptModel
		err := json.NewDecoder(r.Body).Decode(&messageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		convId := cmd.PromptAi(messageData)
		conversation, err := textMessages.LoadConversation(convId)
		if err != nil {
			http.Error(w, "Couldnt fix logMessages", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversation)

	case http.MethodDelete:
		var id string
		err := json.NewDecoder(r.Body).Decode(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = textMessages.DeleteLogFile(id)
		if err != nil {
			http.Error(w, "Could not delete conversation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
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
