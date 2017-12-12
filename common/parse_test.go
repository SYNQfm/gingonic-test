package common

import (
	"log"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDbUrl(t *testing.T) {
	assert := assert.New(t)
	name := ParseDatabaseUrl("")
	assert.Equal("", name)
	name = ParseDatabaseUrl("://abcd")
	assert.Equal("", name)
	name = ParseDatabaseUrl("postgres://user:password@host.com:5432/dbname")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password", name)
	name = ParseDatabaseUrl("postgres://user:password@host.com:5432/dbname?sslmode=disable")
	assert.Equal("host=host.com port=5432 dbname=dbname user=user password=password sslmode=disable", name)
	name = ParseDatabaseUrl(DEFAULT_DB_URL)
	assert.Equal("host=localhost port=5432 dbname=db_test user=circleci password=circleci sslmode=disable", name)
}

func TestIsColor(t *testing.T) {
	log.Println("Testing IsColor")
	assert := assert.New(t)
	assert.True(IsColor("4e565e"))
	assert.True(IsColor("ffffff"))
	assert.False(IsColor("#3ac0a7"))
	assert.False(IsColor("#d44747"))
	assert.False(IsColor("#4e565e"))
}

func TestIsNumber(t *testing.T) {
	log.Println("Testing IsNumber")
	assert := assert.New(t)
	assert.True(IsNumber("1234"))
	assert.False(IsNumber("1234a"))
}

func TestParseBool(t *testing.T) {
	log.Println("Testing ParseBool")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "true")
	values.Set("test2", "false")
	ret := ParseBool("default", values)
	assert.False(ret)
	ret = ParseBool("test", values)
	assert.True(ret)
	ret = ParseBool("test2", values)
	assert.False(ret)
}

func TestParseColor(t *testing.T) {
	log.Println("Testing ParseColor")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "4e565e")
	ret := ParseColor("default", values)
	assert.Equal("", ret)
	ret = ParseColor("default", values, "000000")
	assert.Equal("#000000", ret)
	ret = ParseColor("test", values)
	assert.Equal("#4e565e", ret)
}

func TestParseInt(t *testing.T) {
	log.Println("Testing ParseInt")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "15")
	ret := ParseInt("default", values)
	assert.Equal(0, ret)
	ret = ParseInt("test", values)
	assert.Equal(15, ret)
}

func TestParseString(t *testing.T) {
	log.Println("Testing ParseString")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "whatever")
	ret := ParseString("default", values)
	assert.Equal("", ret)
	ret = ParseString("test", values)
	assert.Equal("whatever", ret)
}
