package api

type Conversation struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Messages []string `json:"messages"`
}

type PromptModel struct {
	Content     string `json:"content"`
	ID          string `json:"id"`
	SystemModel string `json:"system_model"`
}
