package views

import (
	"context"
	"log"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/log/printer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type logView struct {
	app    *TuiApplication
	parent *tview.Flex
	tv     *tview.List
	fields *tview.DeepList
	result client.LogSearchResult
}

func (tv logView) Display(ctx context.Context, result client.LogSearchResult) error {
	go printer.WrapIoWritter(ctx, result, tv.tv, func() {
		// TODO: scroll to end if we are not scroll up
		tv.app.tapp.QueueUpdateDraw(func() {
		})
	})

	return nil
}

func createLogTextView(app *TuiApplication, name string) *logView {
	parentFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	//tv := tview.NewTextView().
	//	SetTextAlign(tview.AlignLeft).
	//	SetScrollable(true)

	tv := tview.NewList()

	tv.SetBorder(true)
	tv.SetTitle(name)

	parentFlex.AddItem(tv, 0, 100, false)

	wrapper := new(logView)
	wrapper.app = app
	wrapper.tv = tv
	wrapper.parent = parentFlex

	parentFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlB {
			if wrapper.fields == nil {
				fields, _, err := wrapper.result.GetFields()
				if err != nil {
					log.Println(err.Error())
					return event
				}

				listPrimitive := tview.NewDeepList()

				i := 0
				for k, values := range fields {
					ii := i
					listPrimitive.AddItem(k, "", rune(0), func() {
						listPrimitive.ToggleSubListDisplay(ii)
					})
					for _, v := range values {
						listPrimitive.AddSubItem(v, "", rune(0), false, nil)
					}
					i += 1
				}

				wrapper.fields = listPrimitive
				parentFlex.AddItem(wrapper.fields, 0, 40, true)
			} else {
				wrapper.parent.RemoveItem(wrapper.fields)
				wrapper.fields = nil
			}
			return nil
		}
		return event
	})

	return wrapper
}
