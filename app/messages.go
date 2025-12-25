package app

import (
	"fmt"
	"time"

	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

// MessagesModel renders a simple list of messages with basic navigation.
type MessagesModel struct {
	width    int
	height   int
	messages []Message
	cursor   int
	identity string
	vp       viewport.Model
}

type Message struct {
	SentAt time.Time
	From   string
	Text   string
}

func NewMessagesModel(identity string) MessagesModel {
	return MessagesModel{
		identity: identity,
		messages: []Message{
			{SentAt: time.Now(), From: "somerand", Text: "Welcome to your inbox."},
			{SentAt: time.Now().Add(2 * time.Minute), From: "somerand2", Text: "This is a placeholder message."},
			{SentAt: time.Now().Add(10 * time.Minute), From: "somerand3", Text: "Navigation: j/k or ↑/↓ to move, esc to return."},
			{SentAt: time.Now().Add(20 * time.Minute), From: "somerand4", Text: "This is that is a phrase used for identification, where this points to something near and that points to something far, both referring to a single item, but it's also the name of a popular Canadian news satire show on CBC Radio. It highlights the basic grammar rule for demonstrative pronouns, contrasting proximity (this) with distance (that) in a humorous, news-style format, as seen in the show This Is That. "},
			{SentAt: time.Now().Add(120 * time.Minute), From: "somerand4", Text: "This is that is a phrase used for identification, where this points to something near and that points to something far, both referring to a single item, but it's also the name of a popular Canadian news satire show on CBC Radio. It highlights the basic grammar rule for demonstrative pronouns, contrasting proximity (this) with distance (that) in a humorous, news-style format, as seen in the show This Is That. "},
			{SentAt: time.Now().Add(120 * time.Minute), From: "somerand4", Text: "This is that is a phrase used for identification, where this points to something near and that points to something far, both referring to a single item, but it's also the name of a popular Canadian news satire show on CBC Radio. It highlights the basic grammar rule for demonstrative pronouns, contrasting proximity (this) with distance (that) in a humorous, news-style format, as seen in the show This Is That. "},
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
			if len(m.messages) > 0 && m.cursor < len(m.messages)-1 {
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
			sender := msgFromStyle.Render(message.From)
			rel := relativeTime(message.SentAt)
			var senderLine string
			if rel != "" {
				date := msgDateStyle.Render(rel)
				senderLine = lipgloss.JoinHorizontal(lipgloss.Left, sender, " ", date)
			} else {
				senderLine = sender
			}
			containerWidth := m.width
			if containerWidth <= 0 {
				containerWidth = 80
			}
			wrapWidth := containerWidth - 6
			if wrapWidth < 20 {
				wrapWidth = 20
			}
			wrapped := wordwrap.String(message.Text, wrapWidth)
			text := msgItemStyle.Render(wrapped)
			if i == m.cursor {
				text = msgSelectedStyle.Render(wrapped)
			}
			items[i] = lipgloss.JoinVertical(lipgloss.Left, senderLine, text)
		}
		list = lipgloss.JoinVertical(lipgloss.Left, items...)
	}

	instructions := msgInstructionStyle.Render("j/k or ↑/↓ to scroll • PgUp/PgDn to scroll • esc to menu • q to quit")

	content := lipgloss.JoinVertical(lipgloss.Center, header, identity, list, instructions)
	// Center horizontally and vertically in the terminal window
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
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
			Padding(0, 2).
			Margin(0, 0).
			Width(80).
			Border(lipgloss.NormalBorder()).
			BorderTop(false)
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
	msgFromStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#D1D6D6")).
			Padding(0, 2).
			Margin(0, 0).
			Foreground(lipgloss.Color("#323A38"))
	msgDateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)
)

func relativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	if d < 0 {
		d = -d
	}
	if d < time.Minute {
		sec := int(d.Seconds())
		if sec <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%ds ago", sec)
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 7 {
		return fmt.Sprintf("%dd ago", days)
	}
	if days < 30 {
		return fmt.Sprintf("%dw ago", days/7)
	}
	if days < 365 {
		return fmt.Sprintf("%dmo ago", days/30)
	}
	return fmt.Sprintf("%dy ago", days/365)
}
