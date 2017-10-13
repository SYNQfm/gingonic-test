package common

import (
	"fmt"
	"strings"
	"time"
)

type Ret struct {
	Label    string
	CountMap map[string]int
	Error    error
	Start    time.Time
}

func NewRet(label string) Ret {
	return Ret{
		Label:    label,
		CountMap: make(map[string]int),
		Error:    nil,
		Start:    time.Now(),
	}
}

func ParseType(type_ string) string {
	t := strings.ToLower(type_)
	switch t {
	case "error":
		t = "errored"
	case "skip":
		t = "skipped"
	case "ct":
		t = "count"
	}
	return t
}

func (r *Ret) Add(type_ string) {
	t := ParseType(type_)
	if _, ok := r.CountMap[t]; !ok {
		r.CountMap[t] = 0
	}
	r.CountMap[t] = r.CountMap[t] + 1
}

func (r *Ret) IsErrored() bool {
	return r.Value("total") == r.Value("errored")
}

func (r *Ret) Value(type_ string) int {
	t := ParseType(type_)
	c, ok := r.CountMap[t]
	if !ok {
		c = 0
	}
	return c
}

func (r *Ret) Eq(type_ string, ct int) bool {
	return r.Value(type_) == ct
}

func (r *Ret) Gte(type_ string, ct int) bool {
	return r.Value(type_) >= ct
}

func (r *Ret) Gt(type_ string, ct int) bool {
	return r.Value(type_) > ct
}

func (r *Ret) Lte(type_ string, ct int) bool {
	return r.Value(type_) <= ct
}

func (r *Ret) Lt(type_ string, ct int) bool {
	return r.Value(type_) < ct
}

func (r *Ret) String() string {
	str := fmt.Sprintf("for %s, processed %d/%d", r.Label, r.Value("ct"), r.Value("total"))
	for k, v := range r.CountMap {
		if k == "count" || k == "total" {
			continue
		}
		str = str + fmt.Sprintf(", %s %d", k, v)
	}
	str = str + "\n"
	dur := time.Since(r.Start)
	ms := int(dur / time.Millisecond)
	str = str + fmt.Sprintf("took %d ms", ms)
	return str
}
