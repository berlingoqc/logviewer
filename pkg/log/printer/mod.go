package printer

import (
	"context"
	"io"
	"text/template"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
)

type LogPrinter interface {
	Display(ctx context.Context, result client.LogSearchResult) error
}

func WrapIoWritter(ctx context.Context, result client.LogSearchResult, writer io.Writer, update func()) error {

	templateConfig := result.GetSearch().PrinterOptions.Template

	if templateConfig.Value == "" {
		templateConfig.S("[{{.Timestamp}}] {{.Message}}")
	}

	template, err3 := template.New("print_printer").Funcs(
		template.FuncMap{"Format": formatDate},
	).Parse(templateConfig.Value + "\n")
	if err3 != nil {
		return err3
	}

	entries, newEntriesChannel, err := result.GetEntries(ctx)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		template.Execute(writer, entry)
	}

	update()

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

			update()

			for entries := range newEntriesChannel {
				if len(entries) > 0 {
					for _, entry := range entries {
						template.Execute(writer, entry)
					}
					update()
				}
			}

		}()

	}

	return nil
}
