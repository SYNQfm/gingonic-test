package ext

import (
	"log"
	"os"
	"strings"

	"github.com/SYNQfm/SYNQ-Golang/synq"
	"github.com/SYNQfm/helpers/common"
	"github.com/jmoiron/sqlx"
)

func SetupDB(def_url ...string) *sqlx.DB {
	dbaddr := common.GetDB(def_url...)
	return sqlx.MustConnect("postgres", dbaddr)
}

func SetupSynq() synq.Api {
	api := SetupSynqApi("v1")
	return api.(synq.Api)
}

func SetupSynqV2() synq.ApiV2 {
	api := SetupSynqApi("v2")
	return api.(synq.ApiV2)
}

func SetupSynqApi(version ...string) (api synq.ApiF) {
	key := os.Getenv("SYNQ_API_KEY")
	if key == "" {
		log.Println("WARNING : no Synq API key specified")
	}
	ver := os.Getenv("SYNQ_API_VERSION")
	if len(version) > 0 {
		ver = version[0]
	}
	if strings.Contains(key, ".") || ver == "v2" {
		api = synq.NewV2(key)
	} else {
		api = synq.NewV1(key)
	}
	url := os.Getenv("SYNQ_API_URL")
	if url != "" {
		api.SetUrl(url)
	}
	return api
}
