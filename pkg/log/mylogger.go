package log

import (
	"io"
	"log"
	"os"
)

type MyLoggerOptions struct {
	// if we output to  stdout
	Stdout bool
	// Path of the file , if present log to it
	Path string
	// What level to log
	Level string
}

func ConfigureMyLogger(options *MyLoggerOptions) {
	var writer io.Writer

	if options.Path != "" {
		logfile, err := os.OpenFile(options.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		if options.Stdout {
			writer = io.MultiWriter(logfile, os.Stdout)
		} else {
			writer = logfile
		}
	} else if options.Stdout {
		writer = os.Stdout
	} else {
		writer, _ = os.OpenFile(os.DevNull, os.O_APPEND, 0666)
	}

	log.SetOutput(writer)

}
