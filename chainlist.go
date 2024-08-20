package main

import (
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type NewChainlist struct {
	Widget     *widgets.List
	Chainlist  []string
	IsMoving   bool
	Em         *EventManager
	ChaineName string
}

func NewChainList(chainName string, em *EventManager) *NewChainlist {
	chainlist, _, _ := IptablesList(chainName)

	l := widgets.NewList()
	l.Rows = chainlist
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(1, 12, termWidth-1, termHeight-4)

	em.AddListener("deleteRule", func(e Event) {
		if e.Data == "yes" {
			IptablesDeleteRule(chainName, l.SelectedRow)
		}

	})

	return &NewChainlist{
		Widget:     l,
		Chainlist:  chainlist,
		IsMoving:   false,
		Em:         em,
		ChaineName: chainName,
	}
}

func (nc *NewChainlist) HandleEvent(e ui.Event, state *UIState) {
	showOtherWidget := false
	switch e.ID {
	case "<Enter>":
		nc.IsMoving = !nc.IsMoving
		// TODO : on enter and rule move insert rule at selected and remove rule at old place
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
	case "a":
		showOtherWidget = true
		newRule := NewRule(nc.ChaineName)
		state.handlers["newRule"] = newRule
		state.SetActive("newRule")
		state.Render()
	case "e":
		showOtherWidget = true
		editRule := EditRule(nc.ChaineName, nc.Chainlist[nc.Widget.SelectedRow])
		state.handlers["editRule"] = editRule
		state.SetActive("editRule")
		state.Render()
	case "P":
		showOtherWidget = true
		info := "Set the policy of chain : " + nc.ChaineName
		selectBox := SelectBox(info, "setPolicy", []string{"DROP", "INPUT", "FORWARD", "ACCEPT", "OUTPUT"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "D":
		showOtherWidget = true
		info := "Are you sure you want to delete chain : " + nc.ChaineName
		selectBox := SelectBox(info, "deleteChain", []string{"yes", "no"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "d":
		showOtherWidget = true
		info := "Are you sure you want to delete rule : " + strconv.Itoa(nc.Widget.SelectedRow) + " in chain : " + nc.ChaineName
		selectBox := SelectBox(info, "deleteRule", []string{"yes", "no"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "F":
		showOtherWidget = true
		info := "Are you sure you want to Flush : " + nc.ChaineName
		selectBox := SelectBox(info, "flushChain", []string{"yes", "no"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	}
	if !showOtherWidget {
		ui.Render(nc.Widget)
		showOtherWidget = false
	}
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
