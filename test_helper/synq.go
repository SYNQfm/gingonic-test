package test_helper

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	VIDEO_ID      = "45d4062f99454c9fb21e5186a09c2119"
	LIVE_VIDEO_ID = "ec37c42b4aab46f18003b33c66e5e641"
)

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
					ke = validVideo(video_id)
					if ke != "" {
						w.WriteHeader(http.StatusBadRequest)
						resp = []byte(ke)
					} else {
						resp, _ = ioutil.ReadFile("../sample/video.json")
					}
				case "/v1/video/create":
					resp, _ = ioutil.ReadFile("../sample/new_video.json")
				case "/v1/video/upload":
					resp, _ = ioutil.ReadFile("../sample/upload.json")
				default:
					w.WriteHeader(http.StatusBadRequest)
					resp = []byte(HTTP_NOT_FOUND)
				}
			}
		}
		w.Write(resp)
	}))
}
