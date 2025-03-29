package app

import (
	"encoding/json"
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/utility"
)

type ChatHandler struct {
	ChatService *ChatService
	logger      utility.Logger
}

func NewChatHandler(service *ChatService, logger *utility.Logger) *ChatHandler {
	return &ChatHandler{
		ChatService: service,
		logger:      *logger,
	}
}

// Define the request structure to include system messages
type ChatRequest struct {
	SystemMessage string `json:"system_message"`
	UserMessage   string `json:"user_message"`
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		h.logger.Error("Failed to decode chat request: " + err.Error())
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	// Call ChatService to stream LLM response
	err := h.ChatService.StreamLLMResponse(req.SystemMessage, req.UserMessage, w)
	if err != nil {
		h.logger.Error("Failed to stream response: " + err.Error())
		http.Error(w, "Failed to process chat", http.StatusInternalServerError)
	}
}
