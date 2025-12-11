package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toyz/ssh-thing/config"
	"github.com/toyz/ssh-thing/ssh"
	"github.com/toyz/ssh-thing/tui/components"
)

type Model struct {
	tabs         []string
	tabContents  []*components.TabContent
	activeTab    int
	tabOffset    int
	quitting     bool
	width        int
	height       int
	colorize     bool
	verticalTabs bool
	keys         KeyMap
	help         help.Model
	statusBar    *components.StatusBar
	config       *config.Config
}

func NewModel(cfg *config.Config, keybindsPath string) (Model, error) {
	var tabs []string
	var tabContents []*components.TabContent

	for _, server := range cfg.Servers {
		tabs = append(tabs, server.Name)

		tab := components.NewTabContent(server.Name)
		tab.ScrollView.Append("Connecting...")
		tabContents = append(tabContents, tab)
	}

	helpModel := help.New()
	helpModel.Styles.ShortKey = components.HelpShortKey
	helpModel.Styles.ShortDesc = components.HelpShortDesc
	helpModel.Styles.ShortSeparator = components.HelpShortSeparator
	helpModel.Styles.FullKey = components.HelpFullKey
	helpModel.Styles.FullDesc = components.HelpFullDesc
	helpModel.Styles.FullSeparator = components.HelpFullSeparator

	return Model{
		tabs:         tabs,
		tabContents:  tabContents,
		activeTab:    0,
		quitting:     false,
		colorize:     false,
		verticalTabs: true,
		keys:         LoadKeyMap(keybindsPath),
		help:         helpModel,
		statusBar:    components.NewStatusBar(),
		config:       cfg,
	}, nil
}

type updateContentMsg struct {
	index int
}

type sshConnectionMsg struct {
	index  int
	client *ssh.Client
	err    error
}

type connectSSHCmd struct {
	index  int
	server *config.SSHServer
}

func connectSSHClient(index int, server *config.SSHServer) tea.Cmd {
	return func() tea.Msg {
		client, err := ssh.NewClient(server)
		return sshConnectionMsg{
			index:  index,
			client: client,
			err:    err,
		}
	}
}

