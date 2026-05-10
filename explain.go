package main

import (
	"fmt"
)

func explain(cmd string) {
	prompt := fmt.Sprintf(
		"Explain the following Linux command in simple terms, highlighting its key flags and purpose:\n\n```bash\n%s\n```",
		cmd,
	)

	fmt.Printf("\n  %s Thinking...", styleInfo.Render("●"))

	result, err := callModel(prompt)
	// Clear the "Thinking..." line
	fmt.Print("\r\033[K")
	if err != nil {
		printError(fmt.Sprintf("Failed to fetch explanation: %v", err))
		return
	}

	printBox(result)
}
