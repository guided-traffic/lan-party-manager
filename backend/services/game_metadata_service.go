package services

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// GameMetadata contains manually curated metadata for a game
type GameMetadata struct {
	MaxPlayers int    `json:"max_players"`
	Notes      string `json:"notes,omitempty"`
}

// GameMetadataService manages game metadata loaded from a JSON file
type GameMetadataService struct {
	metadata map[string]*GameMetadata // appID (string) -> metadata
	mu       sync.RWMutex
	filePath string
}

// NewGameMetadataService creates a new game metadata service
// filePath is the path to the game_metadata.json file
func NewGameMetadataService(filePath string) *GameMetadataService {
	service := &GameMetadataService{
		metadata: make(map[string]*GameMetadata),
		filePath: filePath,
	}
	service.loadMetadata()
	return service
}

// loadMetadata loads the metadata from the JSON file
func (s *GameMetadataService) loadMetadata() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Game metadata file not found at %s, using empty metadata", s.filePath)
		} else {
			log.Printf("Error reading game metadata file: %v", err)
		}
		return
	}

	var metadata map[string]*GameMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		log.Printf("Error parsing game metadata JSON: %v", err)
		return
	}

	s.metadata = metadata
	log.Printf("Loaded game metadata for %d games", len(metadata))
}

// GetMetadata returns the metadata for a given app ID
func (s *GameMetadataService) GetMetadata(appID int) *GameMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert appID to string for map lookup
	appIDStr := intToString(appID)
	return s.metadata[appIDStr]
}

// GetMaxPlayers returns the max players for a game, or 0 if not known
func (s *GameMetadataService) GetMaxPlayers(appID int) int {
	meta := s.GetMetadata(appID)
	if meta == nil {
		return 0
	}
	return meta.MaxPlayers
}

// Reload reloads the metadata from disk
func (s *GameMetadataService) Reload() {
	s.loadMetadata()
}

// intToString converts an int to string without importing strconv
func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	// Max int64 has 19 digits
	digits := make([]byte, 20)
	i := len(digits)

	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		digits[i] = '-'
	}

	return string(digits[i:])
}
