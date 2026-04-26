package shaders

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
)

//go:embed *.glsl
var embedded embed.FS

type Shader struct {
	Name     string
	Filename string
}

func List() []Shader {
	return []Shader{
		{Name: "disable"},
		{Name: "space", Filename: "starfield.glsl"},
		{Name: "retro terminal", Filename: "retro-terminal.glsl"},
		{Name: "under water", Filename: "underwater.glsl"},
		{Name: "drunk", Filename: "drunkard.glsl"},
		{Name: "glitchy", Filename: "glitchy.glsl"},
		{Name: "gradient", Filename: "gradient-shader.glsl"},
		{Name: "fireworks", Filename: "fireworks.glsl"},
		{Name: "matrix", Filename: "matrix.glsl"},
	}
}

func Install(shaderDir string, selected Shader) (string, error) {
	if selected.Filename == "" {
		return "", nil
	}
	if shaderDir == "" {
		return "", errors.New("could not resolve Ghostty shader directory")
	}

	data, err := embedded.ReadFile(selected.Filename)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(shaderDir, 0o755); err != nil {
		return "", err
	}

	shaderPath := filepath.Join(shaderDir, selected.Filename)
	if err := os.WriteFile(shaderPath, data, 0o644); err != nil {
		return "", err
	}

	return shaderPath, nil
}
