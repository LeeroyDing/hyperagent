package tui

import (
"strings"

"github.com/charmbracelet/bubbles/spinner"
"github.com/charmbracelet/bubbles/viewport"
"github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"
)

var (
titleStyle = lipgloss.NewStyle().
Bold(true).
Foreground(lipgloss.Color("#FAFAFA")).
Background(lipgloss.Color("#7D56F4")).
Padding(0, 1)

infoStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#AF87FF"))

thoughtStyle = lipgloss.NewStyle().
Italic(true).
Foreground(lipgloss.Color("#888888"))

toolStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#04B575")).
Bold(true)
)

type Model struct {
Title    string
Thoughts []string
Action   string
Output   string
Spinner  spinner.Model
Viewport viewport.Model
Ready    bool
}

func NewModel(title string) Model {
s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
return Model{
Title:   title,
Spinner: s,
}
}

func (m Model) Init() tea.Cmd {
return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
var cmd tea.Cmd
switch msg := msg.(type) {
case tea.KeyMsg:
if msg.String() == "q" || msg.String() == "ctrl+c" {
return m, tea.Quit
}
case spinner.TickMsg:
m.Spinner, cmd = m.Spinner.Update(msg)
return m, cmd
case tea.WindowSizeMsg:
if !m.Ready {
m.Viewport = viewport.New(msg.Width, msg.Height-10)
m.Ready = true
}
m.Viewport.Width = msg.Width
m.Viewport.Height = msg.Height - 10
}
return m, nil
}

func (m Model) View() string {
if !m.Ready {
return "\n  Initializing..."
}

sb := strings.Builder{}
sb.WriteString(titleStyle.Render(m.Title) + "\n\n")

if len(m.Thoughts) > 0 {
sb.WriteString(infoStyle.Render("💭 Thoughts:") + "\n")
for _, t := range m.Thoughts {
sb.WriteString(thoughtStyle.Render("  • "+t) + "\n")
}
sb.WriteString("\n")
}

if m.Action != "" {
sb.WriteString(toolStyle.Render("🛠️ Action: ") + m.Action + "\n\n")
}

sb.WriteString(m.Spinner.View() + " Processing...\n\n")
m.Viewport.SetContent(m.Output)
sb.WriteString(m.Viewport.View())

return sb.String()
}
