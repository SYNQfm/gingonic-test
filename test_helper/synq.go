package test_helper

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var sample_dir = "sample/"

var testReqs []*http.Request
var testServer *httptest.Server

const (
	VIDEO_ID          = "45d4063d00454c9fb21e5186a09c3115"
	VIDEO_ID2         = "55d4062f99454c9fb21e5186a09c2115"
	LIVE_VIDEO_ID     = "ec37c42b4aab46f18003b33c66e5e641"
	API_KEY           = "aba179c14ab349e0bb0d12b7eca5fa24"
	API_KEY2          = "cba179c14ab349e0bb0d12b7eca5fa25"
	UPLOAD_KEY        = "projects/0a/bf/0abfe1b849154082993f2fce77a16fd9/uploads/videos/55/d4/55d4062f99454c9fb21e5186a09c2115.mp4"
	INVALID_UUID      = `{"url": "http://docs.synq.fm/api/v1/errors/invalid_uuid","name": "invalid_uuid","message": "Invalid uuid. Example: '1c0e3ea4529011e6991554a050defa20'."}`
	VIDEO_NOT_FOUND   = `{"url": "http://docs.synq.fm/api/v1/errors/not_found_video","name": "not_found_video","message": "Video not found."}`
	API_KEY_NOT_FOUND = `{"url": "http://docs.synq.fm/api/v1/errors/not_found_api_key","name": "not_found_api_key","message": "API key not found."}`
	HTTP_NOT_FOUND    = `{"url": "http://docs.synq.fm/api/v1/errors/http_not_found","name": "http_not_found","message": "Not found."}`
)

func SetupServer() string {
	if testServer != nil {
		testServer.Close()
	}
	testReqs = testReqs[:0]
	testServer = SynqStub()
	return testServer.URL
}

func validKey(key string) string {
	if len(key) != 32 {
		return INVALID_UUID
	} else if key != API_KEY {
		return API_KEY_NOT_FOUND
	}
	return ""
}

func validVideo(id string) string {
	if len(id) != 32 {
		return INVALID_UUID
	} else if id != VIDEO_ID || id != LIVE_VIDEO_ID {
		return VIDEO_NOT_FOUND
	}
	return ""
}

func SynqStub() *httptest.Server {
	var resp []byte
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("here in synq response", r.RequestURI)
		testReqs = append(testReqs, r)
		if r.Method == "POST" {
			bytes, _ := ioutil.ReadAll(r.Body)
			//Parse response body
			v, _ := url.ParseQuery(string(bytes))
			key := v.Get("api_key")
			ke := validKey(key)
			if ke != "" {
				w.WriteHeader(http.StatusBadRequest)
				resp = []byte(ke)
			} else {
				switch r.RequestURI {
				case "/v1/video/details":
					video_id := v.Get("video_id")
					log.Println("in video details with video " + video_id)
					ke = validVideo(video_id)
					if ke != "" {
						w.WriteHeader(http.StatusBadRequest)
						resp = []byte(ke)
					} else {
						resp, _ = ioutil.ReadFile(sample_dir + "video.json")
					}
				case "/v1/video/create":
					resp, _ = ioutil.ReadFile(sample_dir + "new_video.json")
				case "/v1/video/upload":
					resp, _ = ioutil.ReadFile(sample_dir + "upload.json")
				default:
					w.WriteHeader(http.StatusBadRequest)
					resp = []byte(HTTP_NOT_FOUND)
				}
			}
		}
		w.Write(resp)
	}))
}
