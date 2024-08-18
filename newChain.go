package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type NewChainObject struct {
	Widget   *widgets.Paragraph
	BaseText string
}

func NewChain() *NewChainObject {
	msgBox := widgets.NewParagraph()

	baseText := "chainName : "

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = true
	msgBox.Text = baseText
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)

	return &NewChainObject{
		Widget:   msgBox,
		BaseText: baseText,
	}
}

func (nc *NewChainObject) HandleEvent(e ui.Event, state *UIState) {
	showOtherWidget := false
	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		// TODO : send command to add new chain
		// TODO: ret of command
		ret := "test"
		msgBox := MsgBox(ret)
		state.handlers["msgBox"] = msgBox
		state.SetActive("msgBox")
		state.Render()
	case "<Backspace>":
		if len(nc.Widget.Text) > len(nc.BaseText) {
			nc.Widget.Text = nc.Widget.Text[:len(nc.Widget.Text)-1]
		}
	default:
		if len(e.ID) == 1 {
			nc.Widget.Text += e.ID
		}
	}
	if !showOtherWidget {
		ui.Render(nc.Widget)
	}
}

func (nr *NewChainObject) Render() {
	ui.Render(nr.Widget)
}
