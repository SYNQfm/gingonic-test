package ext

import (
	"github.com/SYNQfm/helpers/common"
	"github.com/jmoiron/sqlx"
)

func SetupDB(def_url ...string) *sqlx.DB {
	dbaddr := common.GetDB(def_url...)
	return sqlx.MustConnect("postgres", dbaddr)
}
