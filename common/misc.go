package common

import (
	"flag"
	"fmt"
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

type Cli struct {
	Command  string
	Timeout  int
	Simulate bool
	CmdMsg   string
	CmdDef   string
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

func (c *Cli) Parse() {
	var (
		cmd = flag.String(
			"command",
			c.CmdDef,
			c.CmdMsg,
		)
		s = flag.String(
			"simulate",
			"true",
			"simulate the transaction",
		)
		t = flag.Int(
			"timeout",
			120,
			"timeout to use for API call, in seconds, defaults to 120",
		)
		l = flag.Int(
			"limit",
			10,
			"number of actions to run",
		)
	)
	flag.Parse()
	c.Command = *cmd
	c.Timeout = *t
	c.Simulate = *s != "false"
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
