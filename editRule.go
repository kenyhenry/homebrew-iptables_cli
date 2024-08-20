package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type EditRuleObject struct {
	Widget          *widgets.List
	RuleDesc        []string
	BaseTextLengths []int
	ChainName       string
}

func EditRule(chainName string, rule string) *EditRuleObject {
	msgBox := widgets.NewList()
	ruleDesc := []string{
		"table : ",
		"protocol : ",
		"source port : ",
		"dest port : ",
		"source : ",
		"destination : ",
		"module : ",
		"module arg : ",
		"connection states : ",
		"jump : ",
		"jump args : ",
		"in iface : ",
		"out iface : ",
	}

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = false
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)
	msgBox.Rows = ruleDesc

	baseTextLengths := make([]int, len(ruleDesc))
	for i, text := range ruleDesc {
		baseTextLengths[i] = len(text)
	}

	cmd := ExtractAndGenerateCommands(rule, chainName)
	ruleDesc[0] += cmd.Table
	ruleDesc[1] += cmd.Protocol
	ruleDesc[2] += cmd.SPort
	ruleDesc[3] += cmd.DPort
	ruleDesc[4] += cmd.Source
	ruleDesc[5] += cmd.Destination
	ruleDesc[6] += cmd.Module
	ruleDesc[7] += cmd.ModuleArg
	ruleDesc[8] += cmd.ConnectionState
	ruleDesc[9] += cmd.Jump
	ruleDesc[10] += cmd.JumpArg
	ruleDesc[11] += cmd.InIface
	ruleDesc[12] += cmd.OutIface

	return &EditRuleObject{
		Widget:          msgBox,
		RuleDesc:        ruleDesc,
		BaseTextLengths: baseTextLengths,
		ChainName:       chainName,
	}
}

func (nc *EditRuleObject) HandleEvent(e ui.Event, state *UIState) {
	currentRow := nc.Widget.SelectedRow
	baseTextLength := nc.BaseTextLengths[currentRow]
	showOtherWidget := false

	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		// TODO : send command to edit rule
		ret, err := IptablesAddRule(ArraytToCmd(nc.ChainName, nc.RuleDesc, nc.BaseTextLengths))
		if err != nil {
			msgBox := MsgBox(ret)
			state.handlers["msgBox"] = msgBox
			state.SetActive("msgBox")
			state.Render()
		} else {
			ret2, err2 := IptablesAddRule(ArraytToCmd(nc.ChainName, nc.RuleDesc, nc.BaseTextLengths))
			if err2 != nil {
				msgBox := MsgBox(ret2)
				state.handlers["msgBox"] = msgBox
				state.SetActive("msgBox")
				state.Render()
			}
		}
		ui.Clear()
		ui.Render(state.header, state.footer, state.tabpane)
		state.SetActive("chainList")
	case "<Down>":
		nc.Widget.ScrollDown()
	case "<Up>":
		nc.Widget.ScrollUp()
	case "<Backspace>":
		if len(nc.RuleDesc[currentRow]) > baseTextLength {
			nc.RuleDesc[currentRow] = nc.RuleDesc[currentRow][:len(nc.RuleDesc[currentRow])-1]
		}
	case "<Space>":
		nc.RuleDesc[currentRow] += " "
	default:
		if len(e.ID) == 1 {
			nc.RuleDesc[currentRow] += e.ID
		}
	}

	nc.Widget.Rows = nc.RuleDesc
	if !showOtherWidget {
		ui.Render(nc.Widget)
	}
}

func (nr *EditRuleObject) Render() {
	ui.Render(nr.Widget)
}
