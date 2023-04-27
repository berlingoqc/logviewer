package printer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"text/template"
	"time"

	"git.tmaws.io/tmconnect/logexplorer/pkg/log/client"
)

type PrinterOptions struct{
    Template string
}

func formatDate(layout string, t time.Time) string {
    return t.Format(layout)
}

type PrintPrinter struct{
    Options PrinterOptions
}

func (pp PrintPrinter) Append(result client.LogSearchResult) error {

    template, err3 := template.New("print_printer").Funcs(
        template.FuncMap{"Format": formatDate },
    ).Parse(pp.Options.Template + "\n")
    if err3 != nil { return err3 }

    entries, err := result.GetEntries()
    if err != nil { return err }

	for _, entry := range entries {
        template.Execute(os.Stdout, entry)
	}

    ctx := context.Background()

    ctx, cancel := context.WithCancel(ctx)

    defer cancel()

    newEntriesChannel, err := result.OnChange(ctx)

    if err != nil { return err }

    if newEntriesChannel != nil {

        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func() {
            for range c {
                fmt.Println("signal to end")
                cancel()
                return;
            }
        }()

        for newResult := range newEntriesChannel {
            entries, _ := newResult.GetEntries()
        	for _, entry := range entries {
                template.Execute(os.Stdout, entry)
	        }
        }
    }

	return nil
}
