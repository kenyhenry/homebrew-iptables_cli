package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type SelectBoxObject struct {
	Widget      *widgets.List
	SelectItems []string
	Em          *EventManager
}

func SelectBox(textInfo string, eventName string, selectItems []string, em *EventManager) *SelectBoxObject {
	selectBox := widgets.NewList()

	termWidth, termHeight := ui.TerminalDimensions()
	selectBox.SetRect(termWidth/2-25, termHeight/2-5, termWidth/2+25, termHeight/2+5)
	selectBox.Border = true
	selectBox.TitleStyle.Fg = 3
	selectBox.WrapText = false
	selectBox.TextStyle = ui.NewStyle(ui.ColorCyan)
	selectBox.Rows = selectItems

	return &SelectBoxObject{
		Widget:      selectBox,
		SelectItems: selectItems,
		Em:          em,
	}
}

func (nc *SelectBoxObject) HandleEvent(e ui.Event, state *UIState) {
	showOtherWidget := false

	switch e.ID {
	case "<Enter>":
		showOtherWidget = true
		// TODO : send command to add new rule
		ui.Clear()
		ui.Render(state.header, state.footer, state.tabpane)
		state.SetActive("chainList")
		nc.Em.Emit(Event{Name: "deleteChain", Data: nc.SelectItems[nc.Widget.SelectedRow]})
		// time.Sleep(time.Second)
		// TODO should send to widget the result
	case "<Down>":
		nc.Widget.ScrollDown()
	case "<Up>":
		nc.Widget.ScrollUp()
	}
	nc.Widget.Rows = nc.SelectItems
	if !showOtherWidget {
		ui.Render(nc.Widget)
	}
}

func (nr *SelectBoxObject) Render() {
	ui.Render(nr.Widget)
}
