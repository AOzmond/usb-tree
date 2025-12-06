package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AOzmond/usb-tree/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

// compactIndenter overwrite the default tree indenter to use a single space instead of two
func compactIndenter(children tree.Children, index int) string {
	if children.Length()-1 == index {
		return "  "
	}
	return "│ "
}

// compactEnumerator overwrite the default tree enumerator to reduce spacing
func compactEnumerator(children tree.Children, index int) string {
	if children.Length()-1 == index {
		return "└─"
	}
	return "├─"
}

// waitForUpdate consumes the next update message from the provided subscription channel and returns it as a command.
func waitForUpdate(sub chan []lib.Device) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// refreshTreeContent rebuilds the visual tree based on roots and cursor
func (m *model) refreshTreeContent() {
	var deviceTreeSb strings.Builder
	var speeds []string
	idx := 0
	for _, root := range m.roots {
		var deviceTree *tree.Tree
		var rootSpeeds []string
		deviceTree, rootSpeeds, idx = m.buildTreeFromRoot(root, idx)

		deviceTreeSb.WriteString(deviceTree.String())
		deviceTreeSb.WriteByte('\n')
		speeds = append(speeds, rootSpeeds...)
	}
	m.nodeCount = idx

	nameTreeStr := deviceTreeSb.String()
	speedStr := strings.Join(speeds, "\n")

	nameTreeWidth := lipgloss.Width(nameTreeStr)
	speedStrWidth := lipgloss.Width(speedStr)
	gapWidth := m.treeViewport.Width - nameTreeWidth - speedStrWidth

	if gapWidth < 1 {
		gapWidth = 1
	}

	gap := strings.Repeat(" ", gapWidth)

	m.treeContent = lipgloss.JoinHorizontal(lipgloss.Top, nameTreeStr, gap, speedStr)
}

// updateViewportForCursor ensures the cursor is visible in the viewport
func (m *model) updateViewportForCursor() {
	headerHeight := 1
	visualCursor := m.treeCursor + headerHeight

	if visualCursor < m.treeViewport.YOffset {
		m.treeViewport.SetYOffset(visualCursor)
	} else if visualCursor >= m.treeViewport.YOffset+m.treeViewport.Height {
		m.treeViewport.SetYOffset(visualCursor - m.treeViewport.Height + 1)
	}
}

// buildTreeFromRoot iterates over the tree to build the view and track the cursor
// Returns a name tree and a slice of speed strings, as well as the next index to use
func (m *model) buildTreeFromRoot(node *lib.TreeNode, currentIdx int) (*tree.Tree, []string, int) {
	isSelected := currentIdx == m.treeCursor
	idx := currentIdx + 1

	name := node.Name
	speed := formatSpeed(node.Speed)
	nameStyle := lipgloss.NewStyle()
	speedStyle := lipgloss.NewStyle()

	var statusPrefix string
	switch node.State {
	case lib.StateAdded:
		statusPrefix = "+ "
		nameStyle = nameStyle.Foreground(lipgloss.Color(green))
		speedStyle = speedStyle.Foreground(lipgloss.Color(green))

	case lib.StateRemoved:
		statusPrefix = "- "
		nameStyle = nameStyle.Foreground(lipgloss.Color(red))
		speedStyle = speedStyle.Foreground(lipgloss.Color(red))

	default:
		statusPrefix = ""

	}

	// Add expand/collapse indicator for nodes with children
	var childrenIndicator string
	if len(node.Children) > 0 {
		if m.collapsed[node] {
			childrenIndicator = "▶ "
		} else {
			childrenIndicator = "▼ "
		}
	}

	if isSelected {
		nameStyle = nameStyle.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
		speedStyle = speedStyle.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
	}

	nameTree := tree.New().
		Root(nameStyle.Render(childrenIndicator + statusPrefix + name)).
		Indenter(compactIndenter).
		Enumerator(compactEnumerator)

	speeds := []string{speedStyle.Render(speed)}

	// Only render children if not collapsed
	if !m.collapsed[node] {
		for _, child := range node.Children {
			var childDeviceTree *tree.Tree
			var childSpeeds []string
			childDeviceTree, childSpeeds, idx = m.buildTreeFromRoot(child, idx)
			nameTree.Child(childDeviceTree)
			speeds = append(speeds, childSpeeds...)
		}
	}
	return nameTree, speeds, idx
}

// formatSpeed formats the speed string to have a uniform size, and units.
func formatSpeed(speed string) string {
	if speed == "" {
		return ""
	}

	speed = strings.TrimSpace(speed)

	val, err := strconv.ParseFloat(speed, 64)
	if err != nil {
		return fmt.Sprintf("%8s", speed)
	}

	if val >= 1000 {
		// Convert to Gbps
		return fmt.Sprintf("%8s", fmt.Sprintf("%g Gbps", val/1000))
	}

	return fmt.Sprintf("%8s", fmt.Sprintf("%g Mbps", val))
}

// getNodeAtCursor returns the TreeNode at the current cursor position
func (m *model) getNodeAtCursor() *lib.TreeNode {
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
func (m *model) findNodeAtIndex(node *lib.TreeNode, currentIdx int, targetIdx int) (*lib.TreeNode, int) {
	if currentIdx == targetIdx {
		return node, currentIdx + 1
	}
	idx := currentIdx + 1

	if m.collapsed[node] {
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

// getSelectedDeviceInfo returns formatted device info for the currently selected node
func (m *model) getSelectedDeviceInfo() string {
	node := m.getNodeAtCursor()
	if node == nil {
		return "No device selected"
	}

	busStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(plum))
	deviceStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gold))
	vidStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(coralRed))
	pidStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(paleGreen))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gray))
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(skyBlue))

	busString := busStyle.Render("Bus: ", strconv.Itoa(node.Bus))
	deviceString := deviceStyle.Render(" Device: ", strconv.Itoa(node.DevNum))
	vidString := vidStyle.Render(" VID: ", node.VendorID)
	pidString := pidStyle.Render(" PID: ", node.ProductID)

	deviceInfo := lipgloss.JoinHorizontal(lipgloss.Left, busString, deviceString, vidString, pidString)

	nameString := nameStyle.Render(node.Name)
	linkString := linkStyle.Render(getDbAddress(node.VendorID, node.ProductID))

	tooltipString := lipgloss.JoinVertical(lipgloss.Top, deviceInfo, nameString, linkString)

	return tooltipString
}

// getDbAddress returns the USB-ID database link for the given VID and PID
func getDbAddress(vid string, pid string) string {
	baseAddress := "https://the-sz.com/products/usbid/?v="
	return baseAddress + vid + "&p=" + pid
}

func (m *model) formatLogContent() string {
	var sb strings.Builder
	for _, entry := range m.log {
		sb.WriteString(m.formatLogEntry(entry))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *model) formatLogEntry(log lib.Log) string {
	addedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(green))
	removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(red))
	stateStyle := lipgloss.NewStyle()
	stateString := " "
	if log.State == lib.StateRemoved {
		stateStyle = removedStyle
		stateString = "-"
	} else if log.State == lib.StateAdded {
		stateStyle = addedStyle
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
