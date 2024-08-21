package graphical

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/kenyhenry/iptables_cli/events"
	"github.com/kenyhenry/iptables_cli/state"
)

type SelectBoxObject struct {
	Widget      *ui.Grid
	SelectItems []string
	Em          *events.EventManager
	EventName   string
}

func SelectBox(textInfo string, eventName string, selectItems []string, em *events.EventManager) *SelectBoxObject {
	selectBox := widgets.NewList()

	termWidth, termHeight := ui.TerminalDimensions()
	selectBox.SetRect(termWidth/2-5, termHeight/2-5, termWidth/2+5, termHeight/2+5)
	selectBox.Border = true
	selectBox.TitleStyle.Fg = 3
	selectBox.WrapText = true
	selectBox.TextStyle = ui.NewStyle(ui.ColorMagenta)
	selectBox.Rows = selectItems

	paragraph := widgets.NewParagraph()
	paragraph.Text = textInfo
	paragraph.BorderStyle.Fg = ui.ColorRed
	paragraph.SetRect(0, 0, 3, 3)
	paragraph.Border = true

	grid := ui.NewGrid()
	grid.SetRect(termWidth/2-22, termHeight/2-10, termWidth/2+22, termHeight/2+10)
	grid.Set(
		ui.NewRow(0.2, paragraph),
		ui.NewRow(0.8, selectBox),
	)

	return &SelectBoxObject{
		Widget:      grid,
		SelectItems: selectItems,
		Em:          em,
		EventName:   eventName,
	}
}

func getListItem(grid *ui.Grid) *widgets.List {
	gridItem := grid.Items[len(grid.Items)-1]
	if gridItem.IsLeaf {
		widget, ok := gridItem.Entry.(*widgets.List)
		if ok {
			return widget
		}
	}
	return widgets.NewList()
}

func (nc *SelectBoxObject) HandleEvent(e ui.Event, state *state.UIState) {
	showOtherWidget := false

	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		ui.Clear()
		ui.Render(state.Header, state.Footer, state.Tabpane)
		state.SetActive("chainList")
		list := getListItem(nc.Widget)
		nc.Em.Emit(events.Event{Name: nc.EventName, Data: nc.SelectItems[list.SelectedRow]})
		// time.Sleep(time.Second)
	case "<Down>":
		list := getListItem(nc.Widget)
		list.ScrollDown()
	case "<Up>":
		list := getListItem(nc.Widget)
		list.ScrollUp()
	}
	list := getListItem(nc.Widget)
	list.Rows = nc.SelectItems
	if !showOtherWidget {
		ui.Render(nc.Widget)
	}
}

func (nr *SelectBoxObject) Render() {
	ui.Render(nr.Widget)
}
