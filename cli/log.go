package cli

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/AOzmond/usb-tree/lib"
)

func (m *Model) formatLogContent() string {
	var sb strings.Builder
	for _, entry := range m.log {
		sb.WriteString(m.formatLogEntry(entry))
		sb.WriteByte('\n')
	}
	return sb.String()
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
	lhsString := stateStyle.Render(log.Time.Format("15:04:05") + " " + stateString + " " + log.Text + " ")
	rhsString := formatSpeed(log.Speed)
	paddingSize := m.windowWidth - lipgloss.Width(rhsString) - lipgloss.Width(lhsString) - borderSpacing
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
		return
	}
}
