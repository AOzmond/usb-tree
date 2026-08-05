package cli

import "charm.land/lipgloss/v2"

// refreshContent updateChan the UI content, including status line, tree viewport, and log viewport, based on current state.
func (m *Model) refreshContent() {
	lastUpdatedString := "Last Updated: " + m.lastUpdated.Format("15:04:05")
	lastUpdatedWidth := lipgloss.Width(lastUpdatedString) + 1

	helpView := m.helpModel.View(keys)

	helpViewStyle := windowStyle.
		Width(m.windowWidth - lastUpdatedWidth).
		Align(lipgloss.Center)

	renderedHelp := helpViewStyle.Render(helpView)

	m.statusLine = statusStyle.
		Width(m.windowWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, lastUpdatedString, " ", renderedHelp))

	m.recalculateDimensions(m.statusLine)
	m.treeViewport.SetContent(m.renderTree())
	m.logViewport.SetContent(m.logContent)
}

// recalculateDimensions adjusts the dimensions of the tree and log viewports based on window size and status line height.
func (m *Model) recalculateDimensions(statusLine string) {
	m.statusHeight = lipgloss.Height(statusLine)
	remainingHeight := m.windowHeight - m.statusHeight - tooltipHeight

	m.treeViewport.SetHeight(int(float64(remainingHeight)*splitRatio) - borderSpacing)
	m.treeViewport.SetWidth(m.windowWidth - borderSpacing)

	m.logViewport.SetHeight(remainingHeight - m.treeViewport.Height() - (2 * borderSpacing))
	m.logViewport.SetWidth(m.windowWidth - borderSpacing)
}
