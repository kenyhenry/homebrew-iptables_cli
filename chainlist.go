package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type NewChainlist struct {
	Widget    *widgets.List
	Chainlist []string
	IsMoving  bool
}

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

func (nc *NewChainlist) HandleEvent(e ui.Event) {
	switch e.ID {
	case "<Enter>":
		nc.IsMoving = !nc.IsMoving
	case "<Down>":
		nc.Widget.ScrollDown()
		if nc.IsMoving {
			moveDown(nc.Chainlist, nc.Widget.SelectedRow-1)
			nc.Widget.Rows = nc.Chainlist
		}
	case "<Up>":
		nc.Widget.ScrollUp()
		if nc.IsMoving {
			moveUp(nc.Chainlist, nc.Widget.SelectedRow+1)
			nc.Widget.Rows = nc.Chainlist
		}
	}
	ui.Render(nc.Widget)
}

func (nr *NewChainlist) Render() {
	ui.Render(nr.Widget)
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
