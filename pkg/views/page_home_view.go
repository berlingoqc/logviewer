package views

import "github.com/rivo/tview"

const HomePageName = "homePage"

type pageHomeView struct {
	box *tview.TextView
}

func (phv *pageHomeView) GetPrimitive() tview.Primitive {
	return phv.box
}

func (phv *pageHomeView) GetName() string {
	return HomePageName
}

func getPageHomeView() *pageHomeView {
	page := &pageHomeView{}
	page.box = tview.NewTextView().SetText("cacacacacaca")

	return page
}
