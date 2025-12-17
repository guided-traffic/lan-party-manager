package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/yohcop/openid-go"
)

const (
	steamOpenIDEndpoint = "https://steamcommunity.com/openid"
)

// SteamAuth handles Steam OpenID authentication
type SteamAuth struct {
	callbackURL string
	nonceStore  openid.NonceStore
	discovery   openid.DiscoveryCache
}

// NewSteamAuth creates a new SteamAuth instance
func NewSteamAuth(backendURL string) *SteamAuth {
	return &SteamAuth{
		callbackURL: backendURL + "/api/v1/auth/steam/callback",
		nonceStore:  openid.NewSimpleNonceStore(),
		discovery:   openid.NewSimpleDiscoveryCache(),
	}
}

// GetAuthURL returns the Steam OpenID login URL
func (s *SteamAuth) GetAuthURL() (string, error) {
	log.Printf("[STEAM OPENID] Generating auth URL for callback: %s", s.callbackURL)
	authURL, err := openid.RedirectURL(
		steamOpenIDEndpoint,
		s.callbackURL,
		"",
	)
	if err != nil {
		log.Printf("[STEAM OPENID] ERROR - Failed to generate auth URL: %v", err)
		return "", err
	}
	log.Printf("[STEAM OPENID] OK - Auth URL generated, redirecting user to Steam")
	return authURL, nil
}

// ValidateCallback validates the OpenID callback and returns the Steam ID
func (s *SteamAuth) ValidateCallback(fullURL string) (string, error) {
	log.Printf("[STEAM OPENID] Validating callback from Steam...")

	// Verify the OpenID response
	id, err := openid.Verify(fullURL, s.discovery, s.nonceStore)
	if err != nil {
		log.Printf("[STEAM OPENID] ERROR - OpenID verification failed: %v", err)
		return "", fmt.Errorf("failed to verify OpenID response: %w", err)
	}

	// Extract Steam ID from the claimed ID
	// Format: https://steamcommunity.com/openid/id/76561198012345678
	steamID, err := extractSteamID(id)
	if err != nil {
		log.Printf("[STEAM OPENID] ERROR - Failed to extract Steam ID from: %s - %v", id, err)
		return "", err
	}

	log.Printf("[STEAM OPENID] OK - User authenticated successfully (Steam ID: %s)", steamID)
	return steamID, nil
}

// extractSteamID extracts the 64-bit Steam ID from the OpenID claimed identity
func extractSteamID(openIDIdentity string) (string, error) {
	// Steam ID regex pattern - 17 digit number
	pattern := regexp.MustCompile(`^https?://steamcommunity\.com/openid/id/(\d{17})$`)
	matches := pattern.FindStringSubmatch(openIDIdentity)

	if len(matches) != 2 {
		return "", fmt.Errorf("invalid Steam OpenID identity: %s", openIDIdentity)
	}

	return matches[1], nil
}

// BuildFullCallbackURL constructs the full callback URL from the request
func BuildFullCallbackURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	// Check for proxy headers (set by reverse proxy/ingress)
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	// Get the host - prefer X-Forwarded-Host for reverse proxy scenarios
	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	// Build the full URL including query parameters
	fullURL := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	return fullURL.String()
}

// ParseSteamID64 validates that a string is a valid Steam ID 64
func ParseSteamID64(steamID string) (string, error) {
	steamID = strings.TrimSpace(steamID)
	if len(steamID) != 17 {
		return "", fmt.Errorf("invalid Steam ID length: %d", len(steamID))
	}

	// Check if all characters are digits
	for _, c := range steamID {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid Steam ID format")
		}
	}

	return steamID, nil
}
