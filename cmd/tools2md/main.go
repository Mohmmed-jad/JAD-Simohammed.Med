package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/github/github-mcp-server/internal/tools2md"
	"github.com/github/github-mcp-server/pkg/github"
	"github.com/github/github-mcp-server/pkg/translations"
)

var filepath = flag.String("filepath", "schema.md", "filepath to schema file")

func main() {
	tools := github.DefaultTools(translations.NullTranslationHelper)
	md := tools2md.Convert(tools)
	if *filepath == "" {
		panic("filepath cannot be empty")
	}

	err := os.WriteFile(*filepath, []byte(md), 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Schema file generated successfully at", *filepath)
}
