package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"gtm-mcp-server/auth"
	"gtm-mcp-server/config"
	"gtm-mcp-server/gtm"
	"gtm-mcp-server/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "gtm-mcp-server"
	serverVersion = "0.1.0"
)

func main() {
	// Set up structured logging to stderr (stdout is reserved for MCP in stdio mode)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Adjust log level
	if cfg.LogLevel == "debug" {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		slog.SetDefault(logger)
	}

	// Configure GTM account restriction if set
	if cfg.GTMAccountID != "" {
		gtm.SetRestrictedAccountID(cfg.GTMAccountID)
		logger.Info("GTM account restriction enabled", "account_id", cfg.GTMAccountID)
	}

	// Configure GTM container restrictions if set
	if len(cfg.GTMContainerIDs) > 0 {
		gtm.SetRestrictedContainerIDs(cfg.GTMContainerIDs)
		logger.Info("GTM container restriction enabled", "container_ids", cfg.GTMContainerIDs)
	}

	// Run in appropriate mode
	switch cfg.Transport {
	case "stdio":
		runStdioServer(cfg, logger)
	default:
		runHTTPServer(cfg, logger)
	}
}

// runStdioServer runs the MCP server using stdio transport.
// This mode requires GOOGLE_TOKEN_FILE to be set with a valid OAuth token.
func runStdioServer(cfg *config.Config, logger *slog.Logger) {
	logger.Info("starting GTM MCP server in stdio mode")

	// Try to load token if configured
	tokenLoaded := false
	if cfg.GoogleTokenFile != "" {
		tokenSource, err := auth.LoadTokenFromFile(cfg.GoogleTokenFile, cfg.GoogleClientID, cfg.GoogleClientSecret)
		if err != nil {
			logger.Warn("token file not found - will start auth flow",
				"error", err, "path", cfg.GoogleTokenFile)
		} else {
			gtm.SetStdioTokenSource(tokenSource)
			logger.Info("loaded Google token", "path", cfg.GoogleTokenFile)
			tokenLoaded = true
		}
	}

	// If no token, start background auth server and open browser
	if !tokenLoaded && cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		go startAuthServerAndOpenBrowser(cfg, logger)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Register tools (stdio mode uses simplified tool registration)
	registerToolsStdio(mcpServer)

	logger.Info("stdio server ready, waiting for input")

	// Run the server on stdio transport
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Error("stdio server error", "error", err)
		os.Exit(1)
	}
}

// runHTTPServer runs the MCP server using HTTP transport with OAuth.
func runHTTPServer(cfg *config.Config, logger *slog.Logger) {
	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Add logging middleware
	mcpServer.AddReceivingMiddleware(middleware.NewLoggingMiddleware(logger))

	// Register tools
	registerTools(mcpServer)

	// Create HTTP handler for MCP
	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	// Set up HTTP routes
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": serverName,
			"version": serverVersion,
		})
	})

	// OAuth metadata endpoints (always served, no auth required)
	mux.HandleFunc("GET /.well-known/oauth-protected-resource",
		auth.ProtectedResourceMetadataHandler(cfg.BaseURL, cfg.BaseURL))
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", auth.MetadataHandler(cfg.BaseURL))

	// Check if OAuth is configured
	var authServer *auth.Server
	var tokenStore auth.TokenStore
	oauthConfigured := cfg.ValidateAuth() == nil

	if oauthConfigured {
		// Set up OAuth
		tokenStore = auth.NewMemoryTokenStore()
		googleProvider := auth.NewGoogleProvider(
			cfg.GoogleClientID,
			cfg.GoogleClientSecret,
			cfg.BaseURL+"/oauth/callback",
		)
		authServer = auth.NewServer(cfg.BaseURL, googleProvider, tokenStore, logger)

		// Configure token saving if path is set
		if cfg.SaveTokenPath != "" {
			authServer.SetSaveTokenPath(cfg.SaveTokenPath)
			logger.Info("token saving enabled", "path", cfg.SaveTokenPath)
		}

		// OAuth endpoints (no auth required)
		mux.HandleFunc("GET /authorize", authServer.AuthorizeHandler)
		mux.HandleFunc("GET /oauth/callback", authServer.CallbackHandler)
		mux.HandleFunc("POST /token", authServer.TokenHandler)
		mux.HandleFunc("POST /register", authServer.RegistrationHandler)

		// MCP endpoint with REQUIRED auth middleware
		authMiddleware := auth.Middleware(tokenStore, googleProvider, logger, cfg.BaseURL)
		mux.Handle("/", authMiddleware(mcpHandler))

		logger.Info("OAuth configured",
			"authorize_endpoint", cfg.BaseURL+"/authorize",
			"token_endpoint", cfg.BaseURL+"/token",
			"callback_endpoint", cfg.BaseURL+"/oauth/callback",
			"register_endpoint", cfg.BaseURL+"/register",
		)
	} else {
		logger.Warn("OAuth not configured, running without authentication", "error", cfg.ValidateAuth())

		oauthNotConfiguredHandler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error":             "server_error",
				"error_description": "OAuth is not configured. Set GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, and JWT_SECRET.",
			})
		}
		mux.HandleFunc("GET /authorize", oauthNotConfiguredHandler)
		mux.HandleFunc("GET /oauth/callback", oauthNotConfiguredHandler)
		mux.HandleFunc("POST /token", oauthNotConfiguredHandler)
		mux.HandleFunc("POST /register", oauthNotConfiguredHandler)

		mux.Handle("/", mcpHandler)
	}

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // Disabled for SSE streams
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		logger.Info("starting GTM MCP server", "port", cfg.Port, "base_url", cfg.BaseURL)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	logger.Info("server stopped")
}

