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
	log.Println("Testing isColor")
	assert := assert.New(t)
	assert.True(isColor("4e565e"))
	assert.True(isColor("ffffff"))
	assert.False(isColor("#3ac0a7"))
	assert.False(isColor("#d44747"))
	assert.False(isColor("#4e565e"))
}

func TestIsNumber(t *testing.T) {
	log.Println("Testing isNumber")
	assert := assert.New(t)
	assert.True(isNumber("1234"))
	assert.False(isNumber("1234a"))
}

func TestParseBool(t *testing.T) {
	log.Println("Testing parseBool")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "true")
	values.Set("test2", "false")
	ret := parseBool("default", values)
	assert.False(ret)
	ret = parseBool("test", values)
	assert.True(ret)
	ret = parseBool("test2", values)
	assert.False(ret)
}

func TestParseColor(t *testing.T) {
	log.Println("Testing parseColor")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "4e565e")
	ret := parseColor("default", values)
	assert.Equal("", ret)
	ret = parseColor("default", values, "000000")
	assert.Equal("#000000", ret)
	ret = parseColor("test", values)
	assert.Equal("#4e565e", ret)
}

func TestParseInt(t *testing.T) {
	log.Println("Testing parseInt")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "15")
	ret := parseInt("default", values)
	assert.Equal(0, ret)
	ret = parseInt("test", values)
	assert.Equal(15, ret)
}

func TestParseString(t *testing.T) {
	log.Println("Testing parseString")
	assert := assert.New(t)
	values := url.Values{}
	values.Set("test", "whatever")
	ret := parseString("default", values)
	assert.Equal("", ret)
	ret = parseString("test", values)
	assert.Equal("whatever", ret)
}
