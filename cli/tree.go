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
		deviceTree, rootSpeeds, idx = m.buildTreeFromNode(root, idx)

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

// buildTreeFromNode iterates over the tree to build the view and track the cursor
// Returns a name tree and a slice of speed strings, as well as the next index to use
func (m *Model) buildTreeFromNode(node *lib.TreeNode, currentIdx int) (*tree.Tree, []string, int) {
	isSelected := currentIdx == m.treeCursor
	idx := currentIdx + 1

	name := strings.TrimSpace(node.Name)
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
		m.selectedDevice = node.Device
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

			childDeviceTree, childSpeeds, idx = m.buildTreeFromNode(child, idx)
			nameTree.Child(childDeviceTree)
			speeds = append(speeds, childSpeeds...)
		}
	}

	return nameTree, speeds, idx
}

// formatSpeed formats the speed string to have a uniform size and units.
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

// scrollToCursor adjusts the tree viewport's Y offset to keep the cursor centered when possible and stops scrolling past the content.
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
