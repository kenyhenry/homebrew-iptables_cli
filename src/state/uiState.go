package state

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type WidgetHandler interface {
	HandleEvent(e ui.Event, state *UIState)
	Render()
}

type UIState struct {
	ActiveHandler WidgetHandler
	Handlers      map[string]WidgetHandler
	Header        *widgets.Paragraph
	Footer        *widgets.Paragraph
	Tabpane       *widgets.TabPane
}

func NewUIState(header, footer *widgets.Paragraph, tabpane *widgets.TabPane) *UIState {
	return &UIState{
		Handlers: make(map[string]WidgetHandler),
		Header:   header,
		Footer:   footer,
		Tabpane:  tabpane,
	}
}
func (state *UIState) SetActive(name string) {
	if handler, exists := state.Handlers[name]; exists {
		state.ActiveHandler = handler
		state.Render()
	}
}

func (state *UIState) Default() {
	// state.
}

func (state *UIState) HandleEvent(e ui.Event, statePrm *UIState) {
	if state.ActiveHandler != nil {
		state.ActiveHandler.HandleEvent(e, state)
	}
}

func (state *UIState) Render() {
	if state.ActiveHandler != nil {
		state.ActiveHandler.Render()
	}
}
