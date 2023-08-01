package views

import (
	"context"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/rivo/tview"
)

type pageLogView struct {
	view *tview.Flex

	logViews []*logView

	config LogPageConfig

	cancel context.CancelFunc
}

func (mpv *pageLogView) GetPrimitive() tview.Primitive {
	return mpv.view
}

func (mpv *pageLogView) GetName() string {
	return mpv.config.Name
}

func getPageLogView(app *TuiApplication, config LogPageConfig) *pageLogView {
	view := &pageLogView{
		view:     tview.NewFlex(),
		config:   config,
		logViews: make([]*logView, len(config.ContextIds)),
	}

	view.view.SetDirection(tview.FlexRow)

	ctx := context.Background()
	ctx, view.cancel = context.WithCancel(ctx)

	for i, contextId := range config.ContextIds {
		view.logViews[i] = createLogTextView(app, contextId)
		view.view.AddItem(view.logViews[i].parent, 0, 1, true)

		result, err := app.searchFactory.GetSearchResult(contextId, []string{}, client.LogSearch{})
		if err != nil {
			// replace the view by the error view
			panic(err)
		}

		err = view.logViews[i].Display(ctx, result)
		if err != nil {
			// same as previous error
			panic(err)
		}
	}

	return view
}
