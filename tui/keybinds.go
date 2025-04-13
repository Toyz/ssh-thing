package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/pelletier/go-toml/v2"
)

type KeyBindingsMap struct {
	Up                []string `toml:"up"`
	Down              []string `toml:"down"`
	Left              []string `toml:"left"`
	Right             []string `toml:"right"`
	PageUp            []string `toml:"pageUp"`
	PageDown          []string `toml:"pageDown"`
	Home              []string `toml:"home"`
	End               []string `toml:"end"`
	Quit              []string `toml:"quit"`
	ToggleColor       []string `toml:"toggleColor"`
	ResetScroll       []string `toml:"resetScroll"`
	TabNext           []string `toml:"tabNext"`
	TabPrev           []string `toml:"tabPrev"`
	ClearBuffer       []string `toml:"clearBuffer"`
	ToggleWordWrap    []string `toml:"toggleWordWrap"`
	ToggleTabPosition []string `toml:"toggleTabPosition"`
}

type KeyBindingsConfig struct {
	Keybinds KeyBindingsMap `toml:"keybinds"`
}

type KeyMap struct {
	Up                key.Binding
	Down              key.Binding
	Left              key.Binding
	Right             key.Binding
	PageUp            key.Binding
	PageDown          key.Binding
	Home              key.Binding
	End               key.Binding
	Quit              key.Binding
	ToggleColor       key.Binding
	ResetScroll       key.Binding
	TabNext           key.Binding
	TabPrev           key.Binding
	ClearBuffer       key.Binding
	ToggleWordWrap    key.Binding
	ToggleTabPosition key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Up, k.Down, k.ToggleColor, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right, k.TabPrev, k.TabNext},
		{k.Up, k.Down, k.PageUp, k.PageDown, k.Home, k.End, k.ResetScroll},
		{k.ToggleColor, k.ToggleWordWrap, k.ClearBuffer, k.Quit},
	}
}

var bindingDescriptions = map[string]string{
	"up":                "scroll up",
	"down":              "scroll down",
	"left":              "previous tab",
	"right":             "next tab",
	"pageUp":            "page up",
	"pageDown":          "page down",
	"home":              "scroll to top",
	"end":               "scroll to bottom",
	"quit":              "quit",
	"toggleColor":       "toggle colors",
	"resetScroll":       "reset scroll",
	"tabNext":           "next tab",
	"tabPrev":           "previous tab",
	"clearBuffer":       "clear buffer",
	"toggleWordWrap":    "toggle word wrap",
	"toggleTabPosition": "toggle tab position",
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "previous tab"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next tab"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "scroll to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "scroll to bottom"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		ToggleColor: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "toggle colors"),
		),
		ResetScroll: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset scroll"),
		),
		TabNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		TabPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),
		ClearBuffer: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear buffer"),
		),
		ToggleWordWrap: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "toggle word wrap"),
		),
		ToggleTabPosition: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle tab position"),
		),
	}
}

func DefaultKeyBindingsMap() KeyBindingsMap {
	return KeyBindingsMap{
		Up:                []string{"up", "k"},
		Down:              []string{"down", "j"},
		Left:              []string{"left", "h"},
		Right:             []string{"right", "l"},
		PageUp:            []string{"pgup"},
		PageDown:          []string{"pgdown"},
		Home:              []string{"home"},
		End:               []string{"end", "G"},
		Quit:              []string{"q", "ctrl+c"},
		ToggleColor:       []string{"c"},
		ResetScroll:       []string{"r"},
		TabNext:           []string{"tab"},
		TabPrev:           []string{"shift+tab"},
		ClearBuffer:       []string{"ctrl+l"},
		ToggleWordWrap:    []string{"w"},
		ToggleTabPosition: []string{"p"},
	}
}

func (m KeyBindingsMap) ToKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys(m.Up...),
			key.WithHelp(getHelpPrefix(m.Up), bindingDescriptions["up"]),
		),
		Down: key.NewBinding(
			key.WithKeys(m.Down...),
			key.WithHelp(getHelpPrefix(m.Down), bindingDescriptions["down"]),
		),
		Left: key.NewBinding(
			key.WithKeys(m.Left...),
			key.WithHelp(getHelpPrefix(m.Left), bindingDescriptions["left"]),
		),
		Right: key.NewBinding(
			key.WithKeys(m.Right...),
			key.WithHelp(getHelpPrefix(m.Right), bindingDescriptions["right"]),
		),
		PageUp: key.NewBinding(
			key.WithKeys(m.PageUp...),
			key.WithHelp(getHelpPrefix(m.PageUp), bindingDescriptions["pageUp"]),
		),
		PageDown: key.NewBinding(
			key.WithKeys(m.PageDown...),
			key.WithHelp(getHelpPrefix(m.PageDown), bindingDescriptions["pageDown"]),
		),
		Home: key.NewBinding(
			key.WithKeys(m.Home...),
			key.WithHelp(getHelpPrefix(m.Home), bindingDescriptions["home"]),
		),
		End: key.NewBinding(
			key.WithKeys(m.End...),
			key.WithHelp(getHelpPrefix(m.End), bindingDescriptions["end"]),
		),
		Quit: key.NewBinding(
			key.WithKeys(m.Quit...),
			key.WithHelp(getHelpPrefix(m.Quit), bindingDescriptions["quit"]),
		),
		ToggleColor: key.NewBinding(
			key.WithKeys(m.ToggleColor...),
			key.WithHelp(getHelpPrefix(m.ToggleColor), bindingDescriptions["toggleColor"]),
		),
		ResetScroll: key.NewBinding(
			key.WithKeys(m.ResetScroll...),
			key.WithHelp(getHelpPrefix(m.ResetScroll), bindingDescriptions["resetScroll"]),
		),
		TabNext: key.NewBinding(
			key.WithKeys(m.TabNext...),
			key.WithHelp(getHelpPrefix(m.TabNext), bindingDescriptions["tabNext"]),
		),
		TabPrev: key.NewBinding(
			key.WithKeys(m.TabPrev...),
			key.WithHelp(getHelpPrefix(m.TabPrev), bindingDescriptions["tabPrev"]),
		),
		ClearBuffer: key.NewBinding(
			key.WithKeys(m.ClearBuffer...),
			key.WithHelp(getHelpPrefix(m.ClearBuffer), bindingDescriptions["clearBuffer"]),
		),
		ToggleWordWrap: key.NewBinding(
			key.WithKeys(m.ToggleWordWrap...),
			key.WithHelp(getHelpPrefix(m.ToggleWordWrap), bindingDescriptions["toggleWordWrap"]),
		),
		ToggleTabPosition: key.NewBinding(
			key.WithKeys(m.ToggleTabPosition...),
			key.WithHelp(getHelpPrefix(m.ToggleTabPosition), bindingDescriptions["toggleTabPosition"]),
		),
	}
}

