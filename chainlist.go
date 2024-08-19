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

func NewChainList(chainName string) *NewChainlist {
	// TODO send iptables command to get iptables rule and screen
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

	em := NewEventManager()

	l := widgets.NewList()
	l.Rows = chainlist
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(1, 12, termWidth-1, termHeight-4)

	em.AddListener("deleteChain", func(e Event) {
		// TODO : delete chain selected "chainName"
		// TODO : ret of command
		if e.Data == "yes" {
			ret := string(e.Data)
			msgBox := MsgBox(ret)
			ui.Render(msgBox.Widget)
		}

	})

	em.AddListener("deleteRule", func(e Event) {
		// TODO : delete rule selected chainlist[l.SelectedRow]
		// TODO : ret of command
		if e.Data == "yes" {
			ret := string(e.Data)
			msgBox := MsgBox(ret)
			ui.Render(msgBox.Widget)
		}

	})

	em.AddListener("setPolicy", func(e Event) {
		// TODO : set chain policy as selectBox to e.Data
		// TODO : ret of command
		ret := string(e.Data)
		msgBox := MsgBox(ret)
		ui.Render(msgBox.Widget)

	})

	em.AddListener("flushChain", func(e Event) {
		info := "Delete all rules from chain : " + chainName
		selectBox := SelectBox(info, "flushConfirm", []string{"yes", "no"}, em)
		ui.Render(selectBox.Widget)
		ret := string(e.Data)
		msgBox := MsgBox(ret)
		ui.Render(msgBox.Widget)

	})

	em.AddListener("flushConfirm", func(e Event) {
		// TODO : send iptables flush chain except OUTPUT, INPUT, FORWARD
		// TODO : ret of command
		ret := string(e.Data)
		msgBox := MsgBox(ret)
		ui.Render(msgBox.Widget)

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
		editRule := EditRule(nc.Chainlist[nc.Widget.SelectedRow])
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
		info := "Are you sure you want to delete rule : " + strconv.Itoa(nc.Widget.SelectedRow)
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
