package rest

import (
	"fmt"
	"sync"
)

// SSEClient represents a connected SSE client
type SSEClient struct {
	RepoID  string
	Channel chan string
}

// SSEManager manages all SSE connections
type SSEManager struct {
	clients map[string][]*SSEClient
	mu      sync.RWMutex
}

// NewSSEManager creates a new SSE manager
func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string][]*SSEClient),
	}
}

// AddClient adds a new SSE client for a repository
func (m *SSEManager) AddClient(repoID string, client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[repoID] = append(m.clients[repoID], client)
}

// RemoveClient removes an SSE client
func (m *SSEManager) RemoveClient(repoID string, client *SSEClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	clients := m.clients[repoID]
	for i, c := range clients {
		if c == client {
			m.clients[repoID] = append(clients[:i], clients[i+1:]...)
			close(client.Channel)
			break
		}
	}
	
	if len(m.clients[repoID]) == 0 {
		delete(m.clients, repoID)
	}
}

// NotifyClients sends a message to all clients waiting for a specific repository
func (m *SSEManager) NotifyClients(repoID string, message string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	clients := m.clients[repoID]
	for _, client := range clients {
		select {
		case client.Channel <- message:
		default:
			// Channel is full or closed, skip
		}
	}
}

// GetClientCount returns the number of clients for a repository
func (m *SSEManager) GetClientCount(repoID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients[repoID])
}

// Global SSE manager instance
var GlobalSSEManager = NewSSEManager()

// FormatSSEMessage formats a message for SSE
func FormatSSEMessage(data string) string {
	return fmt.Sprintf("data: %s\n\n", data)
}
