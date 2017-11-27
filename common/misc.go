package common

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Ret struct {
	Label    string
	Bytes    int
	CountMap map[string]int
	Error    error
	Start    time.Time
}

func NewRet(label string) Ret {
	return Ret{
		Label:    label,
		Bytes:    0,
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

func (r *Ret) AddBytes(bytes int) {
	r.Bytes = r.Bytes + bytes
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

func Label(dur time.Duration) string {
	if dur == time.Minute {
		return "mins"
	} else if dur == time.Second {
		return "sec"
	} else if dur == time.Millisecond {
		return "ms"
	} else {
		return "ns"
	}
}

// This will determine the right value to use
func (r *Ret) Taken(tDur ...time.Duration) (int, string) {
	var t time.Duration
	dur := time.Since(r.Start)
	if len(tDur) > 0 {
		t = tDur[0]
	}
	if dur >= 1000*time.Second {
		t = time.Minute
	} else if dur >= 10000*time.Millisecond {
		t = time.Second
	} else if dur >= 10000*time.Nanosecond {
		t = time.Millisecond
	} else {
		t = time.Nanosecond
	}
	taken := int(dur / t)
	return taken, Label(t)
}

func (r *Ret) Megs() int {
	megs := r.Bytes / (1000 * 1000)
	return megs
}

func (r *Ret) Speed() string {
	secs, _ := r.Taken(time.Second)
	megs := r.Megs()
	if megs == 0 {
		return ""
	}
	speed := (float64(megs*8) / float64(secs))
	return fmt.Sprintf("%d megs (speed %.1f mbps)", megs, speed)
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
	s, l := r.Taken()
	speed := r.Speed()
	if speed != "" {
		speed = ", " + speed
	}
	str = str + fmt.Sprintf("took %d %s%s", s, l, speed)
	return str
}

// return 32 bytes into 36 bytes
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func ConvertToUUIDFormat(uuid string) string {
	if len(uuid) == 36 && strings.Count(uuid, "-") == 4 {
		return uuid
	}
	if len(uuid) != 32 {
		log.Printf("invalid uuid %s\n", uuid)
		return uuid
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])
}
