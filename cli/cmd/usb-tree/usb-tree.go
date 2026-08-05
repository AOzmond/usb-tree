package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/AOzmond/usb-tree/cli"
)

func main() {
	teaProgram := tea.NewProgram(cli.InitialModel())
	if _, err := teaProgram.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
