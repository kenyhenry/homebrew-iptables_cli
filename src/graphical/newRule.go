package graphical

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kenyhenry/iptables_cli/iptables"
	"github.com/kenyhenry/iptables_cli/state"
)

type NewRuleObject struct {
	Widget          *widgets.List
	RuleDesc        []string
	BaseTextLengths []int
	ChainName       string
}

func NewRule(chaineName string) *NewRuleObject {
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
		"jump arg : ",
		"in iface : ",
		"out iface : ",
	}

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.Title = chaineName
	msgBox.WrapText = false
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)
	msgBox.Rows = ruleDesc

	baseTextLengths := make([]int, len(ruleDesc))
	for i, text := range ruleDesc {
		baseTextLengths[i] = len(text)
	}

	return &NewRuleObject{
		Widget:          msgBox,
		RuleDesc:        ruleDesc,
		BaseTextLengths: baseTextLengths,
		ChainName:       chaineName,
	}
}

func (nc *NewRuleObject) HandleEvent(e ui.Event, state *state.UIState) {
	currentRow := nc.Widget.SelectedRow
	baseTextLength := nc.BaseTextLengths[currentRow]
	showOtherWidget := false

	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		ret, _ := iptables.IptablesAddRule(iptables.ArraytToCmd(nc.ChainName, nc.RuleDesc, nc.BaseTextLengths))
		ui.Clear()
		ui.Render(state.Header, state.Footer, state.Tabpane)
		state.SetActive("chainList")
		msgBox := MsgBox(ret)
		state.Handlers["msgBox"] = msgBox
		state.SetActive("msgBox")
		state.Render()
	case "<Down>":
		nc.Widget.ScrollDown()
	case "<Up>":
		nc.Widget.ScrollUp()
	case "<Backspace>":
		if len(nc.RuleDesc[currentRow]) > baseTextLength {
			nc.RuleDesc[currentRow] = nc.RuleDesc[currentRow][:len(nc.RuleDesc[currentRow])-1]
		}
	case "c":
		if len(e.ID) == 1 {
			nc.RuleDesc[currentRow] += e.ID
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

func (nr *NewRuleObject) Render() {
	ui.Render(nr.Widget)
}
