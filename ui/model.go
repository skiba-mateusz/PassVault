package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skiba-mateusz/PassVault/vault"
)

type Model struct {
	ctx context.Context
	vault *vault.Vault
	loginInput textinput.Model
	form []textinput.Model
	formIndex int
	cursorMode cursor.Mode
	view view
	username string
	err error
}

func NewModel(ctx context.Context, vault *vault.Vault) Model {
	loginInput := textinput.New()
	loginInput.Placeholder = "Password"
	loginInput.Focus()
	loginInput.EchoMode = textinput.EchoPassword
	loginInput.EchoCharacter = '*'

	var formInput textinput.Model
	form := make([]textinput.Model, 2)
	for i := range form {
		formInput = textinput.New()
		
		switch i {
		case 0:
			formInput.Placeholder = "Username"
			formInput.Focus()
		case 1:
			formInput.Placeholder = "Password"
			formInput.EchoMode = textinput.EchoPassword
			formInput.EchoCharacter = '*'
		}

		form[i] = formInput
	}

	currentView := registerView

	if vault.IsSetup(ctx) {
		currentView = loginView
	}

	return Model{
		ctx: ctx,
		vault: vault,
		view: currentView,	
		loginInput: loginInput,
		form: form,
		formIndex: 0,
		err: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}