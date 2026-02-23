package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	registerView view = iota
	loginView
	dashboardView
	addServiceView
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

var tableStyles = func() table.Styles {
    s := table.DefaultStyles()

    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        BorderBottom(true).
        Bold(false)

    s.Selected = s.Selected.
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("57")).
        Bold(false)

    return s
}()

func (m Model) View() string {
	result := fmt.Sprintf("%s\n\n", headerStyle.Render(fmt.Sprintf("PassVault - Your Personal Password Manager | %s", m.username)))
	footer := "esc - quit"

	switch m.view {
	case registerView:
		result += m.registerView()
	case loginView:
		result += m.loginView()
	case dashboardView:
		result += m.dashboardView()
		footer += " | ctrl+n - add service"
	case addServiceView:
		result += m.addServiceView()
	}

	if m.err != nil {
		result += fmt.Sprintf("\n\n%s", m.err)
	}

	return fmt.Sprintf("%s\n\n%s", result, footer)
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
	m.table.SetStyles(tableStyles)
	return fmt.Sprintf("%s\n\n%s", m.message, m.table.View())
}

func (m Model) addServiceView() string {
	return fmt.Sprintf("Add service\n%s", m.editor.View())
}