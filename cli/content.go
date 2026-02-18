package cli

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

// refreshTreeModel rebuilds the tree model based on roots and cursor
func (m *Model) refreshTreeModel() {
	idx := 0
	m.deviceTrees = []*tree.Tree{}
	m.deviceSpeeds = []string{}
	for _, root := range m.roots {
		var deviceTree *tree.Tree
		var rootSpeeds []string
		deviceTree, rootSpeeds, idx = m.buildTreeFromRoot(root, idx)

		m.deviceTrees = append(m.deviceTrees, deviceTree)
		m.deviceSpeeds = append(m.deviceSpeeds, rootSpeeds...)
	}
	m.nodeCount = idx
}

// renderTree renders the tree content to a string
func (m *Model) renderTree() string {
	var deviceTreeSb strings.Builder

	for _, deviceTree := range m.deviceTrees {
		deviceTreeSb.WriteString(deviceTree.String())
		deviceTreeSb.WriteByte('\n')
	}
	nameTreeStr := deviceTreeSb.String()
	speedStr := strings.Join(m.deviceSpeeds, "\n")

	nameTreeWidth := lipgloss.Width(nameTreeStr)
	speedStrWidth := lipgloss.Width(speedStr)
	gapWidth := m.treeViewport.Width - nameTreeWidth - speedStrWidth

	if gapWidth < 1 {
		gapWidth = 1
	}

	gap := strings.Repeat(" ", gapWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, nameTreeStr, gap, speedStr)
}

// buildTreeFromRoot iterates over the tree to build the view and track the cursor
// Returns a name tree and a slice of speed strings, as well as the next index to use
func (m *Model) buildTreeFromRoot(node *lib.TreeNode, currentIdx int) (*tree.Tree, []string, int) {
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
		m.selectedDevice = node
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

// scrollToCursor adjusts the tree viewport's Y offset to keep the cursor centered when possible.
func (m *Model) scrollToCursor() {
	viewportHeight := m.treeViewport.Height
	// Guard against uninitialized or invalid viewport dimensions
	if viewportHeight <= 0 || m.nodeCount <= 0 {
		return
	}

	// Calculate the ideal offset to center the cursor
	idealOffset := m.treeCursor - viewportHeight/2

	// Clamp to valid range: can't scroll above 0
	if idealOffset < 0 {
		idealOffset = 0
	}

	// Can't scroll past the end of content
	maxOffset := m.nodeCount - viewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if idealOffset > maxOffset {
		idealOffset = maxOffset
	}

	m.treeViewport.SetYOffset(idealOffset)
}

// getSelectedDeviceInfo returns formatted device info for the currently selected node
func (m *Model) getSelectedDeviceInfo() string {
	if m.selectedDevice == nil {
		return ""
	}
	node := m.selectedDevice.Device

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

func (m *Model) formatLogContent() string {
	var sb strings.Builder
	for _, entry := range m.log {
		sb.WriteString(m.formatLogEntry(entry))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Model) formatLogEntry(log lib.Log) string {
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
