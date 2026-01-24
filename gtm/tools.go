package gtm

import (
	"context"
	"fmt"

	"gtm-mcp-server/auth"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
)

// restrictedAccountID is the account ID that all operations are restricted to.
// If empty, no restriction is applied.
var restrictedAccountID string

// restrictedContainerIDs is the list of container IDs that operations are restricted to.
// If empty, no container restriction is applied.
var restrictedContainerIDs []string

// stdioTokenSource is the OAuth2 token source used in stdio mode.
var stdioTokenSource oauth2.TokenSource

// SetRestrictedAccountID configures the package to restrict all operations to a specific account.
// Pass an empty string to disable restriction.
func SetRestrictedAccountID(accountID string) {
	restrictedAccountID = accountID
}

// GetRestrictedAccountID returns the currently configured account restriction.
func GetRestrictedAccountID() string {
	return restrictedAccountID
}

// SetRestrictedContainerIDs configures the package to restrict all operations to specific containers.
// Pass nil or empty slice to disable restriction.
func SetRestrictedContainerIDs(containerIDs []string) {
	restrictedContainerIDs = containerIDs
}

// GetRestrictedContainerIDs returns the currently configured container restrictions.
func GetRestrictedContainerIDs() []string {
	return restrictedContainerIDs
}

// SetStdioTokenSource sets the OAuth2 token source for stdio mode.
func SetStdioTokenSource(ts oauth2.TokenSource) {
	stdioTokenSource = ts
}

// HasStdioTokenSource returns true if a stdio token source is configured.
func HasStdioTokenSource() bool {
	return stdioTokenSource != nil
}

// RegisterTools adds all GTM tools to the MCP server.
func RegisterTools(server *mcp.Server) {
	// Read operations
	registerListAccounts(server)
	registerListContainers(server)
	registerListWorkspaces(server)
	registerListTags(server)
	registerGetTag(server)
	registerListTriggers(server)
	registerListVariables(server)
	registerGetVariable(server)
	registerListFolders(server)
	registerGetFolderEntities(server)
	registerListTemplates(server)

	// Write operations
	registerCreateTag(server)
	registerUpdateTag(server)
	registerDeleteTag(server)
	registerCreateTrigger(server)
	registerUpdateTrigger(server)
	registerDeleteTrigger(server)
	registerCreateVariable(server)
	registerDeleteVariable(server)
	registerCreateContainer(server)
	registerDeleteContainer(server)
	registerCreateWorkspace(server)

	// Version operations
	registerCreateVersion(server)
	registerPublishVersion(server)

	// Templates (help LLMs with correct parameter formats)
	registerGetTagTemplates(server)
	registerGetTriggerTemplates(server)

	// Resources (URI-based read access)
	RegisterResources(server)

	// Prompts (template workflows)
	RegisterPrompts(server)
}

// getClient creates a GTM client from the request context with auto-refreshing tokens.
// In stdio mode, it uses the pre-configured token source.
// In HTTP mode, it uses the OAuth token from the request context.
// If restricted account/container IDs are configured, the client will enforce those restrictions.
func getClient(ctx context.Context) (*Client, error) {
	// Check if we're in stdio mode with a pre-configured token source
	if stdioTokenSource != nil {
		return NewClientWithRestriction(ctx, stdioTokenSource, restrictedAccountID, restrictedContainerIDs)
	}

	// HTTP mode: get token from request context
	tokenInfo := auth.GetTokenInfo(ctx)
	if tokenInfo == nil || tokenInfo.GoogleToken == nil {
		return nil, fmt.Errorf("not authenticated - please authenticate with Google first")
	}

	store := auth.GetTokenStore(ctx)
	google := auth.GetGoogleProvider(ctx)

	// Create auto-refreshing token source
	var tokenSource = auth.NewAutoRefreshTokenSource(
		store,
		tokenInfo.AccessToken,
		google.Config(),
		tokenInfo.GoogleToken,
	)

	return NewClientWithRestriction(ctx, tokenSource, restrictedAccountID, restrictedContainerIDs)
}

// RegisterToolsStdio adds all GTM tools to the MCP server for stdio mode.
// Uses the same tools as HTTP mode since getClient() handles both modes.
func RegisterToolsStdio(server *mcp.Server) {
	// All the same tools work in stdio mode
	RegisterTools(server)
}
