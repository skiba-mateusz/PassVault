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
    headerStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("235")).
        Foreground(lipgloss.Color("230")).
        Bold(true).
        PaddingLeft(1).
        PaddingRight(1)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)

    messageStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("49"))

    footerStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("235")).
        Foreground(lipgloss.Color("230")).
        PaddingLeft(1).
        PaddingRight(1).
		Bold(true)
)

var tableStyles = func() table.Styles {
    s := table.DefaultStyles()

    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        BorderBottom(true).
        Foreground(lipgloss.Color("230")).
        Bold(true)

    s.Selected = s.Selected.
        Foreground(lipgloss.Color("193")).
		Bold(true)

    return s
}()

func (m Model) View() string {
    header := headerStyle.Render("PassVault - Your Personal Password Manager")

	if m.username != "" {
    	header = headerStyle.Render(fmt.Sprintf("PassVault - Your Personal Password Manager | welcome, %s", m.username))
	}

    var body string
    footer := "Esc - quit"

	if m.message != "" {
		body += messageStyle.Render(m.message) + "\n\n"
	}
    
    if m.isLoading {
        body += m.loader.View() + " Loading..."
    } else {
        switch m.view {
        case registerView:
            body += m.registerView()
        case loginView:
            body += m.loginView()
        case dashboardView:
            body += m.dashboardView()
            footer += " | Ctrl+N - add | Ctrl+X - del"
        case addServiceView:
            body += m.addServiceView()
            footer += " | Backspace - back"
        }
    }

    if m.err != nil {
        body += "\n\n" + errorStyle.Render(m.err.Error())
    }

    return lipgloss.JoinVertical(lipgloss.Left,
        header,
        "",
        body,
        "",
        footerStyle.Render(footer),
    )
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

	if len(m.table.Rows()) == 0 {
		return "No services found, try adding one"
	}

	return m.table.View()
}

func (m Model) addServiceView() string {
	return fmt.Sprintf("Add service\n%s", m.editor.View())
}