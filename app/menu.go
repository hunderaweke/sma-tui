package app

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	StateMenu = iota
	StateIdentities
	StateMessages
	MessagesScreen = "messages"
)

type MenuModel struct {
	State    int
	cursor   int
	width    int
	height   int
	identity string
}
type ChangeScreenMsg struct {
	NewScreen string
}

func NewModel(a App) MenuModel {
	return MenuModel{State: StateMenu, identity: a.Config.DefaultIdentity}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}
func (m MenuModel) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	ascii := logoStyle.Render(logo)
	title := titleStyle.Render(title)
	subTitle := subTitleStyle.Render(subTitle)
	identitiesNumber := numberStyle.Render(identitiesNumber)
	identitiesText := instructionStyle.Render(identitiesText)
	msgsNumber := numberStyle.Render(msgsNumber)
	msgsText := instructionStyle.Render(msgsText)
	analytics := lipgloss.JoinHorizontal(lipgloss.Center, identitiesNumber, identitiesText, msgsNumber, msgsText)
	instructionLine := bottomInstruction.Render(instructions)
	btn1Style := buttonStyle
	if m.cursor == 0 {
		btn1Style = activeButtonStyle
	}
	btn1 := btn1Style.Render("  Identities")

	btn2Style := buttonStyle
	if m.cursor == 1 {
		btn2Style = activeButtonStyle
	}
	btn2 := btn2Style.Render("󰶌 Your Messages")
	identity := identityStyle.Render(m.identity)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, btn1, btn2)
	content := lipgloss.JoinVertical(lipgloss.Center, ascii, title, subTitle, analytics, buttons, identity, instructionLine)
	page := uiWrapperStyle.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, page)
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		case "i":
			m.cursor = 1
			return m, nil
		case "m":
			m.cursor = 0
			return m, nil
		case "enter":
			return m, m.switchScreen()
		}
	case ChangeScreenMsg:
		switch msg.NewScreen {
		case MessagesScreen:
			log.Println("Messages Screen")
			return NewMessagesModel(m.identity), nil
		}
	}
	return m, nil
}

func (m MenuModel) switchScreen() tea.Cmd {
	if m.State == StateMenu {
		switch m.cursor {
		case 0:
			return func() tea.Msg {
				return ChangeScreenMsg{NewScreen: MessagesScreen}
			}
		case 1:
			log.Println("Identities")
		}
	}
	return func() tea.Msg {
		return "Something"
	}
}

var (
	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Bold(true)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#D1D6D6")).
			Padding(1, 2)

	subTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 2).
			Italic(true)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	instructionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true)
	uiWrapperStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Padding(1, 2).
			MarginTop(1).
			Margin(0, 2)
	bottomInstruction = instructionStyle.MarginTop(10).AlignVertical(lipgloss.Bottom)
	buttonStyle       = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#323A38")).
				Background(lipgloss.Color("#D1D6D6")).
				Padding(1, 3).
				MarginTop(1).
				MarginRight(2)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#D1D6D6")).
				Background(lipgloss.Color("#323A38")).
				Bold(true)

	numberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Bold(true)
	identityStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Border(lipgloss.NormalBorder())
)
