package ghostty

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolvePathsUsesLegacyConfigWhenPresent(t *testing.T) {
	tempDir := t.TempDir()
	configHome := filepath.Join(tempDir, "xdg")
	configDir := filepath.Join(configHome, "ghostty")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	legacyConfig := filepath.Join(configDir, "config")
	if err := os.WriteFile(legacyConfig, []byte("font-size = 12\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", configHome)

	configPath, shaderDir := ResolvePaths()
	if configPath != legacyConfig {
		t.Fatalf("configPath = %q, want %q", configPath, legacyConfig)
	}

	wantShaderDir := filepath.Join(configDir, "lazyghost-shaders")
	if shaderDir != wantShaderDir {
		t.Fatalf("shaderDir = %q, want %q", shaderDir, wantShaderDir)
	}
}

func TestResolvePathsPrefersExistingModernConfig(t *testing.T) {
	tempDir := t.TempDir()
	configHome := filepath.Join(tempDir, "xdg")
	configDir := filepath.Join(configHome, "ghostty")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	modernConfig := filepath.Join(configDir, "config.ghostty")
	if err := os.WriteFile(modernConfig, []byte("font-size = 12\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", configHome)

	configPath, _ := ResolvePaths()
	if configPath != modernConfig {
		t.Fatalf("configPath = %q, want %q", configPath, modernConfig)
	}
}

func TestResolvePathsDarwinPrefersMacConfig(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-specific path resolution")
	}

	tempDir := t.TempDir()
	configHome := filepath.Join(tempDir, "xdg")
	xdgDir := filepath.Join(configHome, "ghostty")
	macDir := filepath.Join(tempDir, "Library", "Application Support", "com.mitchellh.ghostty")

	if err := os.MkdirAll(xdgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(macDir, 0o755); err != nil {
		t.Fatal(err)
	}

	xdgConfig := filepath.Join(xdgDir, "config.ghostty")
	macConfig := filepath.Join(macDir, "config.ghostty")
	if err := os.WriteFile(xdgConfig, []byte("theme = xdg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(macConfig, []byte("theme = mac\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", configHome)

	configPath, _ := ResolvePaths()
	if configPath != macConfig {
		t.Fatalf("configPath = %q, want %q", configPath, macConfig)
	}
}

func TestSaveShaderCreatesMissingConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "ghostty", "config.ghostty")
	shaderPath := "/tmp/lazyghost-shaders/matrix.glsl"

	if err := SaveShader(configPath, shaderPath); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	want := "custom-shader = " + shaderPath
	if string(got) != want {
		t.Fatalf("config contents = %q, want %q", string(got), want)
	}
}
