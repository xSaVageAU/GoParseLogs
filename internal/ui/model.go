package ui

import (
	"goparselogs/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// TUIModel wraps our models.Model to implement tea.Model interface
type TUIModel struct {
	state models.Model
}

// InitialModel creates a new TUIModel with initial state
func InitialModel() TUIModel {
	model := TUIModel{
		state: createInitialState(),
	}
	return model
}

// Init implements tea.Model
func (m TUIModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newState, cmd := Update(msg, m.state)
	m.state = newState
	return m, cmd
}

// View implements tea.Model
func (m TUIModel) View() string {
	return View(m.state)
}
