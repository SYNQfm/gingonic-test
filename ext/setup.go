package ext

import (
	"github.com/SYNQfm/helpers/common"
	"github.com/jmoiron/sqlx"
)

const (
	DEFAULT_DB_URL = "postgres://circleci:circleci@localhost:5432/db_test?sslmode=disable"
)

func SetupDB(def_url ...string) *sqlx.DB {
	dbaddr := common.GetDB()
	return sqlx.MustConnect("postgres", dbaddr)
}
