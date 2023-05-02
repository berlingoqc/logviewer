package views

import (
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/factory"
	"github.com/rivo/tview"
)


// Return the queryBox to display one output of logs
func getQueryBox(searchesId []string) (*tview.Grid, error) {

    newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignLeft).
			SetText(text)
	}

    grid := tview.NewGrid().
		SetColumns(0, 0).
		SetBorders(true)

    for i, v := range searchesId {
        primitive := newPrimitive("log " + v)
        grid.AddItem(primitive,0, i, 1, 1, 0, 0, false)
    }




    return grid, nil
}

func RunQueryViewApp(config config.ContextConfig) error {
    clientFactory, err := factory.GetLogClientFactory(config.Clients)
    if err != nil { return err }

    _, err = factory.GetLogSearchFactory(clientFactory, config)
    if err != nil { return err }

    searchesId := []string{"localSystem", "localSystem"}

    grid, err := getQueryBox(searchesId)
    if err != nil {
        return err
    }

	if err := tview.NewApplication().SetRoot(grid, true).Run(); err != nil {
		return err
	}

	return nil
}
