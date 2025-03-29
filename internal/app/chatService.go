package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Sarvesh-10/ReadEazeBackend/config"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type ChatService struct {
	APIKey string
	logger utility.Logger
}

func NewChatService(apiKey string, logger *utility.Logger) *ChatService {
	return &ChatService{
		APIKey: apiKey,
		logger: *logger,
	}
}

func (s *ChatService) StreamLLMResponse(systemMessage, userMessage string, w http.ResponseWriter) error {
	payload := map[string]interface{}{
		"model": "llama-3.2-3b-preview",
		"messages": []map[string]string{
			{"role": "system", "content": systemMessage},
			{"role": "user", "content": userMessage},
		},
		"stream": true,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", config.AppConfig.LLMURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("Groq API error: " + string(body))
		return errors.New("failed to fetch response from LLM")
	}

	// Prepare SSE response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // Stop on error (including EOF)
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// Remove "data: " prefix
		line = strings.TrimPrefix(line, "data: ")

		// Handle [DONE] message
		if line == "[DONE]" {
			break
		}

		// Parse JSON
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			continue
		}

		// Extract "content" from response
		choices, ok := parsed["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			continue
		}

		delta, ok := choices[0].(map[string]interface{})["delta"].(map[string]interface{})
		if !ok {
			continue
		}

		content, ok := delta["content"].(string)
		if !ok || content == "" {
			continue
		}

		// Send only the content field as SSE
		_, writeErr := fmt.Fprintf(w, "%s\n\n", content)
		if writeErr != nil {
			return writeErr
		}

		// Flush to push data to client
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		}
	}

	return nil
}
