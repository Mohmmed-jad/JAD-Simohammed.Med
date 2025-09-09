package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GetClientFn func(context.Context) (*github.Client, error)

type Access int

const (
	// Zero value is writer, that way if forgotten, it won't be included
	// in read only configuration.
	Write Access = iota
	ReadOnly
)

type Handler func(getClient GetClientFn) server.ToolHandlerFunc

type Tool struct {
	Definition mcp.Tool
	Handler    Handler
	Access     Access
	Category   Category
}

type Category string

const (
	// CategoryUsers is the category for user-related tools.
	CategoryUsers Category = "Users"
	// CategoryIssues is the category for issue-related tools.
	CategoryIssues Category = "Issues"
	// CategoryPullRequests is the category for pull request-related tools.
	CategoryPullRequests Category = "Pull Requests"
	// CategoryRepositories is the category for repository-related tools.
	CategoryRepositories Category = "Repositories"
	// CategorySearch is the category for search-related tools.
	CategorySearch Category = "Search"
	// CategoryCodeScanning is the category for code scanning-related tools.
	CategoryCodeScanning Category = "Code Scanning"
)

type Tools []Tool

func (t Tools) ReadOnly() []Tool {
	var readOnlyTools []Tool
	for _, tool := range t {
		if tool.Access == ReadOnly {
			readOnlyTools = append(readOnlyTools, tool)
		}
	}
	return readOnlyTools
}

func DefaultTools(t translations.TranslationHelperFunc) Tools {
	return []Tool{
		// Users
		GetMe(t),

		// Issues
		GetIssue(t),
		SearchIssues(t),
		ListIssues(t),
		GetIssueComments(t),
		CreateIssue(t),
		AddIssueComment(t),
		UpdateIssue(t),

		// Pull Requests
		GetPullRequest(t),
		ListPullRequests(t),
		GetPullRequestFiles(t),
		GetPullRequestStatus(t),
		GetPullRequestComments(t),
		GetPullRequestReviews(t),
		MergePullRequest(t),
		UpdatePullRequestBranch(t),
		CreatePullRequestReview(t),
		CreatePullRequest(t),
		UpdatePullRequest(t),

		// Repositories
		SearchRepositories(t),
		GetFileContents(t),
		GetCommit(t),
		ListCommits(t),
		CreateOrUpdateFile(t),
		CreateRepository(t),
		ForkRepository(t),
		CreateBranch(t),
		PushFiles(t),

		// Search
		SearchCode(t),
		SearchUsers(t),

		// Code Scanning
		GetCodeScanningAlert(t),
		ListCodeScanningAlerts(t),
	}
}

// NewServer creates a new GitHub MCP server with the specified GH client and logger.
func NewServer(getClient GetClientFn, version string, readOnly bool, t translations.TranslationHelperFunc, opts ...server.ServerOption) *server.MCPServer {
	// Add default options
	defaultOpts := []server.ServerOption{
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	}
	opts = append(defaultOpts, opts...)

	// Create a new MCP server
	s := server.NewMCPServer(
		"github-mcp-server",
		version,
		opts...,
	)

	// // Add GitHub Resources
	s.AddResourceTemplate(GetRepositoryResourceContent(getClient, t))
	s.AddResourceTemplate(GetRepositoryResourceBranchContent(getClient, t))
	s.AddResourceTemplate(GetRepositoryResourceCommitContent(getClient, t))
	s.AddResourceTemplate(GetRepositoryResourceTagContent(getClient, t))
	s.AddResourceTemplate(GetRepositoryResourcePrContent(getClient, t))

	// Add GitHub Tools
	tools := DefaultTools(t)
	if readOnly {
		tools = tools.ReadOnly()
	}

	for _, tool := range tools {
		s.AddTool(tool.Definition, tool.Handler(getClient))
	}

	return s
}

// GetMe creates a tool to get details of the authenticated user.
func GetMe(t translations.TranslationHelperFunc) Tool {
	return Tool{
		Definition: mcp.NewTool("get_me",
			mcp.WithDescription(t("TOOL_GET_ME_DESCRIPTION", "Get details of the authenticated GitHub user. Use this when a request include \"me\", \"my\"...")),
			mcp.WithString("reason",
				mcp.Description("Optional: reason the session was created"),
			),
		),
		Handler: func(getClient GetClientFn) server.ToolHandlerFunc {
			return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client, err := getClient(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to get GitHub client: %w", err)
				}
				user, resp, err := client.Users.Get(ctx, "")
				if err != nil {
					return nil, fmt.Errorf("failed to get user: %w", err)
				}
				defer func() { _ = resp.Body.Close() }()

				if resp.StatusCode != http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return nil, fmt.Errorf("failed to read response body: %w", err)
					}
					return mcp.NewToolResultError(fmt.Sprintf("failed to get user: %s", string(body))), nil
				}

				r, err := json.Marshal(user)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal user: %w", err)
				}

				return mcp.NewToolResultText(string(r)), nil
			}
		},
		Access:   ReadOnly,
		Category: CategoryUsers,
	}
}

