package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Ret struct {
	Label    string
	CountMap map[string]int
	Error    error
	Start    time.Time
	DurMap   map[string]time.Duration
	BytesMap map[string]int64
	IdMap    map[string][]string
}

func NewRet(label string) Ret {
	return Ret{
		Label:    label,
		DurMap:   make(map[string]time.Duration),
		BytesMap: make(map[string]int64),
		CountMap: make(map[string]int),
		Error:    nil,
		Start:    time.Now(),
	}
}

func (r *Ret) AddBytes(bytes int64) {
	r.AddBytesFor("total", bytes)
}

func (r *Ret) AddBytesFor(key string, bytes int64) {
	t := ParseType(key)
	if _, ok := r.BytesMap[t]; !ok {
		r.BytesMap[t] = bytes
	} else {
		r.BytesMap[t] = r.BytesMap[t] + bytes
	}
}

func (r *Ret) AddDurFor(key string, dur time.Duration) {
	t := ParseType(key)
	if _, ok := r.DurMap[t]; !ok {
		r.DurMap[t] = dur
	} else {
		r.DurMap[t] = r.DurMap[t] + dur
	}
}

func (r *Ret) Add(type_ string) {
	r.AddFor(type_, 1)
}

func (r *Ret) AddFor(type_ string, ct int) {
	t := ParseType(type_)
	if _, ok := r.CountMap[t]; !ok {
		r.CountMap[t] = 0
	}
	r.CountMap[t] = r.CountMap[t] + ct
}

func (r *Ret) AddForSave(key string, id string) {
	t := ParseType(key)
	if _, ok := r.IdMap[t]; !ok {
		r.IdMap[t] = []string{}
	}
	r.IdMap[t] = append(r.IdMap[t], id)
}

func (r *Ret) AddError(err error) bool {
	if err != nil {
		r.Error = err
		r.Add("errored")
		return true
	}
	return false
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

func (r *Ret) Bytes(type_ string) int64 {
	t := ParseType(type_)
	b, ok := r.BytesMap[t]
	if !ok {
		b = int64(0)
	}
	return b
}

func (r *Ret) Duration(type_ string) time.Duration {
	t := ParseType(type_)
	d, ok := r.DurMap[t]
	if !ok {
		d = time.Duration(0)
	}
	return d
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

func (r *Ret) LimitReached(limit int) bool {
	return r.Gte("count", limit) || r.Gte("errored", limit)
}

// This will determine the right value to use
func (r *Ret) Taken(tDur ...time.Duration) (int, string) {
	dur := time.Since(r.Start)
	if len(tDur) > 0 {
		t := tDur[0]
		taken := int(dur / t)
		return taken, Label(t)
	} else {
		return DurVal(dur)
	}
}

func (r *Ret) Speed() string {
	if bytes, ok := r.BytesMap["total"]; ok {
		secs, _ := r.Taken(time.Second)
		b, l := BytesVal(bytes)
		speed := (float64(b*8) / float64(secs))
		return fmt.Sprintf("%d %s (speed %.2f %sps)", b, l, speed, l)
	} else {
		return ""
	}
}

func (r *Ret) String() string {
	ct := r.Value("ct")
	total := r.Value("total")
	str := "for " + r.Label
	if ct > 0 || total > 0 {
		str = str + fmt.Sprintf(", processed %d/%d", ct, total)
	}
	for k, v := range r.CountMap {
		if k == "count" || k == "total" {
			continue
		}
		str = str + fmt.Sprintf(", %s %d", k, v)
		bytes := r.Bytes(k)
		dur := r.Duration(k)
		bStr := ""
		dStr := ""
		if bytes > 0 {
			avg := bytes / int64(v)
			b, l := BytesVal(bytes)
			a, l2 := BytesVal(avg)
			bStr = fmt.Sprintf("%d %s, avg %d %s", b, l, a, l2)
		}
		if dur > 0 {
			d, l := DurVal(dur)
			avg := int(dur*time.Nanosecond) / v
			d2, l2 := DurVal(time.Duration(avg))
			dStr = fmt.Sprintf("duration %d %s, avg %d %s", d, l, d2, l2)
		}
		if bStr != "" && dStr != "" {
			str = str + " (" + bStr + ", " + dStr + ")"
		} else if bStr != "" {
			str = str + " ( " + bStr + ")"
		} else if dStr != "" {
			str = str + " ( " + dStr + ")"
		}
	}
	str = str + "\n"
	if r.Error != nil {
		str = str + r.GetErrorString()
	}
	s, l := r.Taken()
	speed := r.Speed()
	if speed != "" {
		speed = ", " + speed
	}
	str = str + fmt.Sprintf("took %d %s%s", s, l, speed)
	return str
}

func (r *Ret) GetErrorString() string {
	return fmt.Sprintf("Error occured : %s\n", r.Error.Error())
}

// This will save any id entries to disk
func (r *Ret) Save(dir string) error {
	for name, list := range r.IdMap {
		file := fmt.Sprintf("%s/%s.json", dir, name)
		log.Printf("Saving %d ids to %s", len(list), file)
		data, _ := json.Marshal(list)
		if err := ioutil.WriteFile(file, data, 0755); err != nil {
			return err
		}
	}
	return nil
}
