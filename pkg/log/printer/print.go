package printer

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"git.tmaws.io/tmconnect/logexplorer/pkg/log/client"
)

type PrintPrinter struct{}

func (pp PrintPrinter) Append(result client.LogSearchResult) error {

    entries, err := result.GetEntries()
    if err != nil { return err }

	for _, entry := range entries {
		fmt.Printf("[%v][%s] %s \n", entry.Timestamp, entry.Level, entry.Message)
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

        for newEntries := range newEntriesChannel {
        	for _, entry := range newEntries {
		        fmt.Printf("[%v][%s] %s \n", entry.Timestamp, entry.Level, entry.Message)
	        }
        }
    }

	return nil
}
