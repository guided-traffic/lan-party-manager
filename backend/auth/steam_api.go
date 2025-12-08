package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	steamAPIBaseURL = "https://api.steampowered.com"
)

// SteamPlayer represents a Steam player's profile data
type SteamPlayer struct {
	SteamID         string `json:"steamid"`
	PersonaName     string `json:"personaname"`
	ProfileURL      string `json:"profileurl"`
	Avatar          string `json:"avatar"`           // 32x32
	AvatarMedium    string `json:"avatarmedium"`     // 64x64
	AvatarFull      string `json:"avatarfull"`       // 184x184
	PersonaState    int    `json:"personastate"`     // 0=Offline, 1=Online, etc.
	CommunityVisibilityState int `json:"communityvisibilitystate"`
	ProfileState    int    `json:"profilestate"`
	LastLogoff      int64  `json:"lastlogoff"`
	RealName        string `json:"realname,omitempty"`
	TimeCreated     int64  `json:"timecreated,omitempty"`
	LocCountryCode  string `json:"loccountrycode,omitempty"`
}

// steamAPIResponse represents the API response structure
type steamAPIResponse struct {
	Response struct {
		Players []SteamPlayer `json:"players"`
	} `json:"response"`
}

// SteamAPIClient handles communication with the Steam Web API
type SteamAPIClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewSteamAPIClient creates a new Steam API client
func NewSteamAPIClient(apiKey string) *SteamAPIClient {
	return &SteamAPIClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetPlayerSummary fetches a single player's profile data
func (c *SteamAPIClient) GetPlayerSummary(steamID string) (*SteamPlayer, error) {
	players, err := c.GetPlayerSummaries([]string{steamID})
	if err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return nil, fmt.Errorf("player not found: %s", steamID)
	}

	return &players[0], nil
}

// GetPlayerSummaries fetches profile data for multiple players (max 100)
func (c *SteamAPIClient) GetPlayerSummaries(steamIDs []string) ([]SteamPlayer, error) {
	if len(steamIDs) == 0 {
		return nil, fmt.Errorf("no Steam IDs provided")
	}

	if len(steamIDs) > 100 {
		return nil, fmt.Errorf("maximum 100 Steam IDs allowed per request")
	}

	if c.apiKey == "" {
		return nil, fmt.Errorf("Steam API key not configured")
	}

	// Build the API URL
	url := fmt.Sprintf(
		"%s/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
		steamAPIBaseURL,
		c.apiKey,
		strings.Join(steamIDs, ","),
	)

	// Make the request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call Steam API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var apiResp steamAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Steam API response: %w", err)
	}

	return apiResp.Response.Players, nil
}

// IsConfigured returns true if the API client has a valid API key
func (c *SteamAPIClient) IsConfigured() bool {
	return c.apiKey != ""
}
