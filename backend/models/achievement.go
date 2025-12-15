package models

// Achievement represents a predefined achievement that users can vote for
type Achievement struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsPositive  bool   `json:"is_positive"`
}

// Predefined achievements
var Achievements = map[string]Achievement{
	// Positive achievements
	"pro-player": {
		ID:          "pro-player",
		Name:        "Pro Player",
		Description: "Zeigt herausragende Fähigkeiten, für seine Verhältnisse.",
		ImageURL:    "/icons/achievements/trophy.svg",
		IsPositive:  true,
	},
	"teamplayer": {
		ID:          "teamplayer",
		Name:        "Teamplayer",
		Description: "Stirbt freiwillig zuerst, damit du looten kannst.",
		ImageURL:    "/icons/achievements/three-friends.svg",
		IsPositive:  true,
	},
	"clutch-king": {
		ID:          "clutch-king",
		Name:        "Clutch King",
		Description: "1v5? Kein Problem. Wo ist die Herausforderung?",
		ImageURL:    "/icons/achievements/muscle-up.svg",
		IsPositive:  true,
	},
	"support-hero": {
		ID:          "support-hero",
		Name:        "Support Hero",
		Description: "Flasht die Gegner, nicht das eigene Team. Ein Wunder!",
		ImageURL:    "/icons/achievements/shaking-hands.svg",
		IsPositive:  true,
	},
	"stratege": {
		ID:          "stratege",
		Name:        "Stratege",
		Description: "Seine Taktik: 'Vertraut mir, Jungs!' alle sterben",
		ImageURL:    "/icons/achievements/chess-king.svg",
		IsPositive:  true,
	},
	"good-sport": {
		ID:          "good-sport",
		Name:        "Gute Manieren",
		Description: "Der einzige der nach dem Match noch Freunde hat.",
		ImageURL:    "/icons/achievements/bow-tie-ribbon.svg",
		IsPositive:  true,
	},

	// Negative achievements
	"rage-quitter": {
		ID:          "rage-quitter",
		Name:        "Rage Quitter",
		Description: "'Das Spiel ist eh buggy' – 0.3 Sekunden nach dem Tod.",
		ImageURL:    "/icons/achievements/enrage.svg",
		IsPositive:  false,
	},
	"toxic": {
		ID:          "toxic",
		Name:        "Toxic",
		Description: "Caps Lock ist sein Standardmodus.",
		ImageURL:    "/icons/achievements/death-juice.svg",
		IsPositive:  false,
	},
	"friendly-fire-expert": {
		ID:          "friendly-fire-expert",
		Name:        "Friendly Fire Expert",
		Description: "Sein Team fürchtet ihn mehr als die Gegner.",
		ImageURL:    "/icons/achievements/backstab.svg",
		IsPositive:  false,
	},
}

// GetAllAchievements returns all achievements as a slice
func GetAllAchievements() []Achievement {
	achievements := make([]Achievement, 0, len(Achievements))
	for _, a := range Achievements {
		achievements = append(achievements, a)
	}
	return achievements
}

// GetAchievement returns an achievement by ID
func GetAchievement(id string) (Achievement, bool) {
	a, ok := Achievements[id]
	return a, ok
}

// IsValidAchievement checks if an achievement ID is valid
func IsValidAchievement(id string) bool {
	_, ok := Achievements[id]
	return ok
}
