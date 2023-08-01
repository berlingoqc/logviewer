package views

import "github.com/rivo/tview"

type MainPagesView struct {
	view *tview.Pages

	pages []PageView
}

func (mpv *MainPagesView) GetPrimitive() tview.Primitive {
	return mpv.view
}

func GetMainPages(app *TuiApplication) *MainPagesView {
	pages := &MainPagesView{
		view:  tview.NewPages(),
		pages: make([]PageView, 1),
	}

	homePage := getPageHomeView()

	pages.view.SetBorder(false)
	pages.view.AddPage(homePage.GetName(), homePage.GetPrimitive(), true, true)
	pages.pages[0] = homePage

	for _, configPage := range app.config.LogPages {
		logPage := getPageLogView(app, configPage)

		pages.view.AddPage(logPage.GetName(), logPage.GetPrimitive(), true, true)
		pages.pages = append(pages.pages, logPage)
	}

	return pages
}