func getHelpPrefix(keys []string) string {
	if len(keys) == 1 {
		return keys[0]
	} else if len(keys) == 2 {
		return fmt.Sprintf("%s/%s", keys[0], keys[1])
	} else if len(keys) > 2 {
		return fmt.Sprintf("%s/...", keys[0])
	}
	return ""
}

func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	sshThingConfigDir := filepath.Join(configDir, "ssh-thing")
	err = os.MkdirAll(sshThingConfigDir, 0755)
	if err != nil {
		return "", err
	}

	return sshThingConfigDir, nil
}

func SaveKeyBindings(bindings KeyBindingsMap, filePath string) error {
	if filePath == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}

		filePath = filepath.Join(configDir, "keybinds.toml")
	}

	config := KeyBindingsConfig{
		Keybinds: bindings,
	}

	tomlData, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal keybinds config: %w", err)
	}

	err = os.WriteFile(filePath, tomlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write keybinds file %s: %w", filePath, err)
	}

	return nil
}

func LoadKeyBindings(filePath string) (KeyBindingsMap, error) {
	defaultBindings := DefaultKeyBindingsMap()

	if filePath == "" {
		configDir, err := GetConfigDir()
		if err != nil {
			return defaultBindings, err
		}

		filePath = filepath.Join(configDir, "keybinds.toml")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := SaveKeyBindings(defaultBindings, ""); err != nil {
				return defaultBindings, fmt.Errorf("failed to create default keybinds file: %w", err)
			}
			return defaultBindings, nil
		}
	}

	tomlData, err := os.ReadFile(filePath)
	if err != nil {
		return defaultBindings, fmt.Errorf("failed to read keybinds file %s: %w", filePath, err)
	}

	var config KeyBindingsConfig
	if err := toml.Unmarshal(tomlData, &config); err != nil {
		return defaultBindings, fmt.Errorf("failed to parse keybinds file %s: %w", filePath, err)
	}

	if len(config.Keybinds.Up) == 0 {
		config.Keybinds.Up = defaultBindings.Up
	}
	if len(config.Keybinds.Down) == 0 {
		config.Keybinds.Down = defaultBindings.Down
	}
	if len(config.Keybinds.Left) == 0 {
		config.Keybinds.Left = defaultBindings.Left
	}
	if len(config.Keybinds.Right) == 0 {
		config.Keybinds.Right = defaultBindings.Right
	}
	if len(config.Keybinds.PageUp) == 0 {
		config.Keybinds.PageUp = defaultBindings.PageUp
	}
	if len(config.Keybinds.PageDown) == 0 {
		config.Keybinds.PageDown = defaultBindings.PageDown
	}
	if len(config.Keybinds.Home) == 0 {
		config.Keybinds.Home = defaultBindings.Home
	}
	if len(config.Keybinds.End) == 0 {
		config.Keybinds.End = defaultBindings.End
	}
	if len(config.Keybinds.Quit) == 0 {
		config.Keybinds.Quit = defaultBindings.Quit
	}
	if len(config.Keybinds.ToggleColor) == 0 {
		config.Keybinds.ToggleColor = defaultBindings.ToggleColor
	}
	if len(config.Keybinds.ResetScroll) == 0 {
		config.Keybinds.ResetScroll = defaultBindings.ResetScroll
	}
	if len(config.Keybinds.TabNext) == 0 {
		config.Keybinds.TabNext = defaultBindings.TabNext
	}
	if len(config.Keybinds.TabPrev) == 0 {
		config.Keybinds.TabPrev = defaultBindings.TabPrev
	}
	if len(config.Keybinds.ClearBuffer) == 0 {
		config.Keybinds.ClearBuffer = defaultBindings.ClearBuffer
	}
	if len(config.Keybinds.ToggleWordWrap) == 0 {
		config.Keybinds.ToggleWordWrap = defaultBindings.ToggleWordWrap
	}
	if len(config.Keybinds.ToggleTabPosition) == 0 {
		config.Keybinds.ToggleTabPosition = defaultBindings.ToggleTabPosition
	}

	return config.Keybinds, nil
}

func LoadKeyMap(filePath string) KeyMap {
	keyBindings, err := LoadKeyBindings(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load keybinds, using defaults: %v\n", err)
		return DefaultKeyMap()
	}

	return keyBindings.ToKeyMap()
}
