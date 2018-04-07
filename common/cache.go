package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Cacheable interface {
	GetCacheFile(string) string
	GetCacheAge() time.Duration
}

func LoadFromCache(name string, c Cacheable, obj interface{}) bool {
	cacheFile := c.GetCacheFile(name)
	if cacheFile != "" {
		if i, e := os.Stat(cacheFile); e == nil {
			if c.GetCacheAge() < 0 || time.Since(i.ModTime()) >= c.GetCacheAge() {
				// do not load it as its too old
				return false
			}
			//log.Printf("loading from cached file %s\n", cacheFile)
			bytes, _ := ioutil.ReadFile(cacheFile)
			json.Unmarshal(bytes, obj)
			return true
		}
	}
	return false
}

func SaveToCache(name string, c Cacheable, obj interface{}) bool {
	cacheFile := c.GetCacheFile(name)
	if cacheFile != "" {
		data, _ := json.Marshal(obj)
		ioutil.WriteFile(cacheFile, data, 0755)
		return true
	}
	return false
}

func PurgeFromCache(name string, c Cacheable) bool {
	cacheFile := c.GetCacheFile(name)
	if _, e := os.Stat(cacheFile); e == nil {
		os.Remove(cacheFile)
		return true
	}
	return false
}
