package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// Structures
type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

// --- OpenAI API request structure ---
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ChatHandler
func Chabgt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Call OpenAI API
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		http.Error(w, "missing OPENAI_API_KEY", http.StatusInternalServerError)
		return
	}

	openaiReq := openAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openAIMessage{
			{Role: "system", Content: "You are a helpful assistant who replies briefly and clearly."},
			{Role: "user", Content: req.Message},
		},
	}

	body, _ := json.Marshal(openaiReq)
	httpReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to contact OpenAI", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var openaiRes openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiRes); err != nil {
		http.Error(w, "invalid response from OpenAI", http.StatusInternalServerError)
		return
	}

	reply := "No response"
	if len(openaiRes.Choices) > 0 {
		reply = openaiRes.Choices[0].Message.Content
	}

	json.NewEncoder(w).Encode(ChatResponse{Reply: reply})
}
