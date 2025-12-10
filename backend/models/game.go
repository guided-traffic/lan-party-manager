package models

// Game represents a Steam game with multiplayer information
type Game struct {
	AppID           int      `json:"app_id"`
	Name            string   `json:"name"`
	HeaderImageURL  string   `json:"header_image_url"`  // 460x215
	CapsuleImageURL string   `json:"capsule_image_url"` // 231x87
	PlaytimeForever int      `json:"playtime_forever"`  // Total playtime in minutes
	Categories      []string `json:"categories"`        // e.g., "Multi-player", "Co-op", etc.
	OwnerCount      int      `json:"owner_count"`       // Number of players who own this game
	Owners          []string `json:"owners"`            // Steam IDs of owners
	IsPinned        bool     `json:"is_pinned"`         // Whether this game is pinned/featured
	// Price information
	IsFree          bool   `json:"is_free"`           // True if free-to-play
	PriceCents      int    `json:"price_cents"`       // Current price in cents (e.g., 5999 = 59.99€)
	OriginalCents   int    `json:"original_cents"`    // Original price before discount
	DiscountPercent int    `json:"discount_percent"`  // Discount percentage (0-100)
	PriceFormatted  string `json:"price_formatted"`   // Formatted price string (e.g., "59,99€" or "Free")
}

// GameOwnership represents a player's ownership of a game
type GameOwnership struct {
	SteamID         string `json:"steam_id"`
	AppID           int    `json:"app_id"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
	IconURL         string `json:"icon_url"`
}

// GamesResponse represents the API response for games
type GamesResponse struct {
	PinnedGames []Game `json:"pinned_games"`
	AllGames    []Game `json:"all_games"`
}

// MultiplayerCategories defines which Steam categories indicate multiplayer capability
var MultiplayerCategories = []string{
	"Multi-player",
	"Co-op",
	"Online Co-op",
	"LAN Co-op",
	"LAN PvP",
}

// IsMultiplayerCategory checks if a category indicates multiplayer
func IsMultiplayerCategory(category string) bool {
	for _, mp := range MultiplayerCategories {
		if mp == category {
			return true
		}
	}
	return false
}

// HasMultiplayerCategory checks if a game has any multiplayer category
func (g *Game) HasMultiplayerCategory() bool {
	for _, cat := range g.Categories {
		if IsMultiplayerCategory(cat) {
			return true
		}
	}
	return false
}
