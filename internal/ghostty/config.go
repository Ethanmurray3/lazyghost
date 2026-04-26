package ghostty

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolvePaths() (configPath, shaderDir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", ""
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(home, ".config")
	}

	xdgDir := filepath.Join(configHome, "ghostty")
	shaderDir = filepath.Join(xdgDir, "lazyghost-shaders")

	xdgCandidates := []string{
		filepath.Join(xdgDir, "config.ghostty"),
		filepath.Join(xdgDir, "config"),
	}

	if runtime.GOOS == "darwin" {
		macDir := filepath.Join(home, "Library", "Application Support", "com.mitchellh.ghostty")
		for _, path := range []string{
			filepath.Join(macDir, "config.ghostty"),
			filepath.Join(macDir, "config"),
		} {
			if fileExists(path) {
				return path, shaderDir
			}
		}
	}

	for _, path := range xdgCandidates {
		if fileExists(path) {
			return path, shaderDir
		}
	}

	return xdgCandidates[0], shaderDir
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func SaveShader(configPath, shaderPath string) error {
	if configPath == "" {
		return errors.New("could not resolve Ghostty config path")
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	config, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		config = nil
	}

	line := "# custom-shader ="
	if shaderPath != "" {
		line = "custom-shader = " + shaderPath
	}

	lines := strings.Split(string(config), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
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
