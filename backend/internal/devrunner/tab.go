package devrunner

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
	csize "github.com/nathan-fiscaletti/consolesize-go"
)

func tabBorderWithBottom(left, middle, right string) lg.Border {
	border := lg.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lg.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lg.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lg.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lg.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lg.Center).Border(lg.NormalBorder()).UnsetBorderTop()
)

func RenderTabHeaders(tabs []string, activeTab int) string {
	doc := strings.Builder{}

	var renderedTabs []string

	consoleWidth, _ := csize.GetConsoleSize()

	tabSize := 0

	lineStyle := lg.NewStyle().Foreground(highlightColor).PaddingTop(2)
	renderedTabs = append(renderedTabs, lineStyle.Render("──"))
	tabSize += 2

	for i, t := range tabs {
		var style lg.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "┘"
		} else if isFirst && !isActive {
			border.BottomLeft = "┴"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		rendered := style.Render(t)
		tabSize += lg.Width(rendered)
		renderedTabs = append(renderedTabs, rendered)
	}

	renderedTabs = append(renderedTabs, lineStyle.Render(strings.Repeat("─", consoleWidth-tabSize+1)+"┐"))

	row := lg.JoinHorizontal(lg.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	//doc.WriteString(windowStyle.Width((lg.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	//doc.WriteString(windowStyle.Width(90).Render(m.TabContent[m.activeTab]))
	//return docStyle.Render(doc.String())
	return doc.String()
}
