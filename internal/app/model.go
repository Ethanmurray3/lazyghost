package app

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/Ethanmurray3/lazyghost/internal/ghostty"
	shaderpkg "github.com/Ethanmurray3/lazyghost/internal/shaders"
)

type model struct {
	configPath string
	shaderDir  string
	shaders    []shaderpkg.Shader
	cursor     int
	status     string
}

type shaderSavedMsg struct {
	name string
}

type shaderSaveFailedMsg struct {
	err error
}

func initialModel() model {
	configPath, shaderDir := ghostty.ResolvePaths()
	return model{
		configPath: configPath,
		shaderDir:  shaderDir,
		shaders:    shaderpkg.List(),
		cursor:     0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.shaders)-1 {
				m.cursor++
			}
		case "enter", "space":
			selected := m.shaders[m.cursor]
			return m, func() tea.Msg {
				shaderPath, err := shaderpkg.Install(m.shaderDir, selected)
				if err != nil {
					return shaderSaveFailedMsg{err: err}
				}
				if err := ghostty.SaveShader(m.configPath, shaderPath); err != nil {
					return shaderSaveFailedMsg{err: err}
				}
				if err := ghostty.Reload(); err != nil {
					return shaderSaveFailedMsg{err: fmt.Errorf("saved %s, but reload failed: %w", selected.Name, err)}
				}
				return shaderSavedMsg{name: selected.Name}
			}
		}
	case shaderSavedMsg:
		m.status = fmt.Sprintf("saved and reloaded: %s", msg.name)
	case shaderSaveFailedMsg:
		m.status = fmt.Sprintf("error: %v", msg.err)
	}

	return m, nil
}

func (m model) View() tea.View {
	var s strings.Builder

	s.WriteString("Change background\n\n")

	for i, choice := range m.shaders {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		fmt.Fprintf(&s, "%s %s\n", cursor, choice.Name)
	}

	s.WriteString("\nPress enter to save, q to quit.\n")
	if m.status != "" {
		fmt.Fprintf(&s, "%s\n", m.status)
	}

	return tea.NewView(s.String())
}

func Run() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
