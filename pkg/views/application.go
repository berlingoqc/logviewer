package views

import (
	"github.com/berlingoqc/logviewer/pkg/log/config"
	"github.com/berlingoqc/logviewer/pkg/log/factory"
	"github.com/berlingoqc/logviewer/pkg/ty"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LogPageConfig struct {
	Name string

	ContextIds []string
}

type TuiApplicationConfig struct {
	ContextFilePath string
	LogPages        []LogPageConfig
}

type TuiApplication struct {
	config   TuiApplicationConfig
	contexts config.ContextConfig

	searchFactory *factory.LogSearchFactory

	statusBarView *statusBarView
	mainPagesView *MainPagesView

	tapp *tview.Application
}

func RunTuiApplication(config TuiApplicationConfig) error {
	app := &TuiApplication{
		config: config,
	}

	// Loading config file

	if err := ty.ReadJsonFile(config.ContextFilePath, &app.contexts); err != nil {
		panic(err)
	}

	// Creating factory for injection of stuff

	clientFactory, err := factory.GetLogClientFactory(app.contexts.Clients)
	if err != nil {
		return err
	}

	app.searchFactory, err = factory.GetLogSearchFactory(clientFactory, app.contexts)
	if err != nil {
		return err
	}

	// Create the application

	app.tapp = tview.NewApplication().EnableMouse(true)

	// create my root
	main := tview.NewGrid().SetRows(0, 1).SetColumns(0)
	main.SetBorder(false)

	app.mainPagesView = GetMainPages(app)
	app.statusBarView = GetStatusBarView(app)

	main.AddItem(app.mainPagesView.GetPrimitive(), 0, 0, 1, 1, 0, 0, true)
	main.AddItem(app.statusBarView.GetPrimitive(), 1, 0, 1, 1, 0, 0, false)

	app.tapp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if err := app.tapp.SetRoot(main, true).Run(); err != nil {
		return err
	}

	return nil
}
