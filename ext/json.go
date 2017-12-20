package common

import (
	"strconv"

	"github.com/buger/jsonparser"
)

func UpdateJson(bytes []byte, vals map[string]interface{}) []byte {
	for k, v := range vals {
		switch val := v.(type) {
		case int:
			str := strconv.Itoa(val)
			bytes, _ = jsonparser.Set(bytes, []byte(str), k)
		case int64:
			str := strconv.FormatInt(val, 10)
			bytes, _ = jsonparser.Set(bytes, []byte(str), k)
		case string:
			bytes, _ = jsonparser.Set(bytes, []byte(`"`+val+`"`), k)
		}
	}
	return bytes
}
