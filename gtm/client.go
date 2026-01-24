// Package gtm provides a client for the Google Tag Manager API.
package gtm

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	tagmanager "google.golang.org/api/tagmanager/v2"
)

// Client wraps the Google Tag Manager API service.
type Client struct {
	Service                *tagmanager.Service
	RestrictedAccountID    string            // If set, restricts all operations to this account only
	RestrictedContainerIDs map[string]bool   // If set, restricts operations to these containers only
}

// NewClient creates a GTM client from an OAuth2 token source.
// The token source should handle automatic refresh.
func NewClient(ctx context.Context, tokenSource oauth2.TokenSource) (*Client, error) {
	return NewClientWithRestriction(ctx, tokenSource, "", nil)
}

// NewClientWithRestriction creates a GTM client with optional account and container restrictions.
// If restrictedAccountID is non-empty, all operations will be restricted to that account.
// If restrictedContainerIDs is non-empty, all operations will be restricted to those containers.
func NewClientWithRestriction(ctx context.Context, tokenSource oauth2.TokenSource, restrictedAccountID string, restrictedContainerIDs []string) (*Client, error) {
	if tokenSource == nil {
		return nil, fmt.Errorf("token source is required")
	}

	service, err := tagmanager.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to create tagmanager service: %w", err)
	}

	// Convert container IDs slice to map for O(1) lookup
	containerMap := make(map[string]bool)
	for _, id := range restrictedContainerIDs {
		containerMap[id] = true
	}

	return &Client{
		Service:                service,
		RestrictedAccountID:    restrictedAccountID,
		RestrictedContainerIDs: containerMap,
	}, nil
}

// ValidateAccountAccess checks if the given account ID is allowed.
// Returns an error if account restriction is enabled and the account doesn't match.
func (c *Client) ValidateAccountAccess(accountID string) error {
	if c.RestrictedAccountID != "" && accountID != c.RestrictedAccountID {
		return fmt.Errorf("access denied: account %s is not allowed (server is restricted to account %s)",
			accountID, c.RestrictedAccountID)
	}
	return nil
}

// ValidateContainerAccess checks if the given container ID is allowed.
// Returns an error if container restriction is enabled and the container is not in the allowed list.
func (c *Client) ValidateContainerAccess(containerID string) error {
	if len(c.RestrictedContainerIDs) > 0 && !c.RestrictedContainerIDs[containerID] {
		return fmt.Errorf("access denied: container %s is not allowed (server is restricted to specific containers)",
			containerID)
	}
	return nil
}

// ValidateAccess checks both account and container access in one call.
func (c *Client) ValidateAccess(accountID, containerID string) error {
	if err := c.ValidateAccountAccess(accountID); err != nil {
		return err
	}
	if err := c.ValidateContainerAccess(containerID); err != nil {
		return err
	}
	return nil
}
