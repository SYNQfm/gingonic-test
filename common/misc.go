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

// seconds by default
func (r *Ret) Taken(type_ ...time.Duration) int {
	t := time.Second
	if len(type_) > 0 {
		t = type_[0]
	}
	dur := time.Since(r.Start)
	v := int(dur / t)
	return v
}

func (r *Ret) Bytes() int {
	bytes, ok := r.CountMap["bytes"]
	if !ok {
		return 0
	}
	megs := bytes / (1000 * 1000)
	return megs
}

func (r *Ret) Speed() string {
	secs := r.Taken()
	megs := r.Bytes()
	if megs == 0 {
		return ""
	}
	speed := (float64(megs) / float64(secs) * 8)
	return fmt.Sprintf("%d megs (speed %f mbps)", megs, speed)
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
	ms := r.Taken(time.Millisecond)
	speed := r.Speed()
	if speed != "" {
		speed = ", " + speed
	}
	str = str + fmt.Sprintf("took %d ms%s", ms, speed)
	return str
}
