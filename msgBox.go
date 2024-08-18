package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type MsgBoxObject struct {
	Widget   *widgets.Paragraph
	BaseText string
}

func MsgBox() *MsgBoxObject {
	msgBox := widgets.NewParagraph()

	baseText := "chainName : "

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = true
	msgBox.Text = baseText
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)

	return &MsgBoxObject{
		Widget:   msgBox,
		BaseText: baseText,
	}
}

func (nc *MsgBoxObject) HandleEvent(e ui.Event) {
	switch e.ID {
	case "<Enter>":
		// TODO : send command to add new rule
		// nc.IsMoving = !nc.IsMoving
	case "<Backspace>":
		if len(nc.Widget.Text) > len(nc.BaseText) {
			nc.Widget.Text = nc.Widget.Text[:len(nc.Widget.Text)-1]
		}
	default:
		if len(e.ID) == 1 {
			nc.Widget.Text += e.ID
		}
	}
	ui.Render(nc.Widget)
}

func (nr *MsgBoxObject) Render() {
	ui.Render(nr.Widget)
}
