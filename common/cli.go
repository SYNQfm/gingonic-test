package common

import (
	"flag"
	"os"
	"strings"
)

type Cli struct {
	Command    string
	Timeout    int
	Simulate   bool
	Limit      int
	Filter     string
	FilterType string
	CacheDir   string
	Flag       *flag.FlagSet
	Args       map[string]interface{}
}

func NewCli() Cli {
	cli := Cli{
		Flag: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
	cli.Args = make(map[string]interface{})
	return cli
}

func (c *Cli) DefaultSetup(msg, def string) {
	c.String("command", msg, def)
	c.String("simulate", "true", "simulate the transaction")
	c.Int("timeout", 120, "timeout to use for API call, in seconds, defaults to 120")
	c.Int("limit", 10, "number of actions to run")
	c.String("cache_dir", "", "cache dir to use for saved values")
}

func (c *Cli) String(name, def, desc string) {
	c.Args[name] = c.Flag.String(name, def, desc)
}
func (c *Cli) Int(name string, def int, desc string) {
	c.Args[name] = c.Flag.Int(name, def, desc)
}

func (c *Cli) GetString(name string) string {
	if _, ok := c.Args[name]; !ok {
		return ""
	}
	return *c.Args[name].(*string)
}

func (c *Cli) GetInt(name string) int {
	if _, ok := c.Args[name]; !ok {
		return -1
	}
	return *c.Args[name].(*int)
}

func (c *Cli) Parse(args ...[]string) {
	var a []string
	if len(args) > 0 {
		a = args[0]
	} else {
		i := 0
		for idx, val := range os.Args {
			if idx == 0 || strings.Contains(val, "-test.run") {
				i = i + 1
				continue
			} else {
				break
			}
		}
		a = os.Args[i:]
	}
	c.Flag.Parse(a)
	c.Command = c.GetString("command")
	c.Timeout = c.GetInt("timeout")
	c.Simulate = c.GetString("simulate") != "false"
	c.Limit = c.GetInt("limit")
	c.CacheDir = c.GetString("cache_dir")
}
