package views

import (
	"context"
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/factory"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	notSelectedColor = tcell.ColorGrey
	selectedColor    = tcell.ColorGreen
)

type tviewWrapper struct {
	app *tview.Application
	tv  *tview.TextView
}

func (tv tviewWrapper) Display(ctx context.Context, result client.LogSearchResult) error {
	go printer.WrapIoWritter(ctx, result, tv.tv, func() {
		// TODO: scroll to end if we are not scroll up
		tv.app.QueueUpdateDraw(func() {
			tv.tv.ScrollToEnd()
		})
	})

	return nil
}

func createLogTextView(app *tview.Application, parent *tview.Grid, name string) *tview.TextView {
	tv := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true)

	tv.SetBorder(true)
	tv.SetBorderColor(notSelectedColor)
	tv.SetTitle(name)
	/*
		tv.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.Mouse {
				log.Println("focusing " + tv.GetTitle())
				app.SetFocus(tv)
				tv.SetBorderColor(selectedColor)
				return action, nil
			}
			return action, event
		})
		tv.SetBlurFunc(func() {
			// hide focus border
			log.Println("blurring " + tv.GetTitle())
			tv.SetBorderColor(notSelectedColor)
		})*/

	return tv
}

// Return the queryBox to display one output of logs
func getQueryBox(app *tview.Application, searchesId []string) (*tview.Grid, map[string]tviewWrapper, error) {

	grid := tview.NewGrid().
		SetColumns(0, 0).
		SetBorders(false)

	tviewWrappers := make(map[string]tviewWrapper)

	for i, v := range searchesId {
		primitive := createLogTextView(app, grid, v)
		grid.AddItem(primitive, 0, i, 1, 1, 0, 0, false)
		tviewWrappers[v] = tviewWrapper{tv: primitive, app: app}
	}

	return grid, tviewWrappers, nil
}

func RunQueryViewApp(config config.ContextConfig, searchIds []string) error {

	app := tview.NewApplication().EnableMouse(true)

	clientFactory, err := factory.GetLogClientFactory(config.Clients)
	if err != nil {
		return err
	}

	if len(searchIds) == 0 {
		return errors.New("required multiple searches for query")
	}

	searchFactory, err := factory.GetLogSearchFactory(clientFactory, config)
	if err != nil {
		return err
	}

	grid, wrappers, err := getQueryBox(app, searchIds)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	//defer cancel()

	for k, v := range wrappers {
		result, err := searchFactory.GetSearchResult(k, []string{}, client.LogSearch{})
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
