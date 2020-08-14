package utils

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
)

// Spacer is a simple util for adding space between terminal messages
func Spacer() {
	fmt.Println("")
}

// RenderMarkdown accepts a markdown formatted string and renders it in the terminal
func RenderMarkdown(body string) string {
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	out := markdown.Render(body, 80, 6)
	return string(fmt.Sprintf("\n%s", out))
}
