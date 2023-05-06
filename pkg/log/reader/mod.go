package reader

import (
	"bufio"
	"context"
	"io"
	"regexp"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)

const maxBatchSize = 10

type ReaderLogResult struct {
	search  client.LogSearch
	scanner *bufio.Scanner
	closer  io.Closer

	// mutex because updated by goroutine
	entries []client.LogEntry
	tags    client.AvailableTags

	regexExtraction *regexp.Regexp
}

func (lr ReaderLogResult) GetSearch() *client.LogSearch {
	return &lr.search
}

func (lr *ReaderLogResult) parseLine(line string) bool {
	entry := client.LogEntry{
		Message: line,
		Fields:  make(ty.MI),
	}

	if lr.regexExtraction != nil {
		match := lr.regexExtraction.FindStringSubmatch(line)
		if len(match) > 0 {
			for i, name := range lr.regexExtraction.SubexpNames() {
				if i != 0 && name != "" {
					lr.tags.AddTag(name, match[i])
					entry.Fields[name] = match[i]
				}
			}
		}
	}

	for k, v := range lr.search.Tags {
		if vv, ok := entry.Fields[k]; ok {
			if v != vv {
				return false
			}
		} else {
			return false
		}
	}

	entry.Level = entry.Fields.GetString("Level")
	lr.entries = append(lr.entries, entry)
	return true
}

func (lr *ReaderLogResult) loadEntries() bool {
	lr.entries = make([]client.LogEntry, 0)

	for lr.scanner.Scan() {
		line := lr.scanner.Text()
		lr.parseLine(line)
	}
	return len(lr.entries) > 0
}

func (lr ReaderLogResult) GetEntries(ctx context.Context) ([]client.LogEntry, chan []client.LogEntry, error) {

	if lr.search.Refresh.Duration.Value == "" {
		lr.loadEntries()
		lr.closer.Close()
		return lr.entries, nil, nil
	} else {
		c := make(chan []client.LogEntry)

		go func() {
			defer close(c)
			defer lr.closer.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					{
						if lr.scanner.Scan() {
							if lr.parseLine(lr.scanner.Text()) {
								c <- []client.LogEntry{lr.entries[len(lr.entries)-1]}
							}
						}
					}
				}
			}
		}()

		return []client.LogEntry{}, c, nil
	}
}

func (lr ReaderLogResult) GetTags() (client.AvailableTags, error) {
	return lr.tags, nil
}

func GetLogResult(
	search client.LogSearch,
	scanner *bufio.Scanner,
	closer io.Closer,
) ReaderLogResult {

	var regexExtraction *regexp.Regexp
	if search.TagExtraction.Regex.Value != "" {
		regexExtraction = regexp.MustCompile(search.TagExtraction.Regex.Value)
	}

	result := ReaderLogResult{
		search:          search,
		scanner:         scanner,
		closer:          closer,
		regexExtraction: regexExtraction,
		tags:            make(client.AvailableTags),
	}

	return result
}
