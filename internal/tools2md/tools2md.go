package tools2md

import (
	"fmt"
	"slices"
	"strings"

	"github.com/github/github-mcp-server/pkg/github"
)

func byCategoryPriority(a, b categorisedTools) int {
	var priorityMap = map[github.Category]int{
		github.CategoryUsers:        0,
		github.CategoryIssues:       1,
		github.CategoryPullRequests: 2,
		github.CategoryRepositories: 3,
		github.CategorySearch:       4,
		github.CategoryCodeScanning: 5,
	}

	pa, oka := priorityMap[a.category]
	pb, okb := priorityMap[b.category]

	// if both categories are not in the map, sort by our priorities
	if oka && okb {
		return pa - pb
	}

	// If either one was in the map then priortise that one.
	if oka {
		return -1
	}

	if okb {
		return 1
	}

	// if neither were in the map, sort alphabetically, which helps with test ordering.
	return strings.Compare(string(a.category), string(b.category))
}

type categorisedToolMap map[github.Category][]github.Tool

func (m categorisedToolMap) add(tool github.Tool) {
	m[tool.Category] = append(m[tool.Category], tool)
}

type categorisedTools struct {
	category github.Category
	tools    []github.Tool
}

type sortedCategorisedTools []categorisedTools

func (m categorisedToolMap) sorted() sortedCategorisedTools {
	var out sortedCategorisedTools
	for category, tools := range m {
		out = append(out, categorisedTools{
			category: category,
			tools:    tools,
		})
	}

	slices.SortStableFunc(out, byCategoryPriority)
	return out
}

// Replace a lot of this with a Go template.
// Much TDD.
func Convert(tools github.Tools) string {
	if len(tools) == 0 {
		return ""
	}

	toolMap := categorisedToolMap{}
	for _, tool := range tools {
		toolMap.add(tool)
	}

	var md markdownBuilder

	md.h2("Tools")
	md.newline()

	sortedToolMap := toolMap.sorted()
	for i, categorisedTools := range sortedToolMap {
		md.h3(string(categorisedTools.category))
		md.newline()

		for j, tool := range categorisedTools.tools {
			md.textf("- %s - %s", bold(tool.Definition.Name), tool.Definition.Description)
			md.newline()

			if len(tool.Definition.InputSchema.Properties) == 0 {
				md.text(" - No parameters required")
				md.newline()
			} else {
				// order the properties alphabetically to maintain a consistent order
				// maybe in future do some kind of grouping like pagination together.
				var propNames []string
				for propName := range tool.Definition.InputSchema.Properties {
					propNames = append(propNames, propName)
				}
				slices.Sort(propNames)

				for _, propName := range propNames {
					prop := tool.Definition.InputSchema.Properties[propName]
					propSchema := prop.(map[string]any)
					required := func() string {
						if slices.Contains(tool.Definition.InputSchema.Required, propName) {
							return "required"
						}
						return "optional"
					}()

					md.textf(" - %s: %s (%s, %s)", code(propName), propSchema["description"], propSchema["type"], required)
					md.newline()
				}
			}

			// if not the last tool in the category, add a newline
			if j < len(categorisedTools.tools)-1 {
				md.newline()
			}
		}

		// if not the last category, add a newline
		if i < len(sortedToolMap)-1 {
			md.newline()
		}
	}

	return md.String()
}

type markdownBuilder struct {
	content strings.Builder
}

func (b *markdownBuilder) h2(text string) {
	b.content.WriteString(fmt.Sprintf("## %s\n", text))
}

func (b *markdownBuilder) h3(text string) {
	b.content.WriteString(fmt.Sprintf("### %s\n", text))
}

func (b *markdownBuilder) bold(text string) {
	b.content.WriteString(fmt.Sprintf("**%s**", text))
}

func (b *markdownBuilder) newline() {
	b.content.WriteString("\n")
}

func (b *markdownBuilder) text(text string) {
	b.content.WriteString(text)
}

func (b *markdownBuilder) textf(format string, args ...any) {
	b.content.WriteString(fmt.Sprintf(format, args...))
}

func (b *markdownBuilder) String() string {
	return b.content.String()
}

func bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func code(text string) string {
	return fmt.Sprintf("`%s`", text)
}
