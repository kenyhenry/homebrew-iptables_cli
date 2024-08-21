package graphical

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kenyhenry/iptables_cli/src/state"
)

type HelperObject struct {
	Widget *widgets.Paragraph
}

func Helper(message string) *HelperObject {
	msgBox := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()
	msgBox.SetRect(termWidth/2-30, termHeight/2-13, termWidth/2+30, termHeight/2+13)
	msgBox.Border = true
	msgBox.TitleStyle.Fg = 3
	msgBox.WrapText = true
	msgBox.Text = message
	msgBox.TextStyle = ui.NewStyle(ui.ColorCyan)

	return &HelperObject{
		Widget: msgBox,
	}
}

func (nc *HelperObject) HandleEvent(e ui.Event, state *state.UIState) {
	// Do Nothing
}

func (nr *HelperObject) Render() {
	ui.Render(nr.Widget)
}
