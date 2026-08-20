package cli

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// refreshContent updateChan the UI content, including status line, tree viewport, and log viewport, based on current state.
func (m *Model) refreshContent() {
	lastUpdatedString := " Last Updated: " + m.lastUpdated.Format("15:04:05")
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

	// The tree omits its bottom border because the tooltip supplies the shared
	// divider, so its frame consumes only one vertical row.
	m.treeViewport.SetHeight(int(float64(remainingHeight)*splitRatio) - (borderSpacing - 1))
	m.treeViewport.SetWidth(m.windowWidth - borderSpacing - (2 * horizontalPadding))

	m.logViewport.SetHeight(remainingHeight - m.treeViewport.Height() - (2 * borderSpacing))
	m.logViewport.SetWidth(m.windowWidth - borderSpacing - (2 * horizontalPadding))
}

// formatSpeed formats the speed string to have a uniform size and units.
func formatSpeed(speed string) string {
	if speed == "" {
		return ""
	}

	speed = strings.TrimSpace(speed)

	val, _ := strconv.ParseFloat(speed, 64)

	if val >= 1000 {
		// Convert to Gbps
		return fmt.Sprintf("%8s", fmt.Sprintf("%g Gbps", val/1000))
	}

	return fmt.Sprintf("%8s", fmt.Sprintf("%g Mbps", val))
}
