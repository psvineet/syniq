package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7561FF")).
			Background(lipgloss.Color("#2E2E2E")).
			Padding(0, 1).
			MarginLeft(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginLeft(2)

	userInputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D7FF")).
			Bold(true)

	aiOutputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00"))

	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
)

type model struct {
	viewport    viewport.Model
	textInput   textinput.Model
	history     []string
	err         error
	loading     bool
	lastQuery   string
	lastCmd     string
	confirmRun  bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask Syniq anything (e.g. how to list files by size?)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	vp := viewport.New(100, 20)
	welcome := "Welcome to Syniq Interactive Chat!\nType your request below."
	vp.SetContent(welcome)

	return model{
		textInput: ti,
		viewport:  vp,
		history:   []string{welcome},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type aiResponseMsg string
type errMsg error

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirmRun {
			switch msg.String() {
			case "y", "Y":
				m.confirmRun = false
				return m, m.executeCommand(m.lastCmd)
			default:
				m.confirmRun = false
				m.history = append(m.history, infoStyle.Render("\nRun cancelled."))
				m.updateViewport()
				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlY:
			if m.lastCmd != "" {
				clipboard.WriteAll(m.lastCmd)
				m.history = append(m.history, successStyle.Render("\nCopied to clipboard!"))
				m.updateViewport()
			}
			return m, nil
		case tea.KeyCtrlR:
			if m.lastCmd != "" {
				m.confirmRun = true
				return m, nil
			}
		case tea.KeyEnter:
			query := m.textInput.Value()
			if query == "" {
				return m, nil
			}
			m.lastQuery = query
			m.textInput.SetValue("")
			m.loading = true
			m.history = append(m.history, fmt.Sprintf("\n%s %s", userInputStyle.Render("❯"), query))
			m.updateViewport()
			return m, m.callAI(query)
		}

	case aiResponseMsg:
		m.loading = false
		
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width-2),
		)
		out, err := r.Render(string(msg))
		if err != nil {
			out = string(msg)
		}

		m.history = append(m.history, fmt.Sprintf("\n%s\n%s", aiOutputStyle.Render("Syniq:"), strings.TrimSpace(out)))
		m.lastCmd = extractCommand(string(msg))
		m.updateViewport()
		saveHistory(m.lastQuery, string(msg))
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		m.history = append(m.history, fmt.Sprintf("\n%s %v", lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("ERROR:"), msg))
		m.updateViewport()
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 12
		if m.viewport.Height < 1 {
			m.viewport.Height = 1
		}
		m.textInput.Width = msg.Width - 10
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *model) updateViewport() {
	if len(m.history) > 0 {
		m.viewport.SetContent(strings.Join(m.history, "\n"))
		m.viewport.GotoBottom()
	}
}

func (m model) View() string {
	loading := ""
	if m.loading {
		loading = lipgloss.NewStyle().Foreground(lipgloss.Color("#5F5F5F")).Render(" (Thinking...)")
	}

	help := "ctrl+c quit • ctrl+y copy • ctrl+r run • enter send"
	if m.confirmRun {
		help = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000")).Render(fmt.Sprintf("Run command: %s? (y/n)", m.lastCmd))
	}

	return fmt.Sprintf(
		"%s%s\n\n%s\n\n%s\n%s",
		titleStyle.Render("SYNIQ INTERACTIVE"),
		loading,
		m.viewport.View(),
		m.textInput.View(),
		infoStyle.Render(help),
	)
}

func (m model) callAI(query string) tea.Cmd {
	return func() tea.Msg {
		res, err := callModel(query)
		if err != nil {
			return errMsg(err)
		}
		return aiResponseMsg(res)
	}
}

func (m model) executeCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("bash", "-c", cmd)
		out, err := c.CombinedOutput()
		if err != nil {
			return errMsg(fmt.Errorf("run failed: %v\n%s", err, string(out)))
		}
		return aiResponseMsg(fmt.Sprintf("Output:\n```\n%s\n```", string(out)))
	}
}

func extractCommand(text string) string {
	re := regexp.MustCompile("(?s)```(?:bash|sh)?\n(.*?)\n```")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func startInteractive() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}
