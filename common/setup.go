package common

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/SYNQfm/SYNQ-Golang/synq"
	"github.com/jmoiron/sqlx"
)

func SetupDB(def_url ...string) *sqlx.DB {
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" && len(def_url) > 0 {
		db_url = def_url[0]
	}
	dbaddr := ParseDatabaseUrl(db_url)
	return sqlx.MustConnect("postgres", dbaddr)
}

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

// this parses the database url and returns it in the format sqlx.DB expects
func ParseDatabaseUrl(dbUrl string) string {
	if dbUrl == "" {
		return ""
	}
	u, e := url.Parse(dbUrl)
	if e != nil {
		log.Printf("Error parsing '%s' : %s\n", dbUrl, e.Error())
		return ""
	}
	str := fmt.Sprintf("host=%s port=%s dbname=%s",
		u.Hostname(), u.Port(), strings.Replace(u.Path, "/", "", -1))
	if u.User != nil && u.User.Username() != "" {
		pass, set := u.User.Password()
		str = str + " user=" + u.User.Username()
		if set {
			str = str + " password=" + pass
		}
	}
	ssl := u.Query().Get("sslmode")
	if ssl != "" {
		str = str + " sslmode=" + ssl
	}
	return str
}