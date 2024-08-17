package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewChainlist represents the chain list widget with its internal state
type NewChainlist struct {
	Widget    *widgets.List
	Chainlist []string
	IsMoving  bool
}

// New creates and returns a new NewChainlist widget
func NewChainList(chainName string) *NewChainlist {
	// TODO get iptables rule and screen
	chainlist := []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}

	// Initialize the List widget
	l := widgets.NewList()
	l.Rows = chainlist
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(1, 12, termWidth-1, termHeight-4)

	return &NewChainlist{
		Widget:    l,
		Chainlist: chainlist,
		IsMoving:  false,
	}
}

// HandleEvent handles the keyboard events for the NewChainlist widget
func (nc *NewChainlist) HandleEvent(e ui.Event) {
	switch e.ID {
	case "<Enter>":
		nc.IsMoving = !nc.IsMoving
	case "j", "<Down>":
		nc.Widget.ScrollDown()
		if nc.IsMoving {
			moveDown(nc.Chainlist, nc.Widget.SelectedRow-1)
			nc.Widget.Rows = nc.Chainlist
		}
	case "k", "<Up>":
		nc.Widget.ScrollUp()
		if nc.IsMoving {
			moveUp(nc.Chainlist, nc.Widget.SelectedRow+1)
			nc.Widget.Rows = nc.Chainlist
		}
	case "<C-d>":
		nc.Widget.ScrollHalfPageDown()
	case "<C-u>":
		nc.Widget.ScrollHalfPageUp()
	case "<C-f>":
		nc.Widget.ScrollPageDown()
	case "<C-b>":
		nc.Widget.ScrollPageUp()
	}
	ui.Render(nc.Widget)
}

func moveUp(slice []string, index int) {
	if index > 0 && index < len(slice) {
		slice[index], slice[index-1] = slice[index-1], slice[index]
	}
}

func moveDown(slice []string, index int) {
	if index >= 0 && index < len(slice)-1 {
		slice[index], slice[index+1] = slice[index+1], slice[index]
	}
}