var (
	infoPattern      = regexp.MustCompile(`(?i)(\[?INFO\]?[:|\s]+)|(INFO\s*-\s*)|(INFO\s{2,})|\bINFO\b`)
	errorPattern     = regexp.MustCompile(`(?i)(\[?ERROR\]?[:|\s]+)|(ERROR\s*-\s*)|(ERROR\s{2,})|\bERROR\b`)
	warnPattern      = regexp.MustCompile(`(?i)(\[?WARN(ING)?\]?[:|\s]+)|(\[?WARNING\]?[:|\s]+)|(WARN(ING)?\s*-\s*)|(WARN(ING)?\s{2,})|\bWARN(ING)?\b`)
	debugPattern     = regexp.MustCompile(`(?i)(\[?DEBUG\]?[:|\s]+)|(DEBUG\s*-\s*)|(DEBUG\s{2,})|\bDEBUG\b`)
	tracePattern     = regexp.MustCompile(`(?i)(\[?TRACE\]?[:|\s]+)|(TRACE\s*-\s*)|(TRACE\s{2,})|\bTRACE\b`)
	fatalPattern     = regexp.MustCompile(`(?i)(\[?FATAL\]?[:|\s]+)|(FATAL\s*-\s*)|\bFATAL\b`)
	timestampPattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}([.,]\d+)?(Z|[+-]\d{2}:\d{2})`)
)

func (m Model) colorizeOutput(content string) string {
	if !m.colorize {
		return content
	}

	content = infoPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#00AFFF")).Render("$0"))
	content = errorPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("$0"))
	content = warnPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render("$0"))
	content = debugPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Render("$0"))
	content = tracePattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Render("$0"))
	content = fatalPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Render("$0"))
	content = timestampPattern.ReplaceAllString(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render("$0"))

	return content
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for i, server := range m.config.Servers {
		cmds = append(cmds, connectSSHClient(i, &server))
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "?" {
			m.help.ShowAll = !m.help.ShowAll

			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				helpHeight := 2
				if m.help.ShowAll {
					helpHeight = len(m.keys.FullHelp())*2 + 4
				}

				m.tabContents[m.activeTab].ScrollView.SetSize(m.width, m.height-1-helpHeight)

				if !m.tabContents[m.activeTab].ScrollView.UserScrolled() {
					m.tabContents[m.activeTab].ScrollView.GotoBottom()
				}
			}

			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(true)
			}

			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
				*vpModel, _ = vpModel.Update(msg)
			}

		case key.Matches(msg, m.keys.Down), key.Matches(msg, m.keys.PageDown):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.ResetUserScrolledIfAtBottom()
			}

			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
				*vpModel, _ = vpModel.Update(msg)
			}

		case key.Matches(msg, m.keys.PageUp), key.Matches(msg, m.keys.Home):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(true)
			}

			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
				*vpModel, _ = vpModel.Update(msg)
			}

		case key.Matches(msg, m.keys.End):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)

				vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
				*vpModel, _ = vpModel.Update(msg)
			}

		case key.Matches(msg, m.keys.Quit):
			for _, tab := range m.tabContents {
				tab.Close()
			}
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.TabNext):
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
				m.tabContents[m.activeTab].ScrollView.GotoBottom()
			}
			// Adjust tab offset if needed
			// This logic is simple: if activeTab is outside the current view, we might need to scroll.
			// But we don't know the view height here easily without recalculating.
			// We'll handle the offset calculation in View() or a helper, but ideally Update() should handle state.
			// For now, let's just ensure we don't crash. The View() will handle the rendering window.
			return m, nil

		case key.Matches(msg, m.keys.Left), key.Matches(msg, m.keys.TabPrev):
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
				m.tabContents[m.activeTab].ScrollView.GotoBottom()
			}
			return m, nil

		case key.Matches(msg, m.keys.ToggleColor):
			m.colorize = !m.colorize

			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				if m.colorize {
					m.tabContents[m.activeTab].ScrollView.UpdateContent(m.colorizeOutput)
				} else {
					m.tabContents[m.activeTab].ScrollView.UpdateContent(nil)
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.ResetScroll):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
				m.tabContents[m.activeTab].ScrollView.GotoBottom()
			}
			return m, nil

		case key.Matches(msg, m.keys.ClearBuffer):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.Clear()
			}
			return m, nil

		case key.Matches(msg, m.keys.ToggleWordWrap):
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				m.tabContents[m.activeTab].ScrollView.ToggleWordWrap()
			}
			return m, nil

		case key.Matches(msg, m.keys.ToggleTabPosition):
			m.verticalTabs = !m.verticalTabs
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				helpHeight := 2
				if m.help.ShowAll {
					helpHeight = len(m.keys.FullHelp())*2 + 4
				}

				tabWidth := 0
				if m.verticalTabs {
					for _, tab := range m.tabs {
						w := lipgloss.Width(tab) + 4
						if w > tabWidth {
							tabWidth = w
						}
					}
					m.tabContents[m.activeTab].ScrollView.SetSize(m.width-tabWidth, m.height-1-helpHeight)
				} else {
					m.tabContents[m.activeTab].ScrollView.SetSize(m.width, m.height-1-helpHeight)
				}

				if !m.tabContents[m.activeTab].ScrollView.UserScrolled() {
					m.tabContents[m.activeTab].ScrollView.GotoBottom()
				}
			}
			return m, nil
		}

	case tea.MouseMsg:
		if m.verticalTabs {
			if msg.Type == tea.MouseLeft {
				tabWidth := 0
				for _, tab := range m.tabs {
					w := lipgloss.Width(tab) + 6
					if w > tabWidth {
						tabWidth = w
					}
				}

				if msg.X < tabWidth {
					// Calculate the clicked index based on scroll offset
					clickedIndex := m.tabOffset + msg.Y
					if clickedIndex >= 0 && clickedIndex < len(m.tabs) {
						m.activeTab = clickedIndex
						if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
							m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
							m.tabContents[m.activeTab].ScrollView.GotoBottom()
						}
						return m, nil
					}
				}
			}
		} else if msg.Type == tea.MouseLeft && msg.Y == 0 {
			xPos := 0
			for i, tab := range m.tabs {
				tabWidth := lipgloss.Width(tab) + 4

				if msg.X >= xPos && msg.X < xPos+tabWidth {
					m.activeTab = i
					if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
						m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
						m.tabContents[m.activeTab].ScrollView.GotoBottom()
					}
					return m, nil
				}

				xPos += tabWidth
			}
		}

		if msg.Type == tea.MouseMotion && msg.Y != 0 {
			if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
				vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
				if !vpModel.AtBottom() {
					m.tabContents[m.activeTab].ScrollView.SetUserScrolled(true)
				} else {
					m.tabContents[m.activeTab].ScrollView.SetUserScrolled(false)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		for _, tab := range m.tabContents {
			if !tab.HasError {
				tab.ScrollView.SetSize(msg.Width, msg.Height-1)
				tab.ScrollView.UpdateContent(nil)
				if m.colorize {
					tab.ScrollView.UpdateContent(m.colorizeOutput)
				}
			}
		}
		return m, nil

	case updateContentMsg:
		if msg.index < len(m.tabContents) && !m.tabContents[msg.index].HasError {
			if m.colorize {
				m.tabContents[msg.index].ScrollView.UpdateContent(m.colorizeOutput)
			} else {
				m.tabContents[msg.index].ScrollView.UpdateContent(nil)
			}
			if !m.tabContents[msg.index].ScrollView.UserScrolled() {
				m.tabContents[msg.index].ScrollView.GotoBottom()
			}
		}
		return m, nil

	case sshConnectionMsg:
		if msg.index < len(m.tabContents) {
			if msg.err != nil {
				m.tabContents[msg.index].HandleError(msg.err)
				m.tabContents[msg.index].ScrollView.Clear()
				m.tabContents[msg.index].ScrollView.Append("Connection failed: " + msg.err.Error())
			} else {
				m.tabContents[msg.index].SetClient(msg.client)
				m.tabContents[msg.index].ScrollView.Clear()
				m.tabContents[msg.index].ScrollView.Append("Connected to " + lipgloss.NewStyle().Bold(true).Render(m.config.Servers[msg.index].Name) + "\n")
				m.tabContents[msg.index].ScrollView.Append("SSH Version: " + lipgloss.NewStyle().Bold(true).Render(string(msg.client.SSHClient.ServerVersion())) + "\n")

				if len(m.config.Servers[msg.index].Commands) > 0 {
					m.tabContents[msg.index].Client.RunCommands(m.config.Servers[msg.index].Commands)
				}

				index := msg.index
				go func() {
					for {
						select {
						case output := <-msg.client.OutputChan:
							m.tabContents[index].ScrollView.Append(output)

							if program != nil {
								program.Send(updateContentMsg{index: index})
							}

						case err := <-msg.client.ErrChan:
							if err != nil {
								m.tabContents[index].ScrollView.Append("Error: " + err.Error())

								if program != nil {
									program.Send(updateContentMsg{index: index})
								}
							}
						}
					}
				}()
			}
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.help, cmd = m.help.Update(msg)

	if m.activeTab < len(m.tabContents) && !m.tabContents[m.activeTab].HasError {
		vpModel := m.tabContents[m.activeTab].ScrollView.ViewportModel()
		var vpCmd tea.Cmd
		*vpModel, vpCmd = vpModel.Update(msg)
		return m, tea.Batch(cmd, vpCmd)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	m.help.Width = m.width
	helpView := m.help.View(m.keys)

	// Prepare status bar data
	var currentTab *components.TabContent
	var serverName, status, scrollPos string

	if m.activeTab < len(m.tabContents) {
		currentTab = m.tabContents[m.activeTab]
		serverName = currentTab.Name
		if currentTab.HasError {
			status = "Error"
		} else if currentTab.Client == nil {
			status = "Connecting"
		} else {
			status = "Connected"
		}
		scrollPos = fmt.Sprintf("%d/%d", currentTab.ScrollView.LineCount(), components.DefaultMaxLines)
	}

	m.statusBar.Width = m.width
	bar := m.statusBar.View(serverName, status, scrollPos, helpView)
	barHeight := lipgloss.Height(bar)

	var content string
	if currentTab != nil {
		if currentTab.HasError {
			content = components.ErrorStyle.Render(currentTab.ErrorMsg)
		} else {
			if m.verticalTabs {
				tabWidth := 0
				for _, tab := range m.tabs {
					w := lipgloss.Width(tab) + 6
					if w > tabWidth {
						tabWidth = w
					}
				}

				if currentTab.ScrollView.ViewportModel().Width != m.width-tabWidth ||
					currentTab.ScrollView.ViewportModel().Height != m.height-barHeight {
					currentTab.ScrollView.SetSize(m.width-tabWidth, m.height-barHeight)
				}

				currentTab.ScrollView.SetBorder(lipgloss.RoundedBorder())
			} else {
				if currentTab.ScrollView.ViewportModel().Width != m.width ||
					currentTab.ScrollView.ViewportModel().Height != m.height-barHeight-1 { // -1 for tab bar
					currentTab.ScrollView.SetSize(m.width, m.height-barHeight-1)
				}

				currentTab.ScrollView.SetBorder(lipgloss.RoundedBorder())
			}
			content = currentTab.ScrollView.View()
		}
	}

	if m.verticalTabs {
		// Calculate available height for tabs
		tabsHeight := m.height - barHeight
		if tabsHeight < 0 {
			tabsHeight = 0
		}

		// Adjust tabOffset to ensure activeTab is visible
		if m.activeTab < m.tabOffset {
			m.tabOffset = m.activeTab
		} else if m.activeTab >= m.tabOffset+tabsHeight {
			m.tabOffset = m.activeTab - tabsHeight + 1
		}
		// Ensure offset is valid
		if m.tabOffset < 0 {
			m.tabOffset = 0
		}

		var verticalTabBar strings.Builder
		// Render only visible tabs
		endIndex := m.tabOffset + tabsHeight
		if endIndex > len(m.tabs) {
			endIndex = len(m.tabs)
		}

		for i := m.tabOffset; i < endIndex; i++ {
			tab := m.tabs[i]
			icon := " "
			if i == m.activeTab {
				verticalTabBar.WriteString(components.ActiveTabStyle.Render(icon + tab))
			} else {
				verticalTabBar.WriteString(components.TabStyle.Render(icon + tab))
			}
			verticalTabBar.WriteString("\n")
		}

		tabWidth := 0
		for _, tab := range m.tabs {
			w := lipgloss.Width(tab) + 6 // Keep it safe
			if w > tabWidth {
				tabWidth = w
			}
		}

		tabsView := lipgloss.NewStyle().
			Width(tabWidth).
			Height(tabsHeight).
			Render(verticalTabBar.String())

		return lipgloss.JoinHorizontal(lipgloss.Left, tabsView, lipgloss.JoinVertical(lipgloss.Top, content, bar))
	} else {
		var tabBar strings.Builder
		xPos := 0
		for i, tab := range m.tabs {
			icon := " "
			renderedTab := ""
			if i == m.activeTab {
				renderedTab = components.ActiveTabStyle.Render(icon + tab)
			} else {
				renderedTab = components.TabStyle.Render(icon + tab)
			}

			tabWidth := lipgloss.Width(renderedTab)
			// paddedTab := renderedTab // No extra padding needed if style handles it

			tabBar.WriteString(renderedTab)
			xPos += tabWidth
		}

		gapWidth := m.width - xPos
		if gapWidth > 0 {
			tabBar.WriteString(components.TabGap.Width(gapWidth).Render(""))
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			tabBar.String(),
			content,
			bar,
		)
	}
}