// registerTools adds MCP tools to the server (HTTP mode with OAuth context).
func registerTools(server *mcp.Server) {
	registerUtilityTools(server)
	gtm.RegisterTools(server)
}

// registerToolsStdio adds MCP tools to the server (stdio mode with file-based token).
func registerToolsStdio(server *mcp.Server) {
	registerUtilityToolsStdio(server)
	gtm.RegisterToolsStdio(server)
}

// registerUtilityTools adds ping and auth_status tools (HTTP mode).
func registerUtilityTools(server *mcp.Server) {
	type PingInput struct {
		Message string `json:"message,omitempty" jsonschema:"Optional message to echo back"`
	}
	type PingOutput struct {
		Reply     string `json:"reply"`
		Timestamp string `json:"timestamp"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "Test connectivity to the GTM MCP server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PingInput) (*mcp.CallToolResult, PingOutput, error) {
		reply := "pong"
		if input.Message != "" {
			reply = fmt.Sprintf("pong: %s", input.Message)
		}
		return nil, PingOutput{Reply: reply, Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
	})

	type AuthStatusInput struct{}
	type AuthStatusOutput struct {
		Authenticated          bool     `json:"authenticated"`
		Message                string   `json:"message"`
		RestrictedAccountID    string   `json:"restrictedAccountId,omitempty"`
		RestrictedContainerIDs []string `json:"restrictedContainerIds,omitempty"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "auth_status",
		Description: "Check authentication status with Google Tag Manager",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AuthStatusInput) (*mcp.CallToolResult, AuthStatusOutput, error) {
		tokenInfo := auth.GetTokenInfo(ctx)
		output := AuthStatusOutput{
			Authenticated:          tokenInfo != nil,
			RestrictedAccountID:    gtm.GetRestrictedAccountID(),
			RestrictedContainerIDs: gtm.GetRestrictedContainerIDs(),
		}
		if tokenInfo != nil {
			var restrictions []string
			if output.RestrictedAccountID != "" {
				restrictions = append(restrictions, fmt.Sprintf("account %s", output.RestrictedAccountID))
			}
			if len(output.RestrictedContainerIDs) > 0 {
				restrictions = append(restrictions, fmt.Sprintf("containers %v", output.RestrictedContainerIDs))
			}
			if len(restrictions) > 0 {
				output.Message = fmt.Sprintf("Authenticated. Access restricted to %s", strings.Join(restrictions, " and "))
			} else {
				output.Message = "Authenticated. Access to all accounts and containers you have permission for"
			}
		} else {
			output.Message = "Not authenticated. GTM tools will require authentication."
		}
		return nil, output, nil
	})
}

