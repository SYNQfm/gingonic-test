package ext

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"
)

func TestUpdateJson(t *testing.T) {
	assert := require.New(t)
	bytes := []byte(`{"wow":"yay"}`)
	vals := make(map[string]interface{})
	strVal := "cool"
	intVal := 123
	int64Val := int64(456)
	vals["str"] = "cool"
	vals["int"] = intVal
	vals["int64"] = int64Val
	b := UpdateJson(bytes, vals)
	v, _ := jsonparser.GetString(b, "str")
	assert.Equal(strVal, v)
	v, _ = jsonparser.GetString(b, "wow")
	assert.Equal("yay", v)
	i, _ := jsonparser.GetInt(b, "int")
	assert.Equal(int64(intVal), i)
	i64, _ := jsonparser.GetInt(b, "int64")
	assert.Equal(int64Val, i64)
}
