package printer

import (
	"context"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
)

type LogPrinter interface {
	Display(ctx context.Context, result client.LogSearchResult) error
}