// registerUtilityToolsStdio adds ping and auth_status tools (stdio mode).
func registerUtilityToolsStdio(server *mcp.Server) {
	type PingInput struct {
		Message string `json:"message,omitempty" jsonschema:"Optional message to echo back"`
	}
	type PingOutput struct {
		Reply     string `json:"reply"`
		Timestamp string `json:"timestamp"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Description: "Test connectivity to the GTM MCP server",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input PingInput) (*mcp.CallToolResult, PingOutput, error) {
		reply := "pong"
		if input.Message != "" {
			reply = fmt.Sprintf("pong: %s", input.Message)
		}
		return nil, PingOutput{Reply: reply, Timestamp: time.Now().UTC().Format(time.RFC3339)}, nil
	})

	type AuthStatusInput struct{}
	type AuthStatusOutput struct {
		Authenticated          bool     `json:"authenticated"`
		Message                string   `json:"message"`
		Mode                   string   `json:"mode"`
		RestrictedAccountID    string   `json:"restrictedAccountId,omitempty"`
		RestrictedContainerIDs []string `json:"restrictedContainerIds,omitempty"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "auth_status",
		Description: "Check authentication status with Google Tag Manager",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AuthStatusInput) (*mcp.CallToolResult, AuthStatusOutput, error) {
		output := AuthStatusOutput{
			Authenticated:          gtm.HasStdioTokenSource(),
			Mode:                   "stdio",
			RestrictedAccountID:    gtm.GetRestrictedAccountID(),
			RestrictedContainerIDs: gtm.GetRestrictedContainerIDs(),
		}
		var restrictions []string
		if output.RestrictedAccountID != "" {
			restrictions = append(restrictions, fmt.Sprintf("account %s", output.RestrictedAccountID))
		}
		if len(output.RestrictedContainerIDs) > 0 {
			restrictions = append(restrictions, fmt.Sprintf("containers %v", output.RestrictedContainerIDs))
		}
		if output.Authenticated {
			if len(restrictions) > 0 {
				output.Message = fmt.Sprintf("Authenticated via token file. Access restricted to %s", strings.Join(restrictions, " and "))
			} else {
				output.Message = "Authenticated via token file. Access to all accounts and containers."
			}
		} else {
			output.Message = "No token loaded. Opening browser for authentication..."
		}
		return nil, output, nil
	})
}

// startAuthServerAndOpenBrowser starts a temporary HTTP server for OAuth and opens the browser.
// Once authentication completes, it loads the token and shuts down the HTTP server.
func startAuthServerAndOpenBrowser(cfg *config.Config, logger *slog.Logger) {
	// Use a different port for the auth server to avoid conflicts
	authPort := 18082
	baseURL := fmt.Sprintf("http://localhost:%d", authPort)

	// Set up OAuth
	tokenStore := auth.NewMemoryTokenStore()
	googleProvider := auth.NewGoogleProvider(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		baseURL+"/oauth/callback",
	)
	authServer := auth.NewServer(baseURL, googleProvider, tokenStore, logger)

	// Configure token saving
	if cfg.GoogleTokenFile != "" {
		authServer.SetSaveTokenPath(cfg.GoogleTokenFile)
	}

	// Set up simple login flow for stdio mode (direct browser auth)
	mux := http.NewServeMux()

	// Simple login - redirects directly to Google
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		// Generate state for CSRF protection
		state, _ := auth.GenerateToken(16)

		// Store state (simple in-memory for this flow)
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			MaxAge:   300,
			HttpOnly: true,
		})

		// Redirect to Google
		authURL := googleProvider.AuthCodeURL(state)
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	// Callback from Google
	mux.HandleFunc("GET /oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		// Verify state
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Check for errors
		if errCode := r.URL.Query().Get("error"); errCode != "" {
			http.Error(w, "OAuth error: "+errCode, http.StatusBadRequest)
			return
		}

		// Exchange code for token
		code := r.URL.Query().Get("code")
		token, err := googleProvider.Exchange(r.Context(), code)
		if err != nil {
			logger.Error("failed to exchange code", "error", err)
			http.Error(w, "Failed to exchange code: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Save token to file
		if cfg.GoogleTokenFile != "" {
			if err := auth.SaveTokenToFile(cfg.GoogleTokenFile, token); err != nil {
				logger.Error("failed to save token", "error", err)
				http.Error(w, "Failed to save token: "+err.Error(), http.StatusInternalServerError)
				return
			}
			logger.Info("saved token to file", "path", cfg.GoogleTokenFile)

			// Load it into gtm package
			if tokenSource, err := auth.LoadTokenFromFile(cfg.GoogleTokenFile, cfg.GoogleClientID, cfg.GoogleClientSecret); err == nil {
				gtm.SetStdioTokenSource(tokenSource)
				logger.Info("authentication successful - token loaded")
			}
		}

		// Show success page
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>GTM MCP - Authenticated</title></head>
<body style="font-family: system-ui; padding: 40px; text-align: center;">
<h1>✓ Authentication Successful!</h1>
<p>You can close this window and return to Claude Code.</p>
<p>GTM tools are now available.</p>
</body>
</html>`))
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", authPort),
		Handler: mux,
	}

	// Start the server
	go func() {
		logger.Info("starting auth server", "port", authPort, "url", baseURL+"/authorize")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("auth server error", "error", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(200 * time.Millisecond)

	// Open browser
	authURL := baseURL + "/login"
	logger.Info("opening browser for authentication", "url", authURL)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", authURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", authURL)
	default: // linux, etc.
		cmd = exec.Command("xdg-open", authURL)
	}

	if err := cmd.Start(); err != nil {
		logger.Warn("failed to open browser automatically", "error", err, "url", authURL)
		logger.Info("please open this URL manually to authenticate", "url", authURL)
	}
}
