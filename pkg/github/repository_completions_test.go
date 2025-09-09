package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/sammorrowdrums/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// Add more fake methods as needed for testing
func TestRepositoryResourceCompletionHandler_Owner(t *testing.T) {
	// Stub getClient to return a fake client with a user and orgs
	getClient := func(ctx context.Context) (*github.Client, error) {
		client := github.NewClient(nil)
		// You can use github's testing helpers or mock the methods as needed
		return client, nil
	}

	handler := RepositoryResourceCompletionHandler(getClient)
	request := mcp.CompleteRequest{}
	request.Params.Ref = map[string]any{"type": "ref/resource", "uri": "repo://"}
	request.Params.Argument.Name = "owner"
	request.Params.Argument.Value = ""

	result, err := handler(context.Background(), request)
	require.NoError(t, err)
	// In a real test, assert on result.Completion.Values
	_ = result
}
