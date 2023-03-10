package devrunner

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type (
	WatcherDetails struct {
		consoleOutput string
		Status        Status
	}

	UIModel struct {
		Shutdown       func()
		ShuttingDown   bool
		BackendData    WatcherDetails
		FrontendData   WatcherDetails
		lastOutput     string
		activeTabIndex int
	}

	MessageWatcherUpdate struct {
		Output OutputData
	}

	Terminate struct{}
)

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Terminate:
		return m, tea.Quit

	case MessageWatcherUpdate:
		var watcherDetails *WatcherDetails
		switch msg.Output.Source {
		case Backend:
			watcherDetails = &m.BackendData
		case Frontend:
			watcherDetails = &m.FrontendData
		}
		switch msg.Output.Status {
		case Data:
			watcherDetails.Status = Running
			o := msg.Output.Output
			watcherDetails.consoleOutput = o
			m.lastOutput = o
		case ErrorS:
			o := msg.Output.Output
			watcherDetails.Status = ErrorS
			watcherDetails.consoleOutput = o
			m.lastOutput = o
		case Running, Loading, Compiling:
			watcherDetails.Status = msg.Output.Status
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			if m.ShuttingDown {
				return m, tea.Quit
			} else {
				m.ShuttingDown = true
				return m, func() tea.Msg {
					m.Shutdown()
					return nil
				}
			}
		case "right":
			m.activeTabIndex = min(m.activeTabIndex+1, 2)
		case "left":
			m.activeTabIndex = max(m.activeTabIndex-1, 0)
		}
	}

	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	statusStyles     = lg.NewStyle().Padding(1, 2)
	tabContentStyles = lg.NewStyle().Padding(0, 2)
	lastOutputStyles = lg.NewStyle().Padding(2, 0)
)

func (m UIModel) View() string {
	if m.ShuttingDown {
		return lg.NewStyle().Padding(2).Render("shutting down...")
	}
	tabs := RenderTabHeaders([]string{"stats", "backend", "frontend"}, m.activeTabIndex)

	var tab string
	switch m.activeTabIndex {
	case 0:
		tab = m.statsPage()
	case 1:
		tab = m.BackendData.consoleOutput
	case 2:
		tab = m.FrontendData.consoleOutput
	}

	return lg.JoinVertical(0,
		statusStyles.Render("Status message"),
		tabs,
		tabContentStyles.Render(tab),
	)
}

func (m UIModel) statsPage() string {
	return lg.JoinVertical(0,
		"backend:  "+string(m.BackendData.Status),
		"frontend: "+string(m.FrontendData.Status),
		lastOutputStyles.Render(m.lastOutput),
	)
}
