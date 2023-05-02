package printer

import (
	"context"
	"os"
	"text/template"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
)

type PrinterOptions struct {
	Template string `json="template"`
}

func formatDate(layout string, t time.Time) string {
	return t.Format(layout)
}

type PrintPrinter struct {
	Options PrinterOptions
}

func (pp PrintPrinter) Display(ctx context.Context, result client.LogSearchResult) error {

	template, err3 := template.New("print_printer").Funcs(
		template.FuncMap{"Format": formatDate},
	).Parse(pp.Options.Template + "\n")
	if err3 != nil {
		return err3
	}

	entries, newEntriesChannel, err := result.GetEntries(ctx)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		template.Execute(os.Stdout, entry)
	}

	if err != nil {
		return err
	}

	if newEntriesChannel != nil {

		/*
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			go func() {
				for range c {
					cancel()
					return
				}
			}()
		*/

		go func() {
			for entries := range newEntriesChannel {
				for _, entry := range entries {
					template.Execute(os.Stdout, entry)
				}
			}

		}()

	}

	return nil
}
