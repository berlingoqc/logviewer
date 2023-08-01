package views

import (
	"fmt"

	"github.com/rivo/tview"
)

type statusBarView struct {
	textElement *tview.TextView
}

func (sbv *statusBarView) GetPrimitive() tview.Primitive {
	return sbv.textElement
}

func GetStatusBarView(app *TuiApplication) *statusBarView {
	bar := &statusBarView{
		textElement: tview.NewTextView(),
	}

	bar.textElement.SetBorder(false)

	txt := ""

	for i, p := range app.mainPagesView.pages {
		txt += fmt.Sprintf("%d - %s ", i, p.GetName())
	}

	bar.textElement.SetText(txt)

	return bar
}
