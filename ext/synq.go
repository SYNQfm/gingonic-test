package ext

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/SYNQfm/SYNQ-Golang/synq"
	"github.com/SYNQfm/helpers/common"
)

func LoadVideos(c common.Cli, api synq.Api) (videos []synq.Video, err error) {
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
		videos, err = api.Query(c.Filter)
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
