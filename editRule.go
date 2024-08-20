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

	// for i := range ruleDesc {
	// 	//  : make ruleSplit match with ruleDesc
	// 	if i < len(ruleSplit) {
	// 		ruleDesc[i] += ruleSplit[i]
	// 	}
	// }

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

func ArraytToCmd(chain string, rules []string, base []int) IptablesCmd {
	return IptablesCmd{
		Chain:           chain,
		Table:           rules[0][base[0]:],
		Protocol:        rules[1][base[1]:],
		DPort:           rules[2][base[2]:],
		SPort:           rules[3][base[3]:],
		Source:          rules[4][base[4]:],
		Destination:     rules[5][base[5]:],
		Module:          rules[6][base[6]:],
		ModuleArg:       rules[7][base[7]:],
		ConnectionState: rules[8][base[8]:],
		Jump:            rules[9][base[9]:],
		JumpArg:         rules[10][base[10]:],
		InIface:         rules[11][base[11]:],
		OutIface:        rules[12][base[12]:],
	}
}
