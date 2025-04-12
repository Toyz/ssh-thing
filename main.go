package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/toyz/ssh-thing/config"
	"github.com/toyz/ssh-thing/tui"
	"github.com/toyz/ssh-thing/util"
)

func main() {
	keyBindsDefaultPath := os.Getenv("HOME") + "/.config/ssh-thing/keybinds.toml"
	if util.IsWindows() {
		keyBindsDefaultPath = os.Getenv("APPDATA") + "\\ssh-thing\\keybinds.toml"
	}

	serverConfigPath := flag.String("servers", "", "Path to the servers configuration file (default: ./servers.toml or ./config.toml)")
	keybindsPath := flag.String("keybinds", "", fmt.Sprint("Path to the keybinds configuration file (default: ", keyBindsDefaultPath, ")"))

	flag.Parse()

	if err := util.EnableVirtualTerminalProcessing(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to enable ANSI terminal support: %v\n", err)
	}

	cfg, err := config.LoadConfig(*serverConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading server config: %v\n", err)
		os.Exit(1)
	}

	if err := tui.StartTUI(cfg, *keybindsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
