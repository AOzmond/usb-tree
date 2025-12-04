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

	if node.State == lib.StateAdded {
		nameStyle = nameStyle.Foreground(lipgloss.Color(green))
		speedStyle = speedStyle.Foreground(lipgloss.Color(green))
	} else if node.State == lib.StateRemoved {
		nameStyle = nameStyle.Foreground(lipgloss.Color(red))
		speedStyle = speedStyle.Foreground(lipgloss.Color(red))
	}

	if isSelected {
		nameStyle = nameStyle.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
		speedStyle = speedStyle.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
	}

	nameTree := tree.New().Root(nameStyle.Render(name)).Indenter(compactIndenter).Enumerator(compactEnumerator)
	speeds := []string{speedStyle.Render(speed)}

	for _, child := range node.Children {
		var childNameTree *tree.Tree
		var childSpeeds []string
		childNameTree, childSpeeds, idx = m.buildTreeFromRoot(child, idx)
		nameTree.Child(childNameTree)
		speeds = append(speeds, childSpeeds...)
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
