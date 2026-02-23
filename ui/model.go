package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skiba-mateusz/PassVault/vault"
)

type Model struct {
	ctx context.Context
	vault *vault.Vault
	editor textinput.Model
	loginInput textinput.Model
	table table.Model
	form []textinput.Model
	formIndex int
	view view
	message string
	username string
	err error
}

func NewModel(ctx context.Context, vault *vault.Vault) Model {
	loginInput := textinput.New()
	loginInput.Placeholder = "Password"
	loginInput.Focus()
	loginInput.EchoMode = textinput.EchoPassword
	loginInput.EchoCharacter = '*'

	editor := textinput.New()
	editor.Placeholder = "Service"
	editor.Focus()

	columns := []table.Column{
		{ Title: "#", Width: 4 },
		{ Title: "Service", Width: 16 },
		{ Title: "Password", Width: 32 },
	}

	table := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(0),
	)

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
		editor: editor,
		loginInput: loginInput,
		form: form,
		table: table,
		formIndex: 0,
		username: "",
		message: "",
		err: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}