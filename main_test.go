package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveGhosttyPathsUsesLegacyConfigWhenPresent(t *testing.T) {
	oldHome := os.Getenv("HOME")
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	t.Cleanup(func() {
		t.Setenv("HOME", oldHome)
		if oldXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
			return
		}
		t.Setenv("XDG_CONFIG_HOME", oldXDG)
	})

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

	configPath, shaderDir := resolveGhosttyPaths()
	if configPath != legacyConfig {
		t.Fatalf("configPath = %q, want %q", configPath, legacyConfig)
	}

	wantShaderDir := filepath.Join(configDir, "lazyghost-shaders")
	if shaderDir != wantShaderDir {
		t.Fatalf("shaderDir = %q, want %q", shaderDir, wantShaderDir)
	}
}

func TestResolveGhosttyPathsPrefersExistingModernConfig(t *testing.T) {
	oldHome := os.Getenv("HOME")
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	t.Cleanup(func() {
		t.Setenv("HOME", oldHome)
		if oldXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
			return
		}
		t.Setenv("XDG_CONFIG_HOME", oldXDG)
	})

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

	configPath, _ := resolveGhosttyPaths()
	if configPath != modernConfig {
		t.Fatalf("configPath = %q, want %q", configPath, modernConfig)
	}
}

func TestSaveShaderCreatesMissingConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "ghostty", "config.ghostty")
	shaderPath := "/tmp/lazyghost-shaders/matrix.glsl"

	if err := saveShader(configPath, shaderPath); err != nil {
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

func TestResolveGhosttyPathsDarwinPrefersMacConfig(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-specific path resolution")
	}

	oldHome := os.Getenv("HOME")
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	t.Cleanup(func() {
		t.Setenv("HOME", oldHome)
		if oldXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
			return
		}
		t.Setenv("XDG_CONFIG_HOME", oldXDG)
	})

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

	configPath, _ := resolveGhosttyPaths()
	if configPath != macConfig {
		t.Fatalf("configPath = %q, want %q", configPath, macConfig)
	}
}

func TestGhosttyLinuxXDoToolArgsUseReloadShortcut(t *testing.T) {
	got := ghosttyLinuxXDoToolArgs()
	want := []string{
		"search", "--class", "ghostty",
		"windowactivate", "--sync",
		"key", "ctrl+shift+comma",
	}

	if len(got) != len(want) {
		t.Fatalf("len(args) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("args[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGhosttyLinuxServicesIncludeStableAndDebugUnits(t *testing.T) {
	want := []string{
		"app-com.mitchellh.ghostty.service",
		"app-com.mitchellh.ghostty-debug.service",
	}

	if len(ghosttyLinuxServices) != len(want) {
		t.Fatalf("len(services) = %d, want %d", len(ghosttyLinuxServices), len(want))
	}

	for i := range want {
		if ghosttyLinuxServices[i] != want[i] {
			t.Fatalf("services[%d] = %q, want %q", i, ghosttyLinuxServices[i], want[i])
		}
	}
}

func TestGhosttyLinuxSignalArgsUseSIGUSR2(t *testing.T) {
	got := ghosttyLinuxSignalArgs()
	want := []string{"-USR2", "-x", "ghostty"}

	if len(got) != len(want) {
		t.Fatalf("len(args) = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("args[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
