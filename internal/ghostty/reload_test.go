package ghostty

import "testing"

func TestLinuxXDoToolArgsUseReloadShortcut(t *testing.T) {
	got := linuxXDoToolArgs()
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

func TestLinuxServicesIncludeStableAndDebugUnits(t *testing.T) {
	want := []string{
		"app-com.mitchellh.ghostty.service",
		"app-com.mitchellh.ghostty-debug.service",
	}

	if len(linuxServices) != len(want) {
		t.Fatalf("len(services) = %d, want %d", len(linuxServices), len(want))
	}

	for i := range want {
		if linuxServices[i] != want[i] {
			t.Fatalf("services[%d] = %q, want %q", i, linuxServices[i], want[i])
		}
	}
}

func TestLinuxSignalArgsUseSIGUSR2(t *testing.T) {
	got := linuxSignalArgs()
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
