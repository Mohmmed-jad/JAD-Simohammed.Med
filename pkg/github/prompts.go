package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewMePrompt(client *github.Client) (mcp.Prompt, server.PromptHandlerFunc) {
	prompt := mcp.NewPrompt("github_me", mcp.WithPromptDescription("GitHub Prompt"))
	prompt.Arguments = []mcp.PromptArgument{}

	return prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		user, resp, err := client.Users.Get(ctx, "")
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		r, err := json.Marshal(user)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user: %w", err)
		}

		return &mcp.GetPromptResult{
			Description: "Your GitHub Identity",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: string(r),
					},
				},
			},
		}, nil

	}

}
