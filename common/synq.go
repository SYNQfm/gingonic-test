package common

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/SYNQfm/SYNQ-Golang/synq"
)

type Cli struct {
	Command    string
	Timeout    int
	Simulate   bool
	Limit      int
	Filter     string
	FilterType string
	CacheDir   string
	CmdMsg     string
	CmdDef     string
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
		cd = flag.String(
			"cache_dir",
			"",
			"cache dir to use for saved values",
		)
	)
	flag.Parse()
	c.Command = *cmd
	c.Timeout = *t
	c.Simulate = *s != "false"
	c.Limit = *l
	c.CacheDir = *cd
}

func (c *Cli) LoadVideos(sApi synq.Api) (videos []synq.Video, err error) {
	cache_file := ""
	if c.CacheDir != "" {
		cache_file = c.CacheDir + "/" + c.FilterType + ".json"
		if _, e := os.Stat(cache_file); e == nil {
			log.Printf("loading from cached file %s\n", cache_file)
			bytes, _ := ioutil.ReadFile(cache_file)
			json.Unmarshal(bytes, &videos)
		}
	}
	if len(videos) == 0 {
		log.Printf("querying '%s'\n", c.Filter)
		videos, err = sApi.Query(c.Filter)
		if err != nil {
			return videos, err

		}
		if cache_file != "" {
			data, _ := json.Marshal(&videos)
			log.Printf("saving %d videos to %s\n", len(videos), cache_file)
			ioutil.WriteFile(cache_file, data, 0755)
		}
	}
	return videos, nil
}
