package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	registerView view = iota
	loginView
	dashboardView
	addPasswordView
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

func (m Model) View() string {
	result := fmt.Sprintf("%s\n\n", headerStyle.Render("PassVault - Your Personal Password Manager"))

	switch m.view {
	case registerView:
		result += m.registerView()
	case loginView:
		result += m.loginView()
	case dashboardView:
		result += m.dashboardView()
	}

	if m.err != nil {
		result += fmt.Sprintf("\n\n%s", m.err)
	}

	return fmt.Sprintf("%s\n\nesc - quit", result)
}

func (m Model) registerView() string {
	var builder strings.Builder

	for i := range m.form {
		builder.WriteString(m.form[i].View())
		if i < len(m.form) - 1 {
			builder.WriteRune('\n')
		}
	}

	return fmt.Sprintf("Register\n%v", builder.String())
}


func (m Model) loginView() string {
	return fmt.Sprintf("Login\n%v", m.loginInput.View())
}

func (m Model) dashboardView() string {
	return "Dashboard"
}

