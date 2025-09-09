package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	stdlog "log"

	iolog "github.com/github/github-mcp-server/pkg/log"
	"github.com/github/github-mcp-server/pkg/translations"
	gogithub "github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RunConfig struct {
	Stdin  io.Reader
	Stdout io.Writer

	Version string

	Token string

	Logger      *log.Logger
	LogCommands bool

	ReadOnly           bool
	ExportTranslations bool
	EnabledToolsets    []string
}

func RunStdioServer(cfg RunConfig) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create GH client
	ghClient := gogithub.NewClient(nil).WithAuthToken(cfg.Token)
	ghClient.UserAgent = fmt.Sprintf("github-mcp-server/%s", cfg.Version)

	host := viper.GetString("host")

	if host != "" {
		var err error
		ghClient, err = ghClient.WithEnterpriseURLs(host, host)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client with host: %w", err)
		}
	}

	t, dumpTranslations := translations.TranslationHelper()

	beforeInit := func(_ context.Context, _ any, message *mcp.InitializeRequest) {
		ghClient.UserAgent = fmt.Sprintf(
			"github-mcp-server/%s (%s/%s)",
			cfg.Version,
			message.Params.ClientInfo.Name,
			message.Params.ClientInfo.Version,
		)
	}

	getClient := func(_ context.Context) (*gogithub.Client, error) {
		return ghClient, nil // closing over client
	}

	hooks := &server.Hooks{
		OnBeforeInitialize: []server.OnBeforeInitializeFunc{beforeInit},
	}
	// Create server
	ghServer := NewServer(cfg.Version, server.WithHooks(hooks))

	enabled := cfg.EnabledToolsets
	// TODO: tear this out
	dynamic := viper.GetBool("dynamic_toolsets")
	if dynamic {
		// filter "all" from the enabled toolsets
		enabled = make([]string, 0, len(cfg.EnabledToolsets))
		for _, toolset := range cfg.EnabledToolsets {
			if toolset != "all" {
				enabled = append(enabled, toolset)
			}
		}
	}

	// Create default toolsets
	toolsets, err := InitToolsets(enabled, cfg.ReadOnly, getClient, t)
	if err != nil {
		cfg.Logger.Fatal("Failed to initialize toolsets:", err)
	}

	context := InitContextToolset(getClient, t)

	// Register resources with the server
	RegisterResources(ghServer, getClient, t)
	// Register the tools with the server
	toolsets.RegisterTools(ghServer)
	context.RegisterTools(ghServer)

	if dynamic {
		dynamic := InitDynamicToolset(ghServer, toolsets, t)
		dynamic.RegisterTools(ghServer)
	}

	stdioServer := server.NewStdioServer(ghServer)

	stdLogger := stdlog.New(cfg.Logger.Writer(), "stdioserver", 0)
	stdioServer.SetErrorLogger(stdLogger)

	if cfg.ExportTranslations {
		// Once server is initialized, all translations are loaded
		dumpTranslations()
	}

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := cfg.Stdin, cfg.Stdout
		if cfg.LogCommands {
			loggedIO := iolog.NewIOLogger(in, out, cfg.Logger)
			in, out = loggedIO, loggedIO
		}

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output github-mcp-server string
	_, _ = fmt.Fprintf(os.Stderr, "GitHub MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		cfg.Logger.Infof("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
