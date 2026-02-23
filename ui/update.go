package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type FailedMsg struct { Err error }

type RegisterSuccessMsg struct { Username string }
type LoginSuccessMsg struct { Username string }
type DashboardSuccessmsg struct { Passwords []table.Row }
type AddServiceSuccessMsg struct { }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case FailedMsg:
		m.err = msg.Err
		m.isLoading = false
		return m, nil
	case LoginSuccessMsg:
		m.view = dashboardView
		m.username = msg.Username
		m.isLoading = true
		return m, m.dashboard()
	case RegisterSuccessMsg:
		m.view = dashboardView
		m.username = msg.Username
		m.isLoading = true
		return m, m.dashboard()
	case DashboardSuccessmsg:
		rows := msg.Passwords
		m.table.SetHeight(len(rows) + 1)
		m.table.SetRows(rows)
		m.isLoading = false
		return m, nil
	case AddServiceSuccessMsg:
		m.table.SetHeight(m.table.Height() + 1)
		m.view = dashboardView
		m.isLoading	= true
		return m, m.dashboard()
	}

	if m.isLoading {
		return m, nil
	}

	switch m.view {
	case registerView:
		return m.registerUpdate(msg)
	case loginView:
		return m.loginUpdate(msg)
	case dashboardView:
		return m.dashboardUpdate(msg)
	case addServiceView:
		return m.addServiceUpdate(msg)
	}

	return m, nil
}

func (m Model) registerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyTab || msg.Type == tea.KeyShiftTab || msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
			keyType := msg.Type
			
			if keyType == tea.KeyEnter && m.formIndex == len(m.form) - 1 {
				username := m.form[0].Value()
				password := m.form[1].Value()
				m.err = nil
				return m, m.register(username, password)
			}

			if keyType == tea.KeyUp || keyType == tea.KeyShiftTab {
				m.formIndex--
			} else {
				m.formIndex++
			}

			if m.formIndex > len(m.form) - 1 {
				m.formIndex = 0
			} else if m.formIndex < 0 {
				m.formIndex = len(m.form) - 1
			}

			cmds := make([]tea.Cmd, len(m.form))
			for i := 0; i < len(m.form); i++ {
				if i == m.formIndex {
					cmds[i] = m.form[i].Focus()
					continue
				}

				m.form[i].Blur()
			}

			m.err = nil
			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) loginUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			password := m.loginInput.Value()
			m.loginInput.Reset()
			m.err = nil
			return m, m.login(password)
		}
	}
	
	var cmd tea.Cmd
	m.loginInput, cmd = m.loginInput.Update(msg)

	return m, cmd
}

func (m Model) dashboardUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	hasRows := len(m.table.Rows()) > 0

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if !hasRows {
				return m, nil
			}

			pass := m.table.SelectedRow()[2]
			clipboard.WriteAll(pass)
			m.message = fmt.Sprintf("%s copied to clipboard", pass)
			return m, nil
		} else if msg.Type == tea.KeyCtrlN {
			m.message = ""
			m.err = nil
			m.view = addServiceView
			return m, textinput.Blink
		}

		if !hasRows && (msg.Type == tea.KeyUp || msg.Type == tea.KeyDown) {
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	
	return m, cmd
}

func (m Model) addServiceUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			service := m.editor.Value()
			if service == "" {
				m.err = fmt.Errorf("Service name cannot be empty")
				return m, nil
			}
			m.err = nil
			m.editor.Reset()
			return m, m.addService(service)
		}
	}

	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(msg)

	return m, cmd
}

func (m Model) login(password string) (tea.Cmd) {
	return func() tea.Msg {
		if err := m.vault.Unlock(m.ctx, password); err != nil {
			return FailedMsg{Err: fmt.Errorf("Incorrect password")}
		}
		return LoginSuccessMsg{ Username: m.vault.Username() }
	}
}

func (m Model) register(username, password string) (tea.Cmd) {
	return func() tea.Msg {
		if err := m.vault.Setup(m.ctx, username, password); err != nil {
			return FailedMsg{Err: err}
		}
		return RegisterSuccessMsg{ Username: m.vault.Username()}
	}
}

func (m Model) dashboard() (tea.Cmd) {
	return func() tea.Msg {
		passwords ,err := m.vault.List(m.ctx)
		if err != nil {
			return FailedMsg{Err: err}
		}

		if len(passwords) == 0 {
			return FailedMsg{Err: fmt.Errorf("Nothing to display, try adding service")}
		}

		rowPasswords := []table.Row{}
		for i, pass := range passwords {
			row := []string{
				fmt.Sprintf("%d", i + 1),
				pass.Service,
				pass.Password,
			}
			rowPasswords = append(rowPasswords, row)
		}

		return DashboardSuccessmsg{Passwords: rowPasswords}
	}
}

func (m Model) addService(service string) tea.Cmd {
	return func() tea.Msg {
		if err := m.vault.AddService(m.ctx, service); err != nil{
			return FailedMsg{Err: err}
		}
		return AddServiceSuccessMsg{}
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.form))

	for i := range m.form {
		m.form[i], cmds[i] = m.form[i].Update(msg)
	}

	return tea.Batch(cmds...)
}