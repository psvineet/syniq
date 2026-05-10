package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")).
			MarginLeft(2)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#875FFF")).
			Padding(1, 2).
			MarginLeft(1).
			Width(80)

	styleInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Italic(true)
)

func printBox(text string) {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(76),
	)

	out, err := r.Render(text)
	if err != nil {
		fmt.Println(text)
		return
	}

	fmt.Println(styleBox.Render(strings.TrimSpace(out)))
}

func printError(msg string) {
	fmt.Printf("\n  %s %s\n\n", 
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("✘"),
		msg,
	)
}

func printSuccess(msg string) {
	fmt.Printf("\n  %s %s\n\n", 
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✔"),
		msg,
	)
}
