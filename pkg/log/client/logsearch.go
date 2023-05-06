package client

import "github.com/berlingoqc/logexplorer/pkg/ty"

type SearchRange struct {
	Lte  ty.Opt[string] `json:"lte"`
	Gte  ty.Opt[string] `json:"gte"`
	Last ty.Opt[string] `json:"last"`
}

type RefreshOptions struct {
	Duration ty.Opt[string] `json:"duration,omitempty"`
}

type TagExtraction struct {
	Regex ty.Opt[string] `json:"regex,omitempty"`
}

type PrinterOptions struct {
	Template ty.Opt[string] `json:"template,omitempty"`
}

type LogSearch struct {
	// Current filterring tags
	Tags ty.MS `json:"tags,omitempty"`
	// Extra rules for filtering tags
	TagsCondition ty.MS `json:"tagsCondition,omitempty"`

	// Range of the log query to do , depends of the system for full availability
	Range SearchRange `json:"range,omitempty"`

	// Max size of the request
	Size ty.Opt[int] `json:"size,omitempty"`

	// Refresh options for live data
	Refresh RefreshOptions `json:"refresh,omitempty"`

	// Options to configure the implementation with specific configuration for the search
	Options ty.MI `json:"options,omitempty"`

	// Extra fields for tag extraction for system without tagging of log entry
	TagExtraction TagExtraction `json:"tagExtraction,omitempty"`

	PrinterOptions PrinterOptions `json:"printerOptions,omitempty"`
}

func (lr *LogSearch) MergeInto(logSeach *LogSearch) error {

	if lr.Tags == nil {
		lr.Tags = ty.MS{}
	}
	if lr.TagsCondition == nil {
		lr.TagsCondition = ty.MS{}
	}
	if lr.Options == nil {
		lr.Options = ty.MI{}
	}

	lr.Tags = ty.MergeM(lr.Tags, logSeach.Tags)
	lr.TagsCondition = ty.MergeM(lr.TagsCondition, logSeach.TagsCondition)

	lr.Size.Merge(&logSeach.Size)
	lr.Refresh.Duration.Merge(&logSeach.Refresh.Duration)

	lr.Options = ty.MergeM(lr.Options, logSeach.Options)
	lr.TagExtraction.Regex.Merge(&logSeach.TagExtraction.Regex)

	lr.PrinterOptions.Template.Merge(&logSeach.PrinterOptions.Template)

	return nil
}
