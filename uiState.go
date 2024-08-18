package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type WidgetHandler interface {
	HandleEvent(e ui.Event)
	Render()
}

type UIState struct {
	activeHandler WidgetHandler
	handlers      map[string]WidgetHandler
	header        *widgets.Paragraph
	footer        *widgets.Paragraph
	tabpane       *widgets.TabPane
}

func NewUIState(header, footer *widgets.Paragraph, tabpane *widgets.TabPane) *UIState {
	return &UIState{
		handlers: make(map[string]WidgetHandler),
		header:   header,
		footer:   footer,
		tabpane:  tabpane,
	}
}
func (state *UIState) SetActive(name string) {
	if handler, exists := state.handlers[name]; exists {
		state.activeHandler = handler
		state.Render()
	}
}

func (state *UIState) Default() {
	// state.
}

func (state *UIState) HandleEvent(e ui.Event) {
	if state.activeHandler != nil {
		state.activeHandler.HandleEvent(e)
	}
}

func (state *UIState) Render() {
	if state.activeHandler != nil {
		state.activeHandler.Render()
	}
}
