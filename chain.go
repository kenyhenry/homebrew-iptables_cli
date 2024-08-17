package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Chain() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	footer := widgets.NewParagraph()
	footer.Text = "<a> new rule | <q> quit"
	footer.SetRect(0, termHeight-3, termWidth, termHeight)
	footer.Border = true
	footer.WrapText = true
	footer.TextStyle.Bg = ui.ColorBlue

	// TODO each name is iptables chain
	chain := []string{"pierwszy", "drugi", "trzeci", "żółw", "four", "five"}
	tabpane := widgets.NewTabPane(chain...)
	tabpane.SetRect(0, 3, termWidth, termHeight-3)
	tabpane.Border = true

	var chainlist *NewChainlist

	renderTab := func() {
		if tabpane.ActiveTabIndex >= 0 && tabpane.ActiveTabIndex < len(chain) {
			// TODO get rule of chain by sending command iptables -C
			chainlist = NewChainList(chain[tabpane.ActiveTabIndex])
			ui.Render(chainlist.Widget)
		}
	}

	// Simulate rendering the active tab
	ui.Render(footer, tabpane)
	renderTab()

	msgBoxActivate := false
	msgBox := NewRule()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Escape>":
			if msgBoxActivate {
				msgBoxActivate = false
				ui.Clear()
				ui.Render(footer, tabpane)
				renderTab()
			}
		case "<Left>":
			tabpane.FocusLeft()
			ui.Clear()
			ui.Render(footer, tabpane)
			renderTab()
		case "<Right>":
			tabpane.FocusRight()
			ui.Clear()
			ui.Render(footer, tabpane)
			renderTab()
		case "<Down>":
			if msgBoxActivate {
				msgBox.RuleHandleEvent(e)
			} else {
				chainlist.HandleEvent(e)
			}
		case "<Up>":
			if msgBoxActivate {
				msgBox.RuleHandleEvent(e)
			} else {
				chainlist.HandleEvent(e)
			}
		case "a":
			if !msgBoxActivate {
				// TODO maybe message box in seperate file
				// TODO handle all key to tape inse of each item list
				msgBoxActivate = true
				ui.Clear()
				ui.Render(footer, tabpane)
				renderTab()
				ui.Render(msgBox.Widget)
			}
		default:
			if !msgBoxActivate {
				chainlist.HandleEvent(e)
			}
		}
	}
}
