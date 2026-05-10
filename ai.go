package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type pollResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func callModel(query string) (string, error) {
	url := "https://text.pollinations.ai/openai"
	
	prompt := fmt.Sprintf(
		"You are Syniq, a high-performance Linux terminal expert. "+
		"Return the requested shell command or explanation in concise markdown. "+
		"Commands MUST be inside triple backtick code blocks with the 'bash' language tag (e.g. ```bash\nls\n```). "+
		"Task: %s",
		query,
	)

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"model": "openai",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var poll pollResponse
	if err := json.Unmarshal(body, &poll); err != nil {
		return "", fmt.Errorf("parsing error: %v", err)
	}

	if len(poll.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return poll.Choices[0].Message.Content, nil
}
