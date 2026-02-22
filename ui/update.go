package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type FailedMsg struct {Err error}

type RegisterSuccessMsg struct {}
type LoginSuccessMsg struct {}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case FailedMsg:
		m.err = msg.Err
		return m, nil
	case LoginSuccessMsg:
		m.view = dashboardView
		return m, nil
	case RegisterSuccessMsg:
		m.view = dashboardView
		return m, nil
	}

	switch m.view {
	case registerView:
		return m.registerUpdate(msg)
	case loginView:
		return m.loginUpdate(msg)
	case dashboardView:
		return m.dashboardUpdate(msg)
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

	return m, cmd}

func (m Model) dashboardUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) login(password string) (tea.Cmd) {
	return func() tea.Msg {
		if err := m.vault.Unlock(m.ctx, password); err != nil {
			return FailedMsg{Err: fmt.Errorf("Incorrect password")}
		}
		return LoginSuccessMsg{}
	}
}

func (m Model) register(username, password string) (tea.Cmd) {
	return func() tea.Msg {
		if err := m.vault.Setup(m.ctx, username, password); err != nil {
			return FailedMsg{Err: err}
		}
		return RegisterSuccessMsg{}
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.form))

	for i := range m.form {
		m.form[i], cmds[i] = m.form[i].Update(msg)
	}

	return tea.Batch(cmds...)
}