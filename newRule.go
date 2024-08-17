package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// NewChainlist represents the chain list widget with its internal state
type NewRuleObject struct {
	Widget   *widgets.List
	IsMoving bool
}

// New creates and returns a new NewChainlist widget
func NewRule() *NewChainlist {
	msgBox := widgets.NewList()
	ruleDesc := []string{
		"protocol :",
		"direction :",
		"port :",
		"module :",
		"connection states :",
		"jump :",
	}
	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = false
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)
	msgBox.Rows = ruleDesc
	return &NewChainlist{
		Widget:   msgBox,
		IsMoving: false,
	}
}

// HandleEvent handles the keyboard events for the NewChainlist widget
func (nc *NewChainlist) RuleHandleEvent(e ui.Event) {
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
	}
	ui.Render(nc.Widget)
}
