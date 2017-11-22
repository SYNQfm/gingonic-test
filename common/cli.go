package common

import (
	"flag"
	"os"
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

func (c *Cli) Parse() {
	c.Flag.Parse(os.Args[1:])
	c.Command = *c.Args["command"].(*string)
	c.Timeout = *c.Args["timeout"].(*int)
	c.Simulate = *c.Args["simulate"].(*string) != "false"
	c.Limit = *c.Args["limit"].(*int)
	c.CacheDir = *c.Args["cache_dir"].(*string)
}
