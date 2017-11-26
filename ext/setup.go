package ext

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SYNQfm/SYNQ-Golang/synq"
	"github.com/SYNQfm/helpers/common"
	"github.com/jmoiron/sqlx"
)

type ApiSetup struct {
	Key     string
	Version string
	Url     string
}

func SetupDB(def_url ...string) *sqlx.DB {
	dbaddr := common.GetDB(def_url...)
	return sqlx.MustConnect("postgres", dbaddr)
}

func SetupSynq() synq.Api {
	api := SetupSynqApi()
	return api.(synq.Api)
}

func SetupSynqV2() synq.ApiV2 {
	config := GetSetupByEnv("v2")
	api := SetupSynqApi(config)
	return api.(synq.ApiV2)
}

func GetSetupByEnv(version string) ApiSetup {
	key := os.Getenv(fmt.Sprintf("SYNQ_API%s_KEY", version))
	url := os.Getenv(fmt.Sprintf("SYNQ_API%s_URL", version))
	return ApiSetup{
		Key:     key,
		Version: version,
		Url:     url,
	}
}

func SetupSynqApi(setup ...ApiSetup) (api synq.ApiF) {
	var config ApiSetup
	if len(setup) > 0 {
		config = setup[0]
	} else {
		config = GetSetupByEnv("")
	}
	if config.Key == "" {
		log.Println("WARNING : no Synq API key specified")
	}
	if strings.Contains(config.Key, ".") || config.Version == "v2" {
		api = synq.NewV2(config.Key)
	} else {
		api = synq.NewV1(config.Key)
	}
	if config.Url != "" {
		api.SetUrl(config.Url)
	}
	return api
}
