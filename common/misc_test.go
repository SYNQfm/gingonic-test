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

func TestDuration(t *testing.T) {
	assert := assert.New(t)
	ret := NewRet("test")
	ret.Add("videos")
	ret.Add("videos")
	ret.AddBytesFor("videos", 1000000)
	ret.AddDurFor("videos", 1*time.Second)
	str := ret.String()
	log.Println(str)
	assert.Contains(str, "videos 2 (1 MB, avg 500 KB, duration 1000 ms, avg 500 ms)")
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
	assert.Contains(str, "1 MB (speed 0.00 MBps)")
	assert.Contains(str, "took 60 mins")
	ret.Start = ret.Start.Add(59 * time.Minute)
	str = ret.String()
	log.Println(str)
	assert.Contains(str, "1 MB (speed 0.13 MBps)")
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

func TestGetAwsSignature(t *testing.T) {
	log.Println("Testing getAwsSignature")
	assert := assert.New(t)
	assert.NotNil(GetAwsSignature("message", "secret"))
}

func TestGetMultipartSignature(t *testing.T) {
	log.Println("Testing getMultipartSignature")
	assert := assert.New(t)
	videoKey := "/path/to/file"
	headers := "POST\n\nvideo/mp4\n\nx-amz-acl:public-read\nx-amz-date:Mon, 23 Oct 2017 18:50:29 GMT\n" + videoKey + "?uploads"
	signature := GetMultipartSignature(headers, "abcd")
	assert.NotEmpty(signature)

	mpsignature := struct {
		Signature string `json:"signature"`
	}{}
	err := json.Unmarshal(signature, &mpsignature)
	assert.Nil(err)
	assert.Equal("TXUvxqMH7sUU/yLcOLrlh7C5su0=", mpsignature.Signature)
}

func TestValidUUID(t *testing.T) {
	assert := assert.New(t)
	assert.False(ValidUUID(""))
	assert.False(ValidUUID("45d4063d00454c9fb21e5186a09c311"))
	assert.True(ValidUUID("45d4063d00454c9fb21e5186a09c3115"))
	assert.True(ValidUUID("9e9dc8c8-f705-41db-88da-b3034894deb9"))
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	list := []string{"a", "b", "c"}
	assert.Equal(0, FindString(list, "a"))
	assert.Equal(1, FindString(list, "b"))
	assert.Equal(2, FindString(list, "c"))
	assert.Equal(-1, FindString(list, "d"))
}

func TestExtToCtype(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("video/mp4", ExtToCtype(".mp4"))
	assert.Equal("application/xml", ExtToCtype(".xml"))
	assert.Equal("application/ttml+xml", ExtToCtype(".ttml"))
	assert.Equal("application/x-subrip", ExtToCtype(".srt"))
}

func TestCtypeToExt(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(".mp4", CtypeToExt("video/mp4"))
	assert.Equal(".xml", CtypeToExt("application/xml"))
	assert.Equal(".ttml", CtypeToExt("application/ttml+xml"))
	assert.Equal(".srt", CtypeToExt("application/x-subrip"))
}

func TestGetFileExtension(t *testing.T) {
	log.Println("Testing GetFileExtension")
	assert := assert.New(t)
	assert.Equal("mp4", GetFileExtension("video/mp4"))
	assert.Equal("xml", GetFileExtension("application/xml"))
	assert.Equal("ttml", GetFileExtension("application/ttml+xml"))
	assert.Equal("srt", GetFileExtension("application/x-subrip"))
}

func TestEmptyJson(t *testing.T) {
	assert := assert.New(t)
	assert.True(EmptyJson(nil))
	assert.True(EmptyJson([]byte{}))
	assert.True(EmptyJson([]byte(`null`)))
	assert.False(EmptyJson([]byte(`{}`)))
}

func TestGetDir(t *testing.T) {
	assert := assert.New(t)
	srcUrl := "s3://synq-frankfurt/videos/9e/9d/9e9dc8c8-f705-41db-88da-b3034894deb9/hls/master_manifest.m3u8"
	dir := GetDir(srcUrl)
	assert.Equal("/videos/9e/9d/9e9dc8c8-f705-41db-88da-b3034894deb9/hls", dir)
}

func TestGetTypeByExt(t *testing.T) {
	assert := assert.New(t)
	ftype := GetTypeByExt("mpd")
	assert.Equal("dash", ftype)
	ftype = GetTypeByExt("ism")
	assert.Equal("smooth", ftype)
	ftype = GetTypeByExt("")
	assert.Equal("hls", ftype)
}

func TestGetExtByType(t *testing.T) {
	assert := assert.New(t)
	ext := GetExtByType("dash")
	assert.Equal("mpd", ext)
	ext = GetExtByType("smooth")
	assert.Equal("ism", ext)
	ext = GetExtByType("hls")
	assert.Equal("m3u8", ext)
}
