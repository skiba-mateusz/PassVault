package ui

import (
	"fmt"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type FailedMsg struct { Err error }

type RegisterSuccessMsg struct { Username string }
type LoginSuccessMsg struct { Username string }
type DashboardSuccessMsg struct { Passwords []table.Row }
type AddServiceSuccessMsg struct { }
type DeleteServiceSuccessMsg struct { }

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

		m.resetStatus()
	}

	switch msg := msg.(type) {
	case FailedMsg:
		m.err = msg.Err
		m.isLoading = false
		return m, nil	
	case LoginSuccessMsg:
		m.username = msg.Username
		return m.moveToDashboard()
	case RegisterSuccessMsg:
		m.username = msg.Username
		return m.moveToDashboard()
	case DashboardSuccessMsg:
		rows := msg.Passwords
		m.table.SetRows(rows)
		
		height := len(rows) + 1
		if height < 3 {
			height = 3
		}
		m.table.SetHeight(height)
        m.table.SetCursor(0)
		
		m.isLoading = false
		m.table.Focus()
		return m, nil
	case AddServiceSuccessMsg:
		return m.moveToDashboard()
	case DeleteServiceSuccessMsg:
		return m.moveToDashboard()
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
				if username == "" || password == "" {
					m.err = fmt.Errorf("Username and password cannot be empty")
					return m, nil
				}
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
			return m, m.login(password)
		}
	}
	
	var cmd tea.Cmd
	m.loginInput, cmd = m.loginInput.Update(msg)

	return m, cmd
}

func (m Model) dashboardUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	hasRows := len(m.table.Rows()) > 0

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if !hasRows {
				return m, nil
			}

			pass := m.table.SelectedRow()[3]
			clipboard.WriteAll(pass)
			m.message = fmt.Sprintf("%s copied to clipboard", pass)
			return m, nil
		} else if msg.Type == tea.KeyCtrlN {
			m.view = addServiceView
			return m, m.editor.Focus()
		} else if msg.Type == tea.KeyCtrlX {
			idStr := m.table.SelectedRow()[0]
			id, _ := strconv.Atoi(idStr)
			return m, m.deleteService(int64(id))
		}

		if !hasRows && (msg.Type == tea.KeyUp || msg.Type == tea.KeyDown) {
			return m, nil
		}
	}

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
			m.editor.Reset()
			return m, m.addService(service)
		} else if msg.Type == tea.KeyBackspace {
			if len(m.editor.Value()) <= 0 {
				m.view = dashboardView
			}
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

		rowPasswords := []table.Row{}
		for i, pass := range passwords {
			row := []string{
				fmt.Sprintf("%d", pass.ID),
				fmt.Sprintf("%d", i + 1),
				pass.Service,
				pass.Password,
			}
			rowPasswords = append(rowPasswords, row)
		}

		return DashboardSuccessMsg{Passwords: rowPasswords}
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

func (m Model) deleteService(id int64) tea.Cmd {
	return func() tea.Msg {
		if err := m.vault.DeleteService(m.ctx, id); err != nil {
			return FailedMsg{Err: err}
		}
		return DeleteServiceSuccessMsg{}
	}
}

func (m Model) moveToDashboard() (tea.Model, tea.Cmd) {
	m.view = dashboardView
	m.isLoading = true
	m.table.Focus()
	return m, m.dashboard()
}

func (m *Model) resetStatus() {
	m.err = nil
	m.message = ""	
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.form))

	for i := range m.form {
		m.form[i], cmds[i] = m.form[i].Update(msg)
	}

	return tea.Batch(cmds...)
}