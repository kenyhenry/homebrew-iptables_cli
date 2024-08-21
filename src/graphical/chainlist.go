package graphical

import (
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kenyhenry/iptables_cli/events"
	"github.com/kenyhenry/iptables_cli/iptables"
	"github.com/kenyhenry/iptables_cli/state"
)

type NewChainlist struct {
	Widget    *widgets.List
	Chainlist []string
	IsMoving  bool
	Em        *events.EventManager
	ChainName string
}

func NewChainList(chainName string, em *events.EventManager) *NewChainlist {
	chainlist, _, _ := iptables.IptablesList(chainName)

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
		Em:        em,
		ChainName: chainName,
	}
}

var ruleStr string
var indexSelected int

func (nc *NewChainlist) HandleEvent(e ui.Event, state *state.UIState) {
	showOtherWidget := false

	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		if !nc.IsMoving {
			indexSelected = nc.Widget.SelectedRow
			ruleStr = nc.Chainlist[indexSelected]
		} else {
			cmd := iptables.ExtractAndGenerateCommands(ruleStr, nc.ChainName)
			cmd.Pos = strconv.Itoa(nc.Widget.SelectedRow + 1)
			if indexSelected <= nc.Widget.SelectedRow {
				cmd.Pos = strconv.Itoa(nc.Widget.SelectedRow + 2)
			}
			ret, err := iptables.IptablesAddRule(cmd)
			if err != nil {
				msgBox := MsgBox(ret)
				state.Handlers["msgBox"] = msgBox
				state.SetActive("msgBox")
				state.Render()
			} else {
				if indexSelected > nc.Widget.SelectedRow {
					indexSelected += 1
				}
				ret2, err2 := iptables.IptablesDeleteRule(nc.ChainName, indexSelected+1)
				if err2 != nil {
					msgBox := MsgBox(ret2)
					state.Handlers["msgBox"] = msgBox
					state.SetActive("msgBox")
					state.Render()
				}
			}
			ui.Clear()
			ui.Render(state.Header, state.Footer, state.Tabpane)
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
		state.Handlers["newRule"] = newRule
		state.SetActive("newRule")
		state.Render()
	case "e":
		if len(nc.Chainlist) == 0 {
			showOtherWidget = true
			editRule := EditRule(nc.ChainName, nc.Chainlist[nc.Widget.SelectedRow], nc.Widget.SelectedRow)
			state.Handlers["editRule"] = editRule
			state.SetActive("editRule")
			state.Render()
		}
	case "P":
		showOtherWidget = true
		info := "Set the policy of chain : " + nc.ChainName
		selectBox := SelectBox(info, "setPolicy", []string{"DROP", "INPUT", "FORWARD", "ACCEPT", "OUTPUT"}, nc.Em)
		state.Handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "D":
		showOtherWidget = true
		info := "Are you sure you want to delete chain : " + nc.ChainName
		selectBox := SelectBox(info, "deleteChain", []string{"yes", "no"}, nc.Em)
		state.Handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "d":
		showOtherWidget = true
		info := "Are you sure you want to delete rule : " + strconv.Itoa(nc.Widget.SelectedRow) + " in chain : " + nc.ChainName
		selectBox := SelectBox(info, "deleteRule", []string{"yes", "no"}, nc.Em)
		state.Handlers["selectBox"] = selectBox
		state.SetActive("selectBox")
		state.Render()
	case "F":
		showOtherWidget = true
		info := "Are you sure you want to Flush : " + nc.ChainName
		selectBox := SelectBox(info, "flushChain", []string{"yes", "no"}, nc.Em)
		state.Handlers["selectBox"] = selectBox
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
