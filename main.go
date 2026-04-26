package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
)

//go:embed shaders/*.glsl
var embeddedShaders embed.FS

type shader struct {
	name     string
	filename string
}

type model struct {
	configPath string
	shaderDir  string
	shaders    []shader
	cursor     int
	status     string
}

func initialModel() model {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	shaderDir := filepath.Join(home, ".config", "ghostty", "lazyghost-shaders")
	return model{
		configPath: filepath.Join(home, ".config", "ghostty", "config.ghostty"),
		shaderDir:  shaderDir,
		shaders: []shader{
			{name: "disable"},
			{name: "space", filename: "starfield.glsl"},
			{name: "retro terminal", filename: "retro-terminal.glsl"},
			{name: "under water", filename: "underwater.glsl"},
			{name: "drunk", filename: "drunkard.glsl"},
			{name: "glitchy", filename: "glitchy.glsl"},
			{name: "gradient", filename: "gradient-shader.glsl"},
			{name: "fireworks", filename: "fireworks.glsl"},
			{name: "matrix", filename: "matrix.glsl"},
		},
		cursor: 0,
	}
}

type shaderSavedMsg struct {
	name string
}

type shaderSaveFailedMsg struct {
	err error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.shaders)-1 {
				m.cursor++
			}

		// The "enter" key and the space bar toggle the selected state
		// for the item that the cursor is pointing at.
		case "enter", "space":
			selected := m.shaders[m.cursor]
			return m, func() tea.Msg {
				shaderPath, err := installShader(m.shaderDir, selected)
				if err != nil {
					return shaderSaveFailedMsg{err: err}
				}
				if err := saveShader(m.configPath, shaderPath); err != nil {
					return shaderSaveFailedMsg{err: err}
				}
				if err := reloadGhostty(); err != nil {
					return shaderSaveFailedMsg{err: fmt.Errorf("saved %s, but reload failed: %w", selected.name, err)}
				}
				return shaderSavedMsg{name: selected.name}
			}
		}
	case shaderSavedMsg:
		m.status = fmt.Sprintf("saved and reloaded: %s", msg.name)
	case shaderSaveFailedMsg:
		m.status = fmt.Sprintf("error: %v", msg.err)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	var s strings.Builder

	// The header
	s.WriteString("Change background\n\n")

	// Iterate over our choices
	for i, choice := range m.shaders {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		fmt.Fprintf(&s, "%s %s\n", cursor, choice.name)
	}

	// The footer
	s.WriteString("\nPress enter to save, q to quit.\n")
	if m.status != "" {
		fmt.Fprintf(&s, "%s\n", m.status)
	}

	// Send the UI for rendering
	return tea.NewView(s.String())
}

func saveShader(configPath, shaderPath string) error {
	if configPath == "" {
		return errors.New("could not resolve Ghostty config path")
	}

	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	line := "# custom-shader ="
	if shaderPath != "" {
		line = "custom-shader = " + shaderPath
	}

	lines := strings.Split(string(config), "\n")
	replaced := false
	for i, existing := range lines {
		trimmed := strings.TrimSpace(existing)
		trimmed = strings.TrimPrefix(trimmed, "#")
		trimmed = strings.TrimSpace(trimmed)
		if strings.HasPrefix(trimmed, "custom-shader") {
			lines[i] = line
			replaced = true
		}
	}

	if !replaced {
		lines = append(lines, line)
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0o644)
}

func installShader(shaderDir string, selected shader) (string, error) {
	if selected.filename == "" {
		return "", nil
	}
	if shaderDir == "" {
		return "", errors.New("could not resolve Ghostty shader directory")
	}

	data, err := embeddedShaders.ReadFile(filepath.Join("shaders", selected.filename))
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(shaderDir, 0o755); err != nil {
		return "", err
	}

	shaderPath := filepath.Join(shaderDir, selected.filename)
	if err := os.WriteFile(shaderPath, data, 0o644); err != nil {
		return "", err
	}

	return shaderPath, nil
}

func reloadGhostty() error {
	switch runtime.GOOS {
	case "darwin":
		return reloadGhosttyDarwin()
	case "linux":
		return reloadGhosttyLinux()
	default:
		return fmt.Errorf("automatic reload is not implemented on %s", runtime.GOOS)
	}
}

func reloadGhosttyDarwin() error {
	script := `
tell application "Ghostty" to activate
delay 0.1
tell application "System Events"
	keystroke "," using {command down, shift down}
end tell
`
	return exec.Command("osascript", "-e", script).Run()
}

func reloadGhosttyLinux() error {
	if err := reloadGhosttyLinuxDBus(); err == nil {
		return nil
	}

	if err := reloadGhosttyLinuxXDoTool(); err == nil {
		return nil
	}

	return errors.New("could not reload Ghostty; install gdbus/xdotool or bind reload_config and reload manually")
}

func reloadGhosttyLinuxDBus() error {
	var errs []string
	for _, appID := range []string{"com.mitchellh.ghostty", "com.mitchellh.ghostty-debug"} {
		objectPath := "/" + strings.ReplaceAll(appID, ".", "/")
		cmd := exec.Command(
			"gdbus",
			"call",
			"--session",
			"--dest", appID,
			"--object-path", objectPath,
			"--method", "org.gtk.Actions.Activate",
			"reload_config",
			"[]",
			"{}",
		)
		if output, err := cmd.CombinedOutput(); err == nil {
			return nil
		} else {
			errs = append(errs, strings.TrimSpace(string(output)))
		}
	}

	return errors.New(strings.Join(errs, "; "))
}

func reloadGhosttyLinuxXDoTool() error {
	cmd := exec.Command(
		"xdotool",
		"search", "--class", "ghostty",
		"windowactivate", "--sync",
		"key", "super+shift+comma",
	)
	return cmd.Run()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
