package common

import (
	"strings"
	"time"
)

type Ret struct {
	Label     string
	Total     int
	Count     int
	Skipped   int
	Errored   int
	isErrored bool
	Start     time.Date
}

func NewRet(label string) Ret {
	return Ret{
		Label:     label,
		Count:     0,
		Total:     0,
		Skipped:   0,
		Errored:   0,
		isErrored: false,
		Start:     time.Now(),
	}
}

func (r *Ret) Add(type_ string) {
	switch strings.ToLower(type_) {
	case "count":
		r.Count = r.Count + 1
	case "total":
		r.Total = r.Total + 1
	case "errored",
		"error":
		r.Errored = r.Errored + 1
	case "skipped",
		"skip":
		r.Skipped = r.Skipped + 1
	}
}
