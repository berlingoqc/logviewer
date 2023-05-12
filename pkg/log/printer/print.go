package printer

import (
	"context"
	"os"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
)

func formatDate(layout string, t time.Time) string {
	return t.Format(layout)
}

type PrintPrinter struct{}

func (pp PrintPrinter) Display(ctx context.Context, result client.LogSearchResult) error {

	return WrapIoWritter(ctx, result, os.Stdout, func() {})
}
