package views

import "github.com/berlingoqc/logexplorer/pkg/log/client"


type textViewLogPrinter struct {

}

func (vlp textViewLogPrinter) Append(entry []client.LogEntry) error {


    return nil
}
