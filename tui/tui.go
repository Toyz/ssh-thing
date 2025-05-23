package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/toyz/ssh-thing/config"
)

var program *tea.Program

func startClientCommands(model *Model) {
	for i, tab := range model.tabContents {
		if tab.Client == nil || tab.HasError {
			continue
		}

		// Clear the loading message before starting to stream output
		tab.ScrollView.Clear()

		if len(tab.Client.Config.Commands) > 0 {
			tab.Client.RunCommands(tab.Client.Config.Commands)
		}

		index := i
		go func() {
			for {
				select {
				case output := <-tab.Client.OutputChan:
					tab.ScrollView.Append(output)

					if program != nil {
						program.Send(updateContentMsg{index: index})
					}
				case err := <-tab.Client.ErrChan:
					tab.HandleError(err)
					if program != nil {
						program.Send(updateContentMsg{index: index})
					}
				}
			}
		}()
	}
}

func StartTUI(cfg *config.Config, keybindsPath string) error {
	model, err := NewModel(cfg, keybindsPath)
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	program = p

	startClientCommands(&model)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
