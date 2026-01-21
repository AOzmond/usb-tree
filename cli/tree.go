package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AOzmond/usb-tree/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// waitForUpdate listens to a channel for a slice of lib.Device updateChan and returns a tea.Cmd to process the message.
func waitForUpdate(sub chan []lib.Device) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// updateNodeCount updateChan the nodeCount based on visible devices.
func (m *Model) updateNodeCount() {
	idx := 0
	for _, root := range m.roots {
		idx = m.countVisibleNodes(root, idx)
	}
	m.nodeCount = idx
}

// countVisibleNodes recursively counts the number of visible nodes under the given node
func (m *Model) countVisibleNodes(node *lib.TreeNode, currentIdx int) int {
	idx := currentIdx + 1
	if !m.collapsed[node.Key()] {
		for _, child := range node.Children {
			idx = m.countVisibleNodes(child, idx)
		}
	}
	return idx
}

// renderTree renders the tree content to a string
func (m *Model) renderTree() string {
	var deviceTreeSb strings.Builder
	idx := 0

	totalWidth := m.treeViewport.Width

	for _, root := range m.roots {
		lines, nextIdx := m.renderNode(root, idx, []bool{}, totalWidth)
		deviceTreeSb.WriteString(strings.Join(lines, "\n"))
		if nextIdx < m.nodeCount {
			deviceTreeSb.WriteByte('\n')
		}
		idx = nextIdx
	}

	return deviceTreeSb.String()
}

// renderNode recursively renders a node and its children
func (m *Model) renderNode(node *lib.TreeNode, currentIdx int, continues []bool, totalWidth int) ([]string, int) {
	isSelected := currentIdx == m.treeCursor
	idx := currentIdx + 1

	rowStyle, contentStyle := m.getNodeStyles(node, isSelected)
	indicators, contentStyle := m.getNodeIndicators(node, contentStyle)
	prefixStr := m.buildTreePrefix(continues)

	line := m.renderNodeLine(node, prefixStr, indicators, rowStyle, contentStyle, totalWidth)
	renderedDevices := []string{line}

	// Only render children if not collapsed
	if !m.collapsed[node.Key()] {
		for i, child := range node.Children {
			isLast := i == len(node.Children)-1
			childContinues := append(continues, !isLast)
			var childLines []string
			childLines, idx = m.renderNode(child, idx, childContinues, totalWidth)
			renderedDevices = append(renderedDevices, childLines...)
		}
	}

	return renderedDevices, idx
}

// getNodeStyles determines and returns the row and content styles for a tree node based on its state and selection status.
func (m *Model) getNodeStyles(node *lib.TreeNode, isSelected bool) (lipgloss.Style, lipgloss.Style) {
	rowStyle := windowStyle
	if isSelected {
		rowStyle = rowStyle.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
	}

	contentStyle := rowStyle
	switch node.State {
	case lib.StateAdded:
		contentStyle = contentStyle.Foreground(lipgloss.Color(green))
	case lib.StateRemoved:
		contentStyle = contentStyle.Foreground(lipgloss.Color(red))
	}

	return rowStyle, contentStyle
}

// getNodeIndicators returns visual indicators for a node, such as collapse/expand symbols and state prefixes.
func (m *Model) getNodeIndicators(node *lib.TreeNode, contentStyle lipgloss.Style) (string, lipgloss.Style) {
	var childrenIndicator string
	if len(node.Children) > 0 {
		if m.collapsed[node.Key()] {
			childrenIndicator = "▶ "
			if m.hasChangedChild(node) {
				contentStyle = contentStyle.Foreground(lipgloss.Color(orange))
			}
		} else {
			childrenIndicator = "▼ "
		}
	}

	var statusPrefix string
	switch node.State {
	case lib.StateAdded:
		statusPrefix = "+ "
	case lib.StateRemoved:
		statusPrefix = "- "
	}

	return childrenIndicator + statusPrefix, contentStyle
}

// buildTreePrefix generates a string representation of a tree structure prefix based on the provided continues slice.
func (m *Model) buildTreePrefix(continues []bool) string {
	var prefix strings.Builder
	for i := 0; i < len(continues); i++ {
		if i == len(continues)-1 {
			if continues[i] {
				prefix.WriteString("├─")
			} else {
				prefix.WriteString("└─")
			}
		} else {
			if continues[i] {
				prefix.WriteString("│ ")
			} else {
				prefix.WriteString("  ")
			}
		}
	}
	return prefix.String()
}

