package synq

import (
	"log"
	"os"

	"github.com/SYNQfm/SYNQ-Golang/synq"
)

func SetupSynq() synq.Api {
	key := os.Getenv("SYNQ_API_KEY")
	if key == "" {
		log.Println("WARNING : no Synq API key specified")
	}
	sApi := synq.New(key)
	url := os.Getenv("SYNQ_API_URL")
	if url != "" {
		sApi.Url = url
	}
	return sApi
}
