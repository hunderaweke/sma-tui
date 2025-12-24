package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunderaweke/sma-tui/config"
)

var (
	logo = `
███████╗    ███╗   ███╗     █████╗ 
██╔════╝    ████╗ ████║    ██╔══██╗
███████╗    ██╔████╔██║    ███████║
╚════██║    ██║╚██╔╝██║    ██╔══██║
███████║ ██ ██║ ╚═╝ ██║ ██ ██║  ██║
╚══════╝    ╚═╝     ╚═╝    ╚═╝  ╚═╝.et
`
	title            = "Welcome to S.M.A"
	subTitle         = "Send Message Anon"
	identitiesText   = " Identities"
	identitiesNumber = "5720"
	msgsNumber       = " 3147"
	msgsText         = " Msgs and Counting..."
	instructions     = "i: Identities • m: Messages • Esc: Close • Enter: Continue • q/Ctrl+C: Quit"
	identity         = "uSOeAQEuQxCW"
)

type model struct {
	cursor int
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) View() string {
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
	identity := identityStyle.Render(identity)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, btn1, btn2)
	content := lipgloss.JoinVertical(lipgloss.Center, ascii, title, subTitle, analytics, buttons, identity, instructionLine)
	page := uiWrapperStyle.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, page)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.cursor = 0
			return m, nil
		case "m":
			m.cursor = 1
			return m, nil
		case "enter":
			fmt.Println("Transitioning...")
		}
	}
	return m, nil
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
				Background(lipgloss.Color("#323A38")). // Bright Cyan
				Bold(true)

	numberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Bold(true)
	identityStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D6D6")).
			Border(lipgloss.NormalBorder())
)

func main() {
	// p := tea.NewProgram(model{}, tea.WithAltScreen()) // WithAltScreen hides your terminal history
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	os.Exit(1)
	// }
	// c, err := NewConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = c.Save("config.json"); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Created Config")
	c, err := config.Load("config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c, _ = config.New()
			c.Save("config.json")
		} else {
			log.Fatal(err)
		}
	}
	log.Println(c)
}
