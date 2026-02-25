package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
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
	loader spinner.Model
	form []textinput.Model
	formIndex int
	view view
	isLoading bool
	message string
	username string
	err error
}

func NewModel(ctx context.Context, vault *vault.Vault) Model {
	loginInput := textinput.New()
	loginInput.Placeholder = "Password"
	loginInput.EchoMode = textinput.EchoPassword
	loginInput.EchoCharacter = '*'

	editor := textinput.New()
	editor.Placeholder = "Service"

	columns := []table.Column{
		{ Title: "ID", Width: 0 },
		{ Title: "#", Width: 4 },
		{ Title: "Service", Width: 16 },
		{ Title: "Password", Width: 32 },
	}

	table := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithHeight(0),
	)

	loader := spinner.New()
	loader.Spinner = spinner.Dot

	var formInput textinput.Model
	form := make([]textinput.Model, 2)
	for i := range form {
		formInput = textinput.New()
		
		switch i {
		case 0:
			formInput.Placeholder = "Username"
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
		loginInput.Focus()
	} else {
		form[0].Focus()
	}

	return Model{
		ctx: ctx,
		vault: vault,
		view: currentView,	
		editor: editor,
		loginInput: loginInput,
		form: form,
		table: table,
		loader: loader,
		formIndex: 0,
		isLoading: false,	
		username: "",
		message: "",
		err: nil,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, m.loader.Tick, textinput.Blink)
}