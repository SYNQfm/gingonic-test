package common

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("errored", ParseType("Error"))
	assert.Equal("count", ParseType("ct"))
	assert.Equal("skipped", ParseType("SKiP"))
	assert.Equal("already", ParseType("already"))
}

func TestRet(t *testing.T) {
	assert := assert.New(t)
	ret := NewRet("test")
	assert.Equal(0, ret.Value("count"))
	ret.Add("count")
	assert.True(ret.Eq("count", 1))
	assert.True(ret.Gte("count", 1))
	assert.True(ret.Lte("count", 1))
	assert.False(ret.Lt("count", 1))
	assert.False(ret.Gt("count", 1))
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	ret := NewRet("test")
	str := ret.String()
	assert.Contains(str, "for test")
	ret.Add("count")
	str = ret.String()
	assert.Contains(str, "processed 1")
	ret.AddBytes(1000000)
	// reset the start date
	ret.Start = ret.Start.Add(-1 * time.Hour)
	str = ret.String()
	log.Println(str)
	assert.Contains(str, "1 megs (speed 0.1 mbps)")
	assert.Contains(str, "took 60 mins")
	ret.Start = ret.Start.Add(59 * time.Minute)
	str = ret.String()
	log.Println(str)
	assert.Contains(str, "1 megs (speed 0.1 mbps)")
	assert.Contains(str, "took 60 sec")
}

func TestConvert(t *testing.T) {
	assert := assert.New(t)
	uuid := "45d4062f99454c9fb21e5186a09c2119"
	vid := ConvertToUUIDFormat(uuid)
	assert.Equal("45d4062f-9945-4c9f-b21e-5186a09c2119", vid)
	vid2 := ConvertToUUIDFormat(vid)
	assert.Equal("45d4062f-9945-4c9f-b21e-5186a09c2119", vid2)
}

func TestGetFileExtension(t *testing.T) {
	log.Println("Testing getFileExtension")
	assert := assert.New(t)
	ext := GetFileExtension("video/mp4")
	assert.Equal(extension, ext)
}

func TestGetAwsSignature(t *testing.T) {
	log.Println("Testing getAwsSignature")
	assert := assert.New(t)
	assert.NotNil(getAwsSignature("message", "secret"))
}

func TestGetMultipartSignature(t *testing.T) {
	log.Println("Testing getMultipartSignature")
	assert := assert.New(t)
	headers := "POST\n\nvideo/mp4\n\nx-amz-acl:public-read\nx-amz-date:Mon, 23 Oct 2017 18:50:29 GMT\n" + videoKey + "?uploads"
	signature := getMultipartSignature(headers, DefaultAWS.SecretKey)
	assert.NotEmpty(signature)

	mpsignature := struct {
		Signature string `json:"signature"`
	}{}
	err := json.Unmarshal(signature, &mpsignature)
	assert.Nil(err)
	assert.Equal("f+w3SFD1Pk5Cur06MAZRJ250zHg=", mpsignature.Signature)
}