// renderNodeLine generates a formatted string representing a tree node line with styles, truncation, and aligned elements.
func (m *Model) renderNodeLine(node *lib.TreeNode, prefixStr, indicators string, rowStyle, contentStyle lipgloss.Style, totalWidth int) string {
	name := strings.TrimSpace(node.Name)
	speed := formatSpeed(node.Speed)

	speedWidth := lipgloss.Width(speed)
	prefixWidth := lipgloss.Width(prefixStr)
	indicatorsWidth := lipgloss.Width(indicators)
	gapWidth := 1

	availableForName := totalWidth - prefixWidth - indicatorsWidth - gapWidth - speedWidth
	if availableForName < 5 {
		availableForName = 5
	}

	truncatedName := middleTruncate(name, availableForName)

	leftPart := prefixStr + indicators + truncatedName
	rightPart := speed

	actualGapWidth := totalWidth - lipgloss.Width(leftPart) - lipgloss.Width(rightPart)
	if actualGapWidth < 1 {
		actualGapWidth = 1
	}
	gap := strings.Repeat(" ", actualGapWidth)

	return rowStyle.Render(prefixStr) + contentStyle.Render(indicators+truncatedName) + rowStyle.Render(gap+rightPart)
}

// middleTruncate shortens a string by replacing its middle with "…" if its length exceeds the specified maxLen.
func middleTruncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}

	half := (maxLen - 1) / 2
	return s[:half] + "…" + s[len(s)-(maxLen-half-1):]
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

// scrollToCursor adjusts the viewport's offset to ensure the cursor is visible when resizing the screen
func (m *Model) scrollToCursor() {
	m.treeViewport.SetYOffset(m.treeCursor)
}

// scrollUpToCursor ensures that the tree cursor remains within the visible portion of the viewport when the cursor moves up
func (m *Model) scrollUpToCursor() {
	viewportHeight := m.treeViewport.Height
	// Guard against uninitialized or invalid viewport dimensions
	if viewportHeight <= 0 || m.nodeCount <= 0 {
		return
	}

	padding := 0
	if (viewportHeight - 2) > 4 {
		padding = 2
	}

	if m.treeCursor < (m.treeViewport.YOffset + padding) {
		m.treeViewport.SetYOffset(m.treeCursor - padding)
	}
}

// scrollUpToCursor ensures that the tree cursor remains within the visible portion of the viewport when the cursor moves down
func (m *Model) scrollDownToCursor() {
	viewportHeight := m.treeViewport.Height
	// Guard against uninitialized or invalid viewport dimensions
	if viewportHeight <= 0 || m.nodeCount <= 0 {
		return
	}

	padding := 1
	if (viewportHeight - 2) > 4 {
		padding = 3
	}

	if m.treeCursor > (m.treeViewport.YOffset + viewportHeight - padding) {
		m.treeViewport.SetYOffset(m.treeCursor - viewportHeight + padding)
	}
}

// getNodeAtCursor returns the TreeNode at the current cursor position
func (m *Model) getNodeAtCursor() *lib.TreeNode {
	idx := 0
	for _, root := range m.roots {
		node, newIdx := m.findNodeAtIndex(root, idx, m.treeCursor)
		if node != nil {
			return node
		}
		idx = newIdx
	}
	return nil
}

// findNodeAtIndex recursively searches for the node at the target index
func (m *Model) findNodeAtIndex(node *lib.TreeNode, currentIdx int, targetIdx int) (*lib.TreeNode, int) {
	if currentIdx == targetIdx {
		return node, currentIdx + 1
	}

	idx := currentIdx + 1

	if m.collapsed[node.Key()] {
		return nil, idx
	}

	for _, child := range node.Children {
		found, newIdx := m.findNodeAtIndex(child, idx, targetIdx)
		if found != nil {
			return found, newIdx
		}
		idx = newIdx
	}

	return nil, idx
}

// hasChangedChild recursively checks if any child of the node has a status update
func (m *Model) hasChangedChild(node *lib.TreeNode) bool {
	for _, child := range node.Children {
		if child.State != lib.StateNormal {
			return true
		}
		if m.hasChangedChild(child) {
			return true
		}
	}
	return false
}

// checkOffscreenChanges returns whether there are changes (added/removed) above and below the visible viewport
func (m *Model) checkOffscreenChanges() (above bool, below bool) {
	idx := 0
	for _, root := range m.roots {
		above, below, idx = m.checkNodeOffscreenChanges(root, idx, above, below)
	}
	return above, below
}

func (m *Model) checkNodeOffscreenChanges(node *lib.TreeNode, currentIdx int, above bool, below bool) (bool, bool, int) {
	if node.State != lib.StateNormal {
		if currentIdx < m.treeViewport.YOffset {
			above = true
		} else if currentIdx >= m.treeViewport.YOffset+m.treeViewport.Height {
			below = true
		}
	}

	// check children within collapsed nodes
	if m.collapsed[node.Key()] {
		if m.hasChangedChild(node) {
			if currentIdx < m.treeViewport.YOffset {
				above = true
			} else if currentIdx >= m.treeViewport.YOffset+m.treeViewport.Height {
				below = true
			}
		}
		return above, below, currentIdx + 1
	}

	// check children within non-collapsed nodes
	idx := currentIdx + 1
	for _, child := range node.Children {
		above, below, idx = m.checkNodeOffscreenChanges(child, idx, above, below)
	}

	return above, below, idx
}
