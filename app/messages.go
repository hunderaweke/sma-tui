package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MessagesModel renders a simple list of messages with basic navigation.
type MessagesModel struct {
	width    int
	height   int
	messages []string
	cursor   int
	identity string
}

func NewMessagesModel(identity string) MessagesModel {
	return MessagesModel{
		identity: identity,
		messages: []string{
			"Welcome to your inbox.",
			"This is a placeholder message.",
			"Navigation: j/k or ↑/↓ to move, esc to return.",
		},
	}
}

func (m MessagesModel) Init() tea.Cmd {
	return nil
}

func (m MessagesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "b":
			return MenuModel{State: StateMenu, width: m.width, height: m.height, identity: m.identity}, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if len(m.messages) == 0 {
				break
			}
			if m.cursor < len(m.messages)-1 {
				m.cursor++
			}
		case "enter":
			// No-op for now; placeholder for future message detail view.
		}
	}

	return m, nil
}

func (m MessagesModel) View() string {

	header := msgTitleStyle.Render("Your Messages")
	identity := msgIdentityStyle.Render(fmt.Sprintf("Identity: %s", m.identity))

	var list string
	if len(m.messages) == 0 {
		list = msgPlaceholderStyle.Render("No messages yet. Send one to see it here.")
	} else {
		items := make([]string, len(m.messages))
		for i, message := range m.messages {
			row := msgItemStyle.Render(message)
			if i == m.cursor {
				row = msgSelectedStyle.Render(message)
			}
			items[i] = row
		}
		list = lipgloss.JoinVertical(lipgloss.Left, items...)
	}

	instructions := msgInstructionStyle.Render("j/k or ↑/↓ to move • enter (noop) • esc to menu • q to quit")

	content := lipgloss.JoinVertical(lipgloss.Left, header, identity, list, instructions)
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}

var (
	msgTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Bold(true).
			Padding(1, 2)
	msgIdentityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(0, 2)
	msgItemStyle = lipgloss.NewStyle().
			Padding(0, 2)
	msgSelectedStyle = msgItemStyle.
				Foreground(lipgloss.Color("#D1D6D6")).
				Background(lipgloss.Color("#323A38")).
				Bold(true)
	msgPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(1, 2)
	msgInstructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(1, 2).
				Italic(true)
)
