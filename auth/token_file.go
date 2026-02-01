package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	tagmanager "google.golang.org/api/tagmanager/v2"
)

// TokenFile represents the structure of a saved OAuth token file.
type TokenFile struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
}

// LoadTokenFromFile loads an OAuth token from a JSON file and returns an auto-refreshing token source.
func LoadTokenFromFile(path, clientID, clientSecret string) (oauth2.TokenSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var tokenFile TokenFile
	if err := json.Unmarshal(data, &tokenFile); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  tokenFile.AccessToken,
		RefreshToken: tokenFile.RefreshToken,
		TokenType:    tokenFile.TokenType,
		Expiry:       tokenFile.Expiry,
	}

	// Create OAuth2 config for token refresh
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			tagmanager.TagmanagerDeleteContainersScope,
			tagmanager.TagmanagerEditContainersScope,
			tagmanager.TagmanagerEditContainerversionsScope,
			tagmanager.TagmanagerPublishScope,
		},
	}

	// Create auto-refreshing token source
	tokenSource := config.TokenSource(context.Background(), token)

	return tokenSource, nil
}

// SaveTokenToFile saves an OAuth token to a JSON file.
func SaveTokenToFile(path string, token *oauth2.Token) error {
	tokenFile := TokenFile{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	data, err := json.MarshalIndent(tokenFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Ensure parent directory exists
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create token directory: %w", err)
		}
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
