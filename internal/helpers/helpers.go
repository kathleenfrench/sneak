package helpers

import (
	"fmt"
	"os"
	"os/exec"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
)

// Spacer is a simple util for adding space between terminal messages
func Spacer() {
	fmt.Println("")
}

// GetKeysFromMap is a helper for fetching a slice of strings from a mpa
func GetKeysFromMap(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// SudoPing becasuse the alpine linux docker imnage does not play nice with normal ping
func SudoPing(ip string) error {
	cmd := exec.Command("sudo", "ping", "-c", fmt.Sprintf("%d", 1), ip)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

// RenderMarkdown accepts a markdown formatted string and renders it in the terminal
func RenderMarkdown(body string) string {
	markdown.BlueBgItalic = color.New(color.FgBlue).SprintFunc()
	out := markdown.Render(body, 80, 6)
	return string(fmt.Sprintf("\n%s", out))
}
