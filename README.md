# SYNIQ 🚀
> **The AI-Powered Linux Terminal Companion — No API Key Required.**

Syniq is a high-performance, local-first CLI tool designed to bridge the gap between human intent and the Linux command line. It translates natural language into safe, accurate shell commands using a free, public AI model.

## 📺 Features at a Glance
- **Zero Configuration**: No API keys, tokens, or signups. It works instantly.
- **Interactive TUI**: A premium, full-screen dashboard for deep exploration.
- **Quote-less CLI**: Ask questions directly without wrapping them in quotes.
- **Pro Rendering**: Markdown support with full syntax highlighting (via Glamour).
- **Safety Engine**: Real-time analysis to block or warn against destructive commands.
- **One-Touch Actions**: 
    - `Ctrl+Y` to copy suggested commands.
    - `Ctrl+R` to execute commands directly (after confirmation).

---

## 🛠️ Installation
### Prerequisites
- **Go 1.22** or higher.

### Build from Source
```bash
git clone https://github.com/Vptsh/syniq.git
cd syniq
go mod tidy
go build -o syniq
sudo mv syniq /usr/local/bin/syniq
```

---

## 🚀 Usage

### 1. Interactive Mode (Recommended)
Start a full-screen chat session to iterate on complex tasks:
```bash
syniq chat
```

### 2. Quick Command Lookup
Get a one-off command for a specific task (quotes are optional):
```bash
syniq ask how to find all files larger than 100MB
```

### 3. Deep Explanation
Break down what a complex command or set of flags actually does:
```bash
syniq explain tar -xzvf archive.tar.gz
```

### 4. History
Review your previous queries and the AI's suggestions:
```bash
syniq history
```

---

## ⌨️ Interactive Shortcuts
While in **`syniq chat`** mode:
| Shortcut | Action |
| :--- | :--- |
| **Enter** | Send your query to the AI |
| **Ctrl + Y** | Copy the suggested command to clipboard |
| **Ctrl + R** | Execute the suggested command (requires `y` confirmation) |
| **Ctrl + C** | Exit the application |
| **Esc** | Exit the application |

---

## 🛡️ Safety & Privacy
Syniq is built with a **Safety First** philosophy:
1. **SAFE**: Standard non-destructive commands are displayed normally.
2. **RISKY**: Commands involving `rm`, `shutdown`, or `systemctl` trigger a warning.
3. **BLOCKED**: Commands that target the root filesystem (`rm -rf /`) or partition tables are hard-blocked.

**Privacy**: Syniq uses an anonymous public model endpoint. No personal data, telemetry, or system information is ever sent to a project backend.

---

## 📂 Project Architecture
- `ai.go`: Integration with public, unauthenticated AI providers.
- `tui.go`: The interactive `bubbletea` state machine and viewport.
- `main.go`: Command routing and the CLI "Thinking" engine.
- `safety.go`: Pattern-matching engine for command validation.
- `ui.go`: Lip Gloss styling and terminal aesthetics.

---

## 📜 License
MIT License. Created by **Vptsh** (psvineet@zohomail.in).
