package graphical

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kenyhenry/iptables_cli/state"
)

type MsgBoxObject struct {
	Widget *widgets.Paragraph
}

func MsgBox(message string) *MsgBoxObject {
	msgBox := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = true
	msgBox.Text = message
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)

	return &MsgBoxObject{
		Widget: msgBox,
	}
}

func (nc *MsgBoxObject) HandleEvent(e ui.Event, state *state.UIState) {
	// Do Nothing
}

func (nr *MsgBoxObject) Render() {
	ui.Render(nr.Widget)
}
