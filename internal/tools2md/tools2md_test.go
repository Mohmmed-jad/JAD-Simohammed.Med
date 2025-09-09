package tools2md_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/github/github-mcp-server/internal/tools2md"
	"github.com/github/github-mcp-server/pkg/github"
	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-cmp/cmp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// func TestSimple(t *testing.T) {
// 	server := github.NewServer(
// 		func(_ context.Context) (*gogithub.Client, error) {
// 			panic("not implemented")
// 		},
// 		"0.0.1",
// 		false,
// 		translations.NullTranslationHelper,
// 	)

// 	// TODO: handle pagination
// 	// as per https://github.com/mark3labs/mcp-go/blob/cc777fcbf3176d0e76634f58047707d1f666cae8/client/stdio.go#L464
// 	request := mcp.JSONRPCRequest{
// 		JSONRPC: mcp.JSONRPC_VERSION,
// 		ID:      "1",
// 		Request: mcp.Request{
// 			Method: "tools/list",
// 		},
// 		Params: mcp.PaginatedRequest{
// 			Request: mcp.Request{},
// 		},
// 	}

// 	raw, err := json.Marshal(request)
// 	if err != nil {
// 		t.Fatalf("failed to marshal request: %v", err)
// 	}

// 	message := server.HandleMessage(
// 		context.Background(),
// 		raw,
// 	)

// 	response, ok := message.(mcp.JSONRPCResponse)
// 	require.True(t, ok, "expected JSONRPCResponse, got %T", message)

// 	listToolsResult, ok := response.Result.(mcp.ListToolsResult)
// 	require.True(t, ok, "expected ListToolsResult, got %T", response.Result)

// 	listToolsResult.Tools[0]

// }

var update = flag.Bool("update", false, "update .golden files")
var diffMd = flag.Bool("diffmd", false, "on failure create a .golden.diff.md file for external comparison")

func TestNoToolsReturnsEmptyString(t *testing.T) {
	tools := github.Tools{}

	result := tools2md.Convert(tools)
	require.Empty(t, result, "expected empty string when converting no tools")
}

func TestOneCategoryOneTool(t *testing.T) {
	goldenFilePath := goldenFilePath(t.Name())

	tools := github.Tools{
		{
			Definition: mcp.Tool{
				Name:        "test_tool",
				Description: "A tool for testing",
			},
			Category: "Test Category",
		},
	}

	md := tools2md.Convert(tools)
	if *update {
		require.NoError(
			t,
			os.WriteFile(goldenFilePath, []byte(md), 0600),
			"failed to update golden file",
		)
	}

	golden, err := os.ReadFile(goldenFilePath)
	require.NoError(t, err, "failed to read golden file")

	if diff := cmp.Diff(string(golden), md); diff != "" {
		if *diffMd {
			diffFilePath := strings.ReplaceAll(goldenFilePath, ".golden.md", ".golden.diff.md")
			require.NoError(
				t,
				os.WriteFile(diffFilePath, []byte(md), 0600),
				"failed to update diff file",
			)
		}

		t.Errorf("golden file mismatch\n%s", diff)
	}
}

func TestMultipleCategoriesMultipleTools(t *testing.T) {
	goldenFilePath := goldenFilePath(t.Name())

	tools := github.Tools{
		{
			Definition: mcp.Tool{
				Name:        "test_tool_1",
				Description: "A tool for testing",
			},
			Category: "Test Category 1",
		},
		{
			Definition: mcp.Tool{
				Name:        "test_tool_2",
				Description: "Another tool for testing",
			},
			Category: "Test Category 2",
		},
		{
			Definition: mcp.Tool{
				Name:        "test_tool_3",
				Description: "Yet another tool for testing",
			},
			Category: "Test Category 1",
		},
	}

	md := tools2md.Convert(tools)
	if *update {
		require.NoError(
			t,
			os.WriteFile(goldenFilePath, []byte(md), 0600),
			"failed to update golden file",
		)
	}

	golden, err := os.ReadFile(goldenFilePath)
	require.NoError(t, err, "failed to read golden file")

	if diff := cmp.Diff(string(golden), md); diff != "" {
		if *diffMd {
			diffFilePath := strings.ReplaceAll(goldenFilePath, ".golden.md", ".golden.diff.md")
			require.NoError(
				t,
				os.WriteFile(diffFilePath, []byte(md), 0600),
				"failed to update diff file",
			)
		}

		t.Errorf("golden file mismatch\n%s", diff)
	}
}

func TestToolsWithProperties(t *testing.T) {
	goldenFilePath := goldenFilePath(t.Name())

	tools := github.Tools{
		{
			Definition: mcp.Tool{
				Name:        "test_tool",
				Description: "A tool for testing",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]any{
						"prop_1": map[string]any{
							"description": "A test property",
							"type":        "string",
						},
						"prop_2": map[string]any{
							"description": "Another test property",
							"type":        "number",
						},
					},
					Required: []string{"prop_1"},
				},
			},
			Category: "Test Category",
		},
	}

	md := tools2md.Convert(tools)
	if *update {
		require.NoError(
			t,
			os.WriteFile(goldenFilePath, []byte(md), 0600),
			"failed to update golden file",
		)
	}

	golden, err := os.ReadFile(goldenFilePath)
	require.NoError(t, err, "failed to read golden file")

	if diff := cmp.Diff(string(golden), md); diff != "" {
		if *diffMd {
			diffFilePath := strings.ReplaceAll(goldenFilePath, ".golden.md", ".golden.diff.md")
			require.NoError(
				t,
				os.WriteFile(diffFilePath, []byte(md), 0600),
				"failed to update diff file",
			)
		}

		t.Errorf("golden file mismatch\n%s", diff)
	}
}

func TestFullSchema(t *testing.T) {
	goldenFilePath := goldenFilePath(t.Name())

	tools := github.DefaultTools(translations.NullTranslationHelper)

	md := tools2md.Convert(tools)
	if *update {
		require.NoError(
			t,
			os.WriteFile(goldenFilePath, []byte(md), 0600),
			"failed to update golden file",
		)
	}

	golden, err := os.ReadFile(goldenFilePath)
	require.NoError(t, err, "failed to read golden file")

	if diff := cmp.Diff(string(golden), md); diff != "" {
		if *diffMd {
			diffFilePath := strings.ReplaceAll(goldenFilePath, ".golden.md", ".golden.diff.md")
			require.NoError(
				t,
				os.WriteFile(diffFilePath, []byte(md), 0600),
				"failed to update diff file",
			)
		}

		t.Errorf("golden file mismatch\n%s", diff)
	}
}

// In case we use subtests
func goldenFilePath(testName string) string {
	return filepath.Join(
		"testdata",
		strings.ReplaceAll(testName+".golden.md", "/", "_"),
	)
}
