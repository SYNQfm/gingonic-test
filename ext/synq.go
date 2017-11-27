package ext

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/SYNQfm/SYNQ-Golang/synq"
	"github.com/SYNQfm/helpers/common"
)

func LoadVideosByQuery(query, name string, c common.Cacheable, api synq.Api) (videos []synq.Video, err error) {
	cacheFile := c.GetCacheFile(name)
	if cacheFile != "" {
		if _, e := os.Stat(cacheFile); e == nil {
			log.Printf("loading from cached file %s\n", cacheFile)
			bytes, _ := ioutil.ReadFile(cacheFile)
			json.Unmarshal(bytes, &videos)
		}
	}
	if len(videos) == 0 {
		log.Printf("querying '%s'\n", query)
		videos, err = api.Query(query)
		if err != nil {
			return videos, err

		}
		if cacheFile != "" {
			data, _ := json.Marshal(&videos)
			log.Printf("saving %d videos to %s\n", len(videos), cacheFile)
			ioutil.WriteFile(cacheFile, data, 0755)
		}
	}
	return videos, nil
}

func LoadVideo(id string, c common.Cacheable, api synq.Api) (video synq.Video, err error) {
	cacheFile := c.GetCacheFile(id)
	if _, err := os.Stat(cacheFile); err == nil {
		bytes, _ := ioutil.ReadFile(cacheFile)
		json.Unmarshal(bytes, &video)
	} else {
		// need to use the v1 api to get the raw video data
		log.Printf("Getting video %s", id)
		video, e := api.GetVideo(id)
		if e != nil {
			return video, e
		}
		if cacheFile != "" {
			bytes, _ := json.Marshal(video)
			ioutil.WriteFile(cacheFile, bytes, 0755)
		}
	}
	return video, err
}
