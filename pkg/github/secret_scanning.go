package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getSecretScanningAlert(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_secret_scanning_alert",
			mcp.WithDescription(t("TOOL_GET_SECRET_SCANNING_ALERT_DESCRIPTION", "Get details of a specific secret scanning alert in a GitHub repository.")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The owner of the repository."),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository."),
			),
			mcp.WithNumber("alert_number",
				mcp.Required(),
				mcp.Description("The number of the alert."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := requiredParam[string](request, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			repo, err := requiredParam[string](request, "repo")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			alertNumber, err := requiredInt(request, "alert_number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			alert, resp, err := client.SecretScanning.GetAlert(ctx, owner, repo, int64(alertNumber))
			if err != nil {
				return nil, fmt.Errorf("failed to get alert: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get alert: %s", string(body))), nil
			}

			r, err := json.Marshal(alert)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal alert: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

func listSecretScanningAlerts(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_secret_scanning_alerts",
			mcp.WithDescription(t("TOOL_LIST_SECRET_SCANNING_ALERTS_DESCRIPTION", "List secret scanning alerts in a GitHub repository.")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The owner of the repository."),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository."),
			),
			mcp.WithString("secret_type",
				mcp.Description("A comma-separated list of secret types to return. All default secret patterns are returned. To return generic patterns, pass the token name(s) in the parameter."),
			),
			mcp.WithString("state",
				mcp.Description("State of the secret scanning alerts to list. Set to open or resolved to only list secret scanning alerts in a specific state."),
				mcp.DefaultString("open"),
			),
			mcp.WithString("resolution",
				mcp.Description("A comma-separated list of resolutions. Only secret scanning alerts with one of these resolutions are listed. Valid resolutions are false_positive, wont_fix, revoked, pattern_edited, pattern_deleted or used_in_tests."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := requiredParam[string](request, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			repo, err := requiredParam[string](request, "repo")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			secretType, err := optionalParam[string](request, "secret_type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			state, err := optionalParam[string](request, "state")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resolution, err := optionalParam[string](request, "resolution")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			alerts, resp, err := client.SecretScanning.ListAlertsForRepo(ctx, owner, repo, &github.SecretScanningAlertListOptions{SecretType: secretType, State: state, Resolution: resolution})
			if err != nil {
				return nil, fmt.Errorf("failed to list alerts: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list alerts: %s", string(body))), nil
			}

			r, err := json.Marshal(alerts)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal alerts: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
