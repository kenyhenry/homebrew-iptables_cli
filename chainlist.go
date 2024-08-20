package main

import (
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type NewChainlist struct {
	Widget    *widgets.List
	Chainlist []string
	IsMoving  bool
	Em        *EventManager
	ChainName string
}

func NewChainList(chainName string, em *EventManager) *NewChainlist {
	chainlist, _, _ := IptablesList(chainName)

	l := widgets.NewList()
	l.Rows = chainlist
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(1, 12, termWidth-1, termHeight-4)

	// BUG : add listener each new chainlist, induce multiple deleteCRule
	em.AddListener("deleteRule", func(e Event) {
		if e.Data == "yes" {
			IptablesDeleteRule(chainName, l.SelectedRow)
		}

	})

	return &NewChainlist{
		Widget:    l,
		Chainlist: chainlist,
		IsMoving:  false,
		Em:        em,
		ChainName: chainName,
	}
}

var ruleStr string
var indexSelected int

func (nc *NewChainlist) HandleEvent(e ui.Event, state *UIState) {
	showOtherWidget := false

	switch e.ID {
	// BUG : move iptable rule not working index error
	case "<Enter>":
		if !nc.IsMoving {
			indexSelected = nc.Widget.SelectedRow
			ruleStr = nc.Chainlist[indexSelected]
		} else {
			cmd := ExtractAndGenerateCommands(ruleStr, nc.ChainName)
			if indexSelected <= nc.Widget.SelectedRow-1 {
				cmd.Pos = strconv.Itoa(nc.Widget.SelectedRow + 1)
			} else {
				cmd.Pos = strconv.Itoa(nc.Widget.SelectedRow - 1)
			}
			ret, err := IptablesAddRule(cmd)
			if err != nil {
				msgBox := MsgBox(ret)
				state.handlers["msgBox"] = msgBox
				state.SetActive("msgBox")
				state.Render()
			} else {
				// if indexSelected < nc.Widget.SelectedRow {
				// 	indexSelected++
				// }
				// ret2, err2 := IptablesDeleteRule(nc.ChainName, indexSelected)
				// if err2 != nil {
				// msgBox := MsgBox(ret2)
				// state.handlers["msgBox"] = msgBox
				// state.SetActive("msgBox")
				// state.Render()
				// }
			}
			ui.Clear()
			ui.Render(state.header, state.footer, state.tabpane)
			state.SetActive("chainList")
			state.Render()

		}
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
	case "a":
		showOtherWidget = true
		newRule := NewRule(nc.ChainName)
		state.handlers["newRule"] = newRule
		state.SetActive("newRule")
		state.Render()
	case "e":
		showOtherWidget = true
		editRule := EditRule(nc.ChainName, nc.Chainlist[nc.Widget.SelectedRow], nc.Widget.SelectedRow)
		state.handlers["editRule"] = editRule
		state.SetActive("editRule")
		state.Render()
	case "P":
		showOtherWidget = true
		info := "Set the policy of chain : " + nc.ChainName
		selectBox := SelectBox(info, "setPolicy", []string{"DROP", "INPUT", "FORWARD", "ACCEPT", "OUTPUT"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "D":
		showOtherWidget = true
		info := "Are you sure you want to delete chain : " + nc.ChainName
		selectBox := SelectBox(info, "deleteChain", []string{"yes", "no"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "d":
		showOtherWidget = true
		info := "Are you sure you want to delete rule : " + strconv.Itoa(nc.Widget.SelectedRow) + " in chain : " + nc.ChainName
		selectBox := SelectBox(info, "deleteRule", []string{"yes", "no"}, nc.Em)
		state.handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "F":
		showOtherWidget = true
		info := "Are you sure you want to Flush : " + nc.ChainName
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
