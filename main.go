package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func printUsage() {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D7FF")).Render("SYNIQ")
	fmt.Printf("\n  %s — The AI-Powered Terminal Companion (v2.0.0)\n\n", title)
	fmt.Println("  Usage:")
	fmt.Println("    syniq chat            Start interactive chat session (FREE)")
	fmt.Println("    syniq ask <query>     Ask a quick question")
	fmt.Println("    syniq explain <cmd>   Explain a specific command")
	fmt.Println("    syniq history         Show recent command history")
	fmt.Println("    syniq version         Show version information")
	fmt.Println("\n  Interactive Shortcuts:")
	fmt.Println("    Ctrl+Y                Copy suggested command to clipboard")
	fmt.Println("    Ctrl+R                Run suggested command (with confirmation)")
	fmt.Println("\n  Examples:")
	fmt.Println("    syniq ask find all large logs")
	fmt.Println("    syniq explain tar -xzvf")
	fmt.Println()
}

func ask(question string) {
	fmt.Printf("\n  %s Thinking...", styleInfo.Render("●"))
	
	result, err := callModel(question)
	// Clear the "Thinking..." line
	fmt.Print("\r\033[K")

	if err != nil {
		fmt.Println("⚠️ Online lookup failed, checking local cache...")

		if cached, ok := findCachedAnswer(question); ok {
			fmt.Println("[CACHED RESULT — fuzzy match]")
			printBox(cached)
			return
		}

		printError(fmt.Sprintf("No cached result available: %v", err))
		return
	}

	status, _ := checkSafety(result)

	if status == "BLOCK" {
		printError("BLOCKED: This command is extremely dangerous.")
		return
	}

	if status == "WARN" {
		fmt.Printf("  %s %s\n", 
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAF00")).Render("⚠"),
			"WARNING: This command may be risky.",
		)
		fmt.Print("  Do you want to see it anyway? (yes/no): ")

		var answer string
		fmt.Scanln(&answer)

		if answer != "yes" {
			fmt.Println("  Cancelled.")
			return
		}
	}

	printBox(result)
	saveHistory(question, result)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "chat":
		startInteractive()
	case "ask":
		if len(os.Args) < 3 {
			printError("Usage: syniq ask <your question>")
			return
		}
		ask(strings.Join(os.Args[2:], " "))
	case "explain":
		if len(os.Args) < 3 {
			printError("Usage: syniq explain <command>")
			return
		}
		explain(strings.Join(os.Args[2:], " "))
	case "history":
		showHistory()
	case "version", "-v", "--version":
		fmt.Println("syniq v2.0.0")
	case "help":
		printUsage()
	default:
		printError(fmt.Sprintf("Unknown command: %s", command))
		printUsage()
	}
}
