package cli

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/AOzmond/usb-tree/lib"
)

func (m *Model) formatLogContent() string {
	var sb strings.Builder
	for i, entry := range m.log {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(m.formatLogEntry(entry))
	}
	return sb.String()
}

func (m *Model) scrollLogUp() {
	m.logViewport.ScrollUp(1)
	m.clampLogViewport()
}

func (m *Model) scrollLogDown() {
	m.logViewport.ScrollDown(1)
	m.updateLogScrollState()
	m.clampLogViewport()
}

func (m *Model) updateLogScrollState() {
	if m.logViewport.AtBottom() {
		m.logHasNew = false
	}
}

func (m *Model) formatLogEntry(log lib.Log) string {
	stateString := " "
	if log.State == lib.StateRemoved {
		stateStyle = removedLogStyle
		stateString = "-"
	} else if log.State == lib.StateAdded {
		stateStyle = addedLogStyle
		stateString = "+"
	}
	rhsString := formatSpeed(log.Speed)
	logPrefix := log.Time.Format("15:04:05") + " " + stateString + " "
	availableForName := m.logViewport.Width() - lipgloss.Width(logPrefix) - lipgloss.Width(rhsString)
	if availableForName < 0 {
		availableForName = 0
	}
	name := middleTruncate(log.Text, availableForName)
	lhsString := stateStyle.Render(logPrefix + name)

	paddingSize := m.logViewport.Width() - lipgloss.Width(rhsString) - lipgloss.Width(lhsString)
	if paddingSize < 0 {
		paddingSize = 0
	}
	padding := strings.Repeat(" ", paddingSize)
	return lipgloss.JoinHorizontal(lipgloss.Left, lhsString, padding, rhsString)
}

// clampLogViewport prevents the log viewport from scrolling past the available content.
func (m *Model) clampLogViewport() {
	if m.logViewport.Height() <= 0 {
		return
	}

	contentHeight := lipgloss.Height(m.logContent)
	maxYOffset := contentHeight - m.logViewport.Height()
	if maxYOffset < 0 {
		maxYOffset = 0
	}
	if m.logViewport.YOffset() > maxYOffset {
		m.logViewport.SetYOffset(maxYOffset)
	} else if m.logViewport.YOffset() < 0 {
		m.logViewport.SetYOffset(0)
	}
}
