package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func showHistory() {
	history := loadHistory()

	if len(history) == 0 {
		printInfo("No history available yet.")
		return
	}

	fmt.Printf("\n  %s\n\n", styleTitle.Render("COMMAND HISTORY"))

	for i, h := range history {
		num := lipgloss.NewStyle().Foreground(lipgloss.Color("#5F5F5F")).Render(fmt.Sprintf("%d.", i+1))
		q := lipgloss.NewStyle().Bold(true).Render(h.Question)
		c := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Render(h.Command)

		fmt.Printf("  %s %s\n", num, q)
		fmt.Printf("     %s %s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#875FFF")).Render("→"), c)
	}
}

func printInfo(msg string) {
	fmt.Printf("\n  %s %s\n\n", 
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7FF")).Render("ℹ"),
		msg,
	)
}
