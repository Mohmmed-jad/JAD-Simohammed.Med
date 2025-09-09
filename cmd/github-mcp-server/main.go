package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/github/github-mcp-server/pkg/github"
	iolog "github.com/github/github-mcp-server/pkg/log"
	"github.com/github/github-mcp-server/pkg/translations"
	gogithub "github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "GitHub MCP Server",
		Long:  `A GitHub MCP server that handles various tools and resources.`,
	}

	stdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Start stdio server",
		Long:  `Start a server that communicates via standard input/output streams using JSON-RPC messages.`,
		Run: func(cmd *cobra.Command, args []string) {
			logFile := viper.GetString("log-file")
			readOnly := viper.GetBool("read-only")
			exportTranslations := viper.GetBool("export-translations")
			logger, err := initLogger(logFile)
			if err != nil {
				stdlog.Fatal("Failed to initialize logger:", err)
			}
			logCommands := viper.GetBool("enable-command-logging")
			if err := runStdioServer(readOnly, logger, logCommands, exportTranslations); err != nil {
				stdlog.Fatal("failed to run stdio server:", err)
			}
		},
	}

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Start HTTP server",
		Long:  `Start a server that communicates via HTTP using Server-Sent Events (SSE).`,
		Run: func(cmd *cobra.Command, args []string) {
			logFile := viper.GetString("log-file")
			readOnly := viper.GetBool("read-only")
			exportTranslations := viper.GetBool("export-translations")
			port := viper.GetString("port")
			logger, err := initLogger(logFile)
			if err != nil {
				stdlog.Fatal("Failed to initialize logger:", err)
			}
			logCommands := viper.GetBool("enable-command-logging")
			if err := runHTTPServer(readOnly, logger, logCommands, exportTranslations, port); err != nil {
				stdlog.Fatal("failed to run http server:", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Add global flags that will be shared by all commands
	rootCmd.PersistentFlags().Bool("read-only", false, "Restrict the server to read-only operations")
	rootCmd.PersistentFlags().String("log-file", "", "Path to log file")
	rootCmd.PersistentFlags().Bool("enable-command-logging", false, "When enabled, the server will log all command requests and responses to the log file")
	rootCmd.PersistentFlags().Bool("export-translations", false, "Save translations to a JSON file")
	rootCmd.PersistentFlags().String("gh-host", "", "Specify the GitHub hostname (for GitHub Enterprise etc.)")

	// Add HTTP specific flags
	httpCmd.Flags().String("port", "8080", "Port for the HTTP server")

	// Bind flags to viper
	viper.BindPFlag("read-only", rootCmd.PersistentFlags().Lookup("read-only"))
	viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("enable-command-logging", rootCmd.PersistentFlags().Lookup("enable-command-logging"))
	viper.BindPFlag("export-translations", rootCmd.PersistentFlags().Lookup("export-translations"))
	viper.BindPFlag("gh-host", rootCmd.PersistentFlags().Lookup("gh-host"))
	viper.BindPFlag("port", httpCmd.Flags().Lookup("port"))

	// Add subcommands
	rootCmd.AddCommand(stdioCmd)
	rootCmd.AddCommand(httpCmd)
}

func initConfig() {
	// Initialize Viper configuration
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
}

func initLogger(outPath string) (*log.Logger, error) {
	if outPath == "" {
		return log.New(), nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetOutput(file)

	return logger, nil
}

func runStdioServer(readOnly bool, logger *log.Logger, logCommands bool, exportTranslations bool) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create GH client
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		logger.Fatal("GITHUB_PERSONAL_ACCESS_TOKEN not set")
	}
	ghClient := gogithub.NewClient(nil).WithAuthToken(token)

	// Check GH_HOST env var first, then fall back to viper config
	host := os.Getenv("GH_HOST")
	if host == "" {
		host = viper.GetString("gh-host")
	}

	if host != "" {
		var err error
		ghClient, err = ghClient.WithEnterpriseURLs(host, host)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client with host: %w", err)
		}
	}

	t, dumpTranslations := translations.TranslationHelper()

	// Create
	ghServer := github.NewServer(ghClient, readOnly, t)
	stdioServer := server.NewStdioServer(ghServer)

	stdLogger := stdlog.New(logger.Writer(), "stdioserver", 0)
	stdioServer.SetErrorLogger(stdLogger)

	if exportTranslations {
		// Once server is initialized, all translations are loaded
		dumpTranslations()
	}

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

		if logCommands {
			loggedIO := iolog.NewIOLogger(in, out, logger)
			in, out = loggedIO, loggedIO
		}

		errC <- stdioServer.Listen(ctx, in, out)
	}()

	// Output github-mcp-server string
	_, _ = fmt.Fprintf(os.Stderr, "GitHub MCP Server running on stdio\n")

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		logger.Infof("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}

func runHTTPServer(readOnly bool, logger *log.Logger, logCommands bool, exportTranslations bool, port string) error {
	// Create app context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create GH client
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		logger.Fatal("GITHUB_PERSONAL_ACCESS_TOKEN not set")
	}
	ghClient := gogithub.NewClient(nil).WithAuthToken(token)

	// Check GH_HOST env var first, then fall back to viper config
	host := os.Getenv("GH_HOST")
	if host == "" {
		host = viper.GetString("gh-host")
	}

	if host != "" {
		var err error
		ghClient, err = ghClient.WithEnterpriseURLs(host, host)
		if err != nil {
			return fmt.Errorf("failed to create GitHub client with host: %w", err)
		}
	}

	t, dumpTranslations := translations.TranslationHelper()

	// Create GitHub server
	ghServer := github.NewServer(ghClient, readOnly, t)

	if exportTranslations {
		// Once server is initialized, all translations are loaded
		dumpTranslations()
	}

	// Create SSE server
	sseServer := server.NewSSEServer(ghServer)

	// Start listening for messages
	errC := make(chan error, 1)
	go func() {
		// Configure and start HTTP server
		mux := http.NewServeMux()

		// Add SSE handler with logging middleware if enabled
		var handler http.Handler = sseServer
		if logCommands {
			handler = loggingMiddleware(handler, logger)
		}
		mux.Handle("/", handler)

		srv := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		// Graceful shutdown
		go func() {
			<-ctx.Done()
			if err := srv.Shutdown(context.Background()); err != nil {
				logger.Errorf("HTTP server shutdown error: %v", err)
			}
		}()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errC <- err
		}
	}()

	// Output github-mcp-server string
	_, _ = fmt.Fprintf(os.Stderr, "GitHub MCP Server running on http://localhost:%s\n", port)

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		logger.Infof("shutting down server...")
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}

// loggingMiddleware wraps an http.Handler and logs requests
func loggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Received request")

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
