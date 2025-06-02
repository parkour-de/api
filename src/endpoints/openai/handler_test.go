package openai

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"pkv/api/src/repository/dpv"
	tr "pkv/api/src/repository/t"
	"testing"
)

func TestProxyChatCompletions_Localhost(t *testing.T) {
	// Initialize config and set GeminiUrl to test server
	var err error
	config, err := dpv.NewConfig("../../../config.yml")
	if err != nil {
		t.Fatalf("could not initialise config instance: %v", err)
	}
	if err := tr.LoadDE(config); err != nil {
		log.Printf("Could not load strings_de.ini: %v", err)
	}
	// Start a test Gemini API server
	geminiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var geminiReq struct {
			Contents         []map[string]interface{} `json:"contents"`
			GenerationConfig struct {
				MaxOutputTokens  int     `json:"maxOutputTokens"`
				Temperature      float64 `json:"temperature"`
				ResponseMimeType string  `json:"responseMimeType"`
			} `json:"generationConfig"`
		}
		if err := json.NewDecoder(r.Body).Decode(&geminiReq); err != nil {
			http.Error(w, "invalid Gemini request", http.StatusBadRequest)
			return
		}
		if len(geminiReq.Contents) == 0 {
			http.Error(w, "missing contents", http.StatusBadRequest)
			return
		}
		if geminiReq.GenerationConfig.MaxOutputTokens == 0 {
			http.Error(w, "missing maxOutputTokens", http.StatusBadRequest)
			return
		}
		if geminiReq.GenerationConfig.ResponseMimeType != "application/json" {
			http.Error(w, "invalid responseMimeType", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"candidates": [{
				"content": {
					"parts": [{"text": "Hello from Gemini!"}],
					"role": "model"
				}
			}]
		}`))
	}))
	defer geminiSrv.Close()

	config.Auth.GeminiUrl = geminiSrv.URL + "?key="
	config.Auth.GeminiApiKey = "test-key"
	dpv.ConfigInstance = config

	// Prepare OpenAI-style request
	reqBody := map[string]interface{}{
		"model": "gemini-1.5-flash",
		"messages": []map[string]string{
			{"role": "user", "content": "Say hello"},
		},
		"max_tokens":  32,
		"temperature": 0.2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/openai/v1/chat/completions", bytes.NewReader(body))
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()

	ProxyChatCompletions(rr, req, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content != "Hello from Gemini!" {
		t.Errorf("unexpected content: %+v", resp)
	}
}
func TestProxyChatCompletions_NonLocalhost(t *testing.T) {
	// Initialize config and set GeminiUrl to test server
	var err error
	config, err := dpv.NewConfig("../../../config.yml")
	if err != nil {
		t.Fatalf("could not initialise config instance: %v", err)
	}
	if err := tr.LoadDE(config); err != nil {
		log.Printf("Could not load strings_de.ini: %v", err)
	}

	config.Auth.GeminiUrl = "invalid://invalid?key="
	config.Auth.GeminiApiKey = "test-key"
	dpv.ConfigInstance = config

	// Prepare OpenAI-style request
	reqBody := map[string]interface{}{
		"model": "gemini-1.5-flash",
		"messages": []map[string]string{
			{"role": "user", "content": "Say hello"},
		},
		"max_tokens":  32,
		"temperature": 0.2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/openai/v1/chat/completions", bytes.NewReader(body))
	req.RemoteAddr = "192.168.1.100:54321"
	rr := httptest.NewRecorder()

	ProxyChatCompletions(rr, req, nil)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rr.Code)
	}
}
