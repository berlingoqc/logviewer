package views

/*
import (
	"context"
	"log"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/log/printer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	notSelectedColor = tcell.ColorGrey
	selectedColor    = tcell.ColorGreen
)

// Return the queryBox to display one output of logs
func getQueryBox(app *tview.Application, searchesId []string) (*tview.Flex, map[string]*logView, error) {

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	tviewWrappers := make(map[string]*logView)

	elementProportion := 100 / len(searchesId)

	for _, v := range searchesId {
		wrapper := createLogTextView(app, v)
		flex.AddItem(wrapper.parent, 0, elementProportion, false)
		tviewWrappers[v] = wrapper
	}

	return flex, tviewWrappers, nil
}

func RunQueryViewApp(config config.ContextConfig, searchIds []string) error {

	app := tview.NewApplication().EnableMouse(true)

	grid, wrappers, err := getQueryBox(app, searchIds)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	//defer cancel()

	for k, v := range wrappers {
		result, err := searchFactory.GetSearchResult(k, []string{}, client.LogSearch{})
		v.result = result
		if err != nil {
			return err
		}

		err = v.Display(ctx, result)
		if err != nil {
			panic(err)
		}
	}

	if err := app.SetRoot(grid, true).Run(); err != nil {
		return err
	}

	return nil
}
*/
