package main

import (
	"strings"

	"github.com/AOzmond/usb-tree/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

func compactIndenter(children tree.Children, index int) string {
	if children.Length()-1 == index {
		return "  "
	}
	return "│ "
}

func compactEnumerator(children tree.Children, index int) string {
	if children.Length()-1 == index {
		return "└─"
	}
	return "├─"
}

func waitForUpdate(sub chan []lib.Device) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// refreshTreeContent rebuilds the visual tree based on roots and cursor
func (m *model) refreshTreeContent() {
	var sb strings.Builder
	idx := 0
	for _, root := range m.roots {
		var newTree *tree.Tree
		newTree, idx = m.buildTreeFromRoot(root, idx)
		sb.WriteString(newTree.String())
		sb.WriteByte('\n')
	}
	m.nodeCount = idx
	m.treeContent = sb.String()
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
func (m *model) buildTreeFromRoot(node *lib.TreeNode, currentIdx int) (*tree.Tree, int) {
	isSelected := currentIdx == m.treeCursor
	idx := currentIdx + 1

	name := node.Name
	style := lipgloss.NewStyle()

	if node.State == lib.StateAdded {
		style = style.Foreground(lipgloss.Color(green))
	} else if node.State == lib.StateRemoved {
		style = style.Foreground(lipgloss.Color(red))
	}

	if isSelected {
		style = style.Background(lipgloss.Color(white)).Foreground(lipgloss.Color("0"))
	}

	newTree := tree.New().Root(style.Render(name)).Indenter(compactIndenter).Enumerator(compactEnumerator)
	for _, child := range node.Children {
		var childTree *tree.Tree
		childTree, idx = m.buildTreeFromRoot(child, idx)
		newTree.Child(childTree)
	}
	return newTree, idx
}
