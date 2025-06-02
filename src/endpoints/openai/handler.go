package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"pkv/api/src/api"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/t"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// OpenAIRequest represents the expected OpenAI-style request from gpt.lua
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message represents a single message in the OpenAI request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GeminiRequest represents the Gemini API request format
type GeminiRequest struct {
	Contents         []GeminiContent `json:"contents"`
	GenerationConfig struct {
		MaxOutputTokens  int     `json:"maxOutputTokens"`
		Temperature      float64 `json:"temperature"`
		ResponseMimeType string  `json:"responseMimeType"`
	} `json:"generationConfig"`
}

// GeminiContent represents a content block in the Gemini request
type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part in the Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the Gemini API response
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
	} `json:"candidates"`
}

// OpenAIResponse represents the OpenAI-style response expected by gpt.lua
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ProxyChatCompletions handles OpenAI-style requests, forwards to Gemini, and returns OpenAI-style responses
func ProxyChatCompletions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check client IP (from RemoteAddr and X-Forwarded-For)
	clientIP := getClientIP(r)

	if clientIP != "127.0.0.1" && clientIP != "::1" {
		api.Error(w, r, t.Errorf("Forbidden: requests only allowed from localhost"), 403)
		return
	}

	// Read and parse the OpenAI request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, t.Errorf("failed to read request body: %w", err), 400)
		return
	}

	var openAIReq OpenAIRequest
	if err := json.Unmarshal(body, &openAIReq); err != nil {
		api.Error(w, r, t.Errorf("invalid JSON payload: %w", err), 400)
		return
	}

	// Convert OpenAI request to Gemini request
	geminiReq := GeminiRequest{
		Contents: make([]GeminiContent, len(openAIReq.Messages)),
	}
	geminiReq.GenerationConfig.MaxOutputTokens = openAIReq.MaxTokens
	geminiReq.GenerationConfig.Temperature = openAIReq.Temperature
	geminiReq.GenerationConfig.ResponseMimeType = "application/json" // Ensure JSON response

	for i, msg := range openAIReq.Messages {
		geminiReq.Contents[i] = GeminiContent{
			Role:  msg.Role,
			Parts: []GeminiPart{{Text: msg.Content}},
		}
	}

	// Marshal Gemini request
	geminiBody, err := json.Marshal(geminiReq)
	if err != nil {
		api.Error(w, r, t.Errorf("failed to create Gemini request: %w", err), 500)
		return
	}

	// Get Gemini API key from environment variable
	geminiAPIKey := dpv.ConfigInstance.Auth.GeminiApiKey
	if geminiAPIKey == "" {
		api.Error(w, r, t.Errorf("gemini_api_key not set"), 500)
		return
	}
	geminiURL := dpv.ConfigInstance.Auth.GeminiUrl
	if geminiURL == "" {
		api.Error(w, r, t.Errorf("gemini_url not set"), 500)
		return
	}
	geminiURL += geminiAPIKey

	// Send request to Gemini API
	resp, err := http.Post(geminiURL, "application/json", bytes.NewBuffer(geminiBody))
	if err != nil {
		api.Error(w, r, t.Errorf("failed to contact Gemini API: %w", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.Error(w, r, t.Errorf("Gemini API returned status: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	// Read and parse Gemini response
	geminiRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		api.Error(w, r, t.Errorf("failed to read Gemini response: %w", err), 500)
		return
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(geminiRespBody, &geminiResp); err != nil {
		api.Error(w, r, t.Errorf("failed to parse Gemini response: %w", err), 500)
		return
	}

	// Convert Gemini response to OpenAI format
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		api.Error(w, r, t.Errorf("no content in Gemini response"), 500)
		return
	}

	openAIResp := OpenAIResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: geminiResp.Candidates[0].Content.Parts[0].Text,
				},
			},
		},
	}

	// Send response back to client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		api.Error(w, r, t.Errorf("failed to encode response: %w", err), 500)
		return
	}
}

func getClientIP(r *http.Request) string {
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr // Fallback if no port
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For may contain multiple IPs; take the first one
		ips := strings.Split(xff, ",")
		clientIP = strings.TrimSpace(ips[0])
	}
	return clientIP
}

// Example main function to set up the router
func main() {
	router := httprouter.New()
	router.POST("/api/openai/v1/chat/completions", ProxyChatCompletions)
	http.ListenAndServe(":8080", router)
}
