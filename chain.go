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

	logo := ` I)iiii P)ppppp  T)tttttt   A)aa   B)bbbb   L)       E)eeeeee  S)ssss       C)ccc  L)       I)iiii
			   I)   P)    pp    T)     A)  aa  B)   bb  L)       E)       S)    ss     C)   cc L)         I)
			   I)   P)ppppp     T)    A)    aa B)bbbb   L)       E)eeeee   S)ss       C)       L)         I)
			   I)   P)          T)    A)aaaaaa B)   bb  L)       E)            S)     C)       L)         I)
			   I)   P)          T)    A)    aa B)    bb L)       E)       S)    ss     C)   cc L)         I)
			 I)iiii P)          T)    A)    aa B)bbbbb  L)llllll E)eeeeee  S)ssss       C)ccc  L)llllll I)iiii`

	termWidth, termHeight := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = logo
	header.SetRect(0, 0, termWidth, 10)
	header.TextStyle.Fg = ui.ColorGreen

	footer := widgets.NewParagraph()
	footer.Text = "<a> new rule | <e> edit rule | <q> quit"
	footer.SetRect(0, termHeight-3, termWidth, termHeight)
	footer.Border = true
	footer.WrapText = true
	footer.TextStyle.Bg = ui.ColorBlue

	var chain []string
	tabpane := widgets.NewTabPane(chain...)
	tabpane.SetRect(0, 10, termWidth, termHeight-3)
	tabpane.Border = true

	state := NewUIState(header, footer, tabpane)

	renderTab := func() {
		// TODO get rule of chain by sending command iptables -C
		chain = []string{"pierwszy", "drugi", "trzeci", "żółw", "four", "five"}
		tabpane.TabNames = chain
		if tabpane.ActiveTabIndex >= 0 && tabpane.ActiveTabIndex < len(chain) {
			chainlist := NewChainList(chain[tabpane.ActiveTabIndex])
			state.handlers["chainList"] = chainlist
			state.SetActive("chainList")
		}
	}

	newRule := NewRule()
	// msgBox := MsgBox()
	editRule := EditRule()
	newChain := NewChain()

	state.handlers["newRule"] = newRule
	state.handlers["editRule"] = editRule
	state.handlers["newChain"] = newChain

	renderTab()
	ui.Render(header, footer, tabpane)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {

		case "q", "<C-c>":
			return
		case "<Escape>":
			ui.Render(header, footer, tabpane)
			renderTab()
		case "<Left>":
			tabpane.FocusLeft()
			ui.Clear()
			state.Render()
			ui.Render(header, footer, tabpane)
			renderTab()
		case "<Right>":
			tabpane.FocusRight()
			ui.Clear()
			state.Render()
			ui.Render(header, footer, tabpane)
			renderTab()
		case "a":
			state.SetActive("newRule")
		case "e":
			state.SetActive("editRule")
		case "c":
			state.SetActive("newChain")
		default:
			state.HandleEvent(e)
		}
	}
}