// OptionalParamOK is a helper function that can be used to fetch a requested parameter from the request.
// It returns the value, a boolean indicating if the parameter was present, and an error if the type is wrong.
func OptionalParamOK[T any](r mcp.CallToolRequest, p string) (value T, ok bool, err error) {
	// Check if the parameter is present in the request
	val, exists := r.Params.Arguments[p]
	if !exists {
		// Not present, return zero value, false, no error
		return
	}

	// Check if the parameter is of the expected type
	value, ok = val.(T)
	if !ok {
		// Present but wrong type
		err = fmt.Errorf("parameter %s is not of type %T, is %T", p, value, val)
		ok = true // Set ok to true because the parameter *was* present, even if wrong type
		return
	}

	// Present and correct type
	ok = true
	return
}

// isAcceptedError checks if the error is an accepted error.
func isAcceptedError(err error) bool {
	var acceptedError *github.AcceptedError
	return errors.As(err, &acceptedError)
}

// requiredParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request.
// 2. Checks if the parameter is of the expected type.
// 3. Checks if the parameter is not empty, i.e: non-zero value
func requiredParam[T comparable](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", p, zero)
	}

	if r.Params.Arguments[p].(T) == zero {
		return zero, fmt.Errorf("missing required parameter: %s", p)

	}

	return r.Params.Arguments[p].(T), nil
}

// RequiredInt is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request.
// 2. Checks if the parameter is of the expected type.
// 3. Checks if the parameter is not empty, i.e: non-zero value
func RequiredInt(r mcp.CallToolRequest, p string) (int, error) {
	v, err := requiredParam[float64](r, p)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// OptionalParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, it checks if the parameter is of the expected type and returns it
func OptionalParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, nil
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T, is %T", p, zero, r.Params.Arguments[p])
	}

	return r.Params.Arguments[p].(T), nil
}

// OptionalIntParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, it checks if the parameter is of the expected type and returns it
func OptionalIntParam(r mcp.CallToolRequest, p string) (int, error) {
	v, err := OptionalParam[float64](r, p)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// OptionalIntParamWithDefault is a helper function that can be used to fetch a requested parameter from the request
// similar to optionalIntParam, but it also takes a default value.
func OptionalIntParamWithDefault(r mcp.CallToolRequest, p string, d int) (int, error) {
	v, err := OptionalIntParam(r, p)
	if err != nil {
		return 0, err
	}
	if v == 0 {
		return d, nil
	}
	return v, nil
}

// OptionalStringArrayParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, iterates the elements and checks each is a string
func OptionalStringArrayParam(r mcp.CallToolRequest, p string) ([]string, error) {
	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return []string{}, nil
	}

	switch v := r.Params.Arguments[p].(type) {
	case nil:
		return []string{}, nil
	case []string:
		return v, nil
	case []any:
		strSlice := make([]string, len(v))
		for i, v := range v {
			s, ok := v.(string)
			if !ok {
				return []string{}, fmt.Errorf("parameter %s is not of type string, is %T", p, v)
			}
			strSlice[i] = s
		}
		return strSlice, nil
	default:
		return []string{}, fmt.Errorf("parameter %s could not be coerced to []string, is %T", p, r.Params.Arguments[p])
	}
}

// WithPagination returns a ToolOption that adds "page" and "perPage" parameters to the tool.
// The "page" parameter is optional, min 1. The "perPage" parameter is optional, min 1, max 100.
func WithPagination() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (min 1)"),
			mcp.Min(1),
		)(tool)

		mcp.WithNumber("perPage",
			mcp.Description("Results per page for pagination (min 1, max 100)"),
			mcp.Min(1),
			mcp.Max(100),
		)(tool)
	}
}

type PaginationParams struct {
	page    int
	perPage int
}

// OptionalPaginationParams returns the "page" and "perPage" parameters from the request,
// or their default values if not present, "page" default is 1, "perPage" default is 30.
// In future, we may want to make the default values configurable, or even have this
// function returned from `withPagination`, where the defaults are provided alongside
// the min/max values.
func OptionalPaginationParams(r mcp.CallToolRequest) (PaginationParams, error) {
	page, err := OptionalIntParamWithDefault(r, "page", 1)
	if err != nil {
		return PaginationParams{}, err
	}
	perPage, err := OptionalIntParamWithDefault(r, "perPage", 30)
	if err != nil {
		return PaginationParams{}, err
	}
	return PaginationParams{
		page:    page,
		perPage: perPage,
	}, nil
}
