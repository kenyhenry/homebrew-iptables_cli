package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
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

	help := `Helper

			Rules
			<a> 			add new rule on current chain
			<e> 			edit rule selected on current chain
			<d> 			delete the rule selected on the current chain

			Chain
			<c> 			add new chain
			<D> 			delete the current chain
			<F> 			flush or delete all rule inside of the current chain
			<E> 			rename the current chain
			<P> 			set current chain policy only supported ("DROP", "INPUT", "FORWARD", "ACCEPT", "OUTPUT")

			General
			<Enter> 		on tape on rule can move rule tape again to valid new rule emplacement
			<Up & Down> 	to move in the list of rule
			<Left & Right> 	to navigate into chain
			<ctrl-c> 		quit iptables_cli
	`

	termWidth, termHeight := ui.TerminalDimensions()

	header := widgets.NewParagraph()
	header.Text = logo
	header.SetRect(0, 0, termWidth, 10)
	header.TextStyle.Fg = ui.ColorGreen

	footer := widgets.NewParagraph()
	footer.Text = "<a> new rule | <e> edit rule | <c> add chain | <d> delete rule | <D> delete chain | <P> set chain policy | <E> rename chain | <F> flush chain | <ctrl-c> quit"
	footer.SetRect(0, termHeight-3, termWidth, termHeight)
	footer.Border = true
	footer.WrapText = true
	footer.TextStyle.Fg = ui.ColorCyan

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

	renderTab()
	newChain := NewChain()
	state.handlers["newChain"] = newChain
	ui.Render(header, footer, tabpane)
	renderTab()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		switch e.ID {

		case "<C-c>":
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
		case "c":
			if isDifferentFromKnownHandlers(state) {
				state.SetActive("newChain")
			} else {
				state.HandleEvent(e, state)
			}
		case "E":
			if isDifferentFromKnownHandlers(state) {
				if !ContainString(chain[tabpane.ActiveTabIndex], []string{"DROP", "INPUT", "FORWARD", "ACCEPT", "OUTPUT"}) {
					// TODO : from oldchainename and newchainname send iptables command
					state.SetActive("newChain")
				}
			} else {
				state.HandleEvent(e, state)
			}
		case "h":
			helperBox := Helper(help)
			ui.Render(helperBox.Widget)
		default:
			state.HandleEvent(e, state)

		}
	}
}

func isDifferentFromKnownHandlers(state *UIState) bool {
	_, isNewRule := state.activeHandler.(*NewRuleObject)
	_, isEditRule := state.activeHandler.(*EditRuleObject)
	_, isNewChain := state.activeHandler.(*NewChainObject)

	return !isNewRule && !isEditRule && !isNewChain
}
