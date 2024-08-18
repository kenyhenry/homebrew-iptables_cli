package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type EditRuleObject struct {
	Widget          *widgets.List
	RuleDesc        []string
	IsMoving        bool
	BaseTextLengths []int
}

func EditRule() *EditRuleObject {
	msgBox := widgets.NewList()
	ruleDesc := []string{
		"protocol : ",
		"direction : ",
		"port : ",
		"module : ",
		"connection states : ",
		"jump : ",
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

	return &EditRuleObject{
		Widget:          msgBox,
		RuleDesc:        ruleDesc,
		IsMoving:        false,
		BaseTextLengths: baseTextLengths,
	}
}

func (nc *EditRuleObject) HandleEvent(e ui.Event) {
	currentRow := nc.Widget.SelectedRow
	baseTextLength := nc.BaseTextLengths[currentRow]

	switch e.ID {
	case "<Enter>":
		// TODO : send command to add new rule
		nc.IsMoving = !nc.IsMoving
	case "<Down>":
		nc.Widget.ScrollDown()
	case "<Up>":
		nc.Widget.ScrollUp()
	case "<Backspace>":
		if len(nc.RuleDesc[currentRow]) > baseTextLength {
			nc.RuleDesc[currentRow] = nc.RuleDesc[currentRow][:len(nc.RuleDesc[currentRow])-1]
		}
	default:
		if len(e.ID) == 1 {
			nc.RuleDesc[currentRow] += e.ID
		}
	}

	nc.Widget.Rows = nc.RuleDesc
	ui.Render(nc.Widget)
}

func (nr *EditRuleObject) Render() {
	ui.Render(nr.Widget)
}