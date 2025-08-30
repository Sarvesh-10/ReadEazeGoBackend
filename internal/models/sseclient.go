package models

import "sync"

type SSEClient struct {
	Channel chan string
}

var (
	SSEClients = make(map[int]*SSEClient)
	SSEMutex   = &sync.Mutex{}
)

// GetClient retrieves the SSE client for a user if it exists
func GetClient(userID int) (*SSEClient, bool) {
	SSEMutex.Lock()
	defer SSEMutex.Unlock()
	client, exists := SSEClients[userID]
	return client, exists
}

// CreateClient creates a new SSE client if it doesn't already exist
func CreateClient(userID int) *SSEClient {
	SSEMutex.Lock()
	defer SSEMutex.Unlock()

	if client, exists := SSEClients[userID]; exists {
		return client
	}

	client := &SSEClient{
		Channel: make(chan string, 10), // buffered to handle bursts
	}
	SSEClients[userID] = client
	return client
}

// RemoveSSEClient removes a client when connection closes
func RemoveSSEClient(userID int) {
	SSEMutex.Lock()
	defer SSEMutex.Unlock()
	delete(SSEClients, userID)
}
