package common

import (
	"fmt"
	"strings"
	"time"
)

type Ret struct {
	Label   string
	Total   int
	Count   int
	Skipped int
	Already int
	Errored int
	Error   error
	Start   time.Time
}

func NewRet(label string) Ret {
	return Ret{
		Label:   label,
		Count:   0,
		Total:   0,
		Already: 0,
		Skipped: 0,
		Errored: 0,
		Error:   nil,
		Start:   time.Now(),
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
	case "already":
		r.Already = r.Already + 1
	case "skipped",
		"skip":
		r.Skipped = r.Skipped + 1
	}
}

func (r *Ret) IsErrored() bool {
	return r.Total == r.Errored
}

func (r *Ret) String() string {
	return fmt.Sprintf("for %s, processed %d/%d, skipped %d, errored %d", r.Label, r.Count, r.Total, r.Skipped, r.Errored)
}
