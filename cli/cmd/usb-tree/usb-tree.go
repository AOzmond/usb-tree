package main

import (
	"fmt"
	"os"

	"github.com/AOzmond/usb-tree/cli"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	teaProgram := tea.NewProgram(cli.InitialModel(), tea.WithAltScreen())
	if _, err := teaProgram.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
