package ghostty

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var linuxServices = []string{
	"app-com.mitchellh.ghostty.service",
	"app-com.mitchellh.ghostty-debug.service",
}

func Reload() error {
	switch runtime.GOOS {
	case "darwin":
		return reloadDarwin()
	case "linux":
		return reloadLinux()
	default:
		return fmt.Errorf("automatic reload is not implemented on %s", runtime.GOOS)
	}
}

func reloadDarwin() error {
	script := `
tell application "Ghostty" to activate
delay 0.1
tell application "System Events"
	keystroke "," using {command down, shift down}
end tell
`
	return exec.Command("osascript", "-e", script).Run()
}

func reloadLinux() error {
	var errs []string

	if err := reloadLinuxSystemd(); err == nil {
		return nil
	} else if err.Error() != "" {
		errs = append(errs, "systemd: "+err.Error())
	}

	if err := reloadLinuxDBus(); err == nil {
		return nil
	} else if err.Error() != "" {
		errs = append(errs, "dbus: "+err.Error())
	}

	if err := reloadLinuxSignal(); err == nil {
		return nil
	} else if err.Error() != "" {
		errs = append(errs, "signal: "+err.Error())
	}

	if err := reloadLinuxXDoTool(); err == nil {
		return nil
	} else if err.Error() != "" {
		errs = append(errs, "xdotool: "+err.Error())
	}

	return fmt.Errorf("could not reload Ghostty (%s)", strings.Join(errs, "; "))
}

func reloadLinuxSystemd() error {
	var errs []string
	for _, service := range linuxServices {
		cmd := exec.Command("systemctl", "reload", "--user", service)
		if output, err := cmd.CombinedOutput(); err == nil {
			return nil
		} else {
			errs = append(errs, strings.TrimSpace(string(output)))
		}
	}

	return errors.New(strings.Join(errs, "; "))
}

func reloadLinuxDBus() error {
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

func reloadLinuxXDoTool() error {
	cmd := exec.Command("xdotool", linuxXDoToolArgs()...)
	return cmd.Run()
}

func reloadLinuxSignal() error {
	cmd := exec.Command("pkill", linuxSignalArgs()...)
	if output, err := cmd.CombinedOutput(); err == nil {
		return nil
	} else {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return errors.New(msg)
	}
}

func linuxSignalArgs() []string {
	return []string{"-USR2", "-x", "ghostty"}
}

func linuxXDoToolArgs() []string {
	return []string{
		"search", "--class", "ghostty",
		"windowactivate", "--sync",
		"key", "ctrl+shift+comma",
	}
}
