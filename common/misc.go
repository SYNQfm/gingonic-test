package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

func ParseType(type_ string) string {
	t := strings.ToLower(type_)
	switch t {
	case "error":
		t = "errored"
	case "skip":
		t = "skipped"
	case "ct":
		t = "count"
	}
	return t
}

func Label(dur time.Duration) string {
	if dur == time.Hour {
		return "hrs"
	} else if dur == time.Minute {
		return "mins"
	} else if dur == time.Second {
		return "sec"
	} else if dur == time.Millisecond {
		return "ms"
	} else {
		return "ns"
	}
}

func DurFromLabel(label string) time.Duration {
	if label == "hrs" {
		return time.Hour
	} else if label == "mins" {
		return time.Minute
	} else if label == "sec" {
		return time.Second
	} else if label == "ms" {
		return time.Millisecond
	} else {
		return time.Nanosecond
	}

}

func DurVal(dur time.Duration) (int, string) {
	var t time.Duration
	if dur >= 1000*time.Minute {
		t = time.Hour
	} else if dur >= 1000*time.Second {
		t = time.Minute
	} else if dur >= 10000*time.Millisecond {
		t = time.Second
	} else if dur >= 10000*time.Nanosecond {
		t = time.Millisecond
	} else {
		t = time.Nanosecond
	}
	taken := int(dur / t)
	return taken, Label(t)
}

func BytesVal(bytes int64) (int64, string) {
	if bytes == 0 {
		return 0, ""
	}
	meg := int64(1000 * 1000)
	if bytes < 50000 {
		return bytes, "B"
	} else if bytes < meg {
		return bytes / 1000, "KB"
	}
	b := bytes / meg
	label := "MB"
	if b > 500000 {
		b = b / (1000 * 1000)
		label = "TB"
	} else if b > 5000 {
		b = b / 1000
		label = "GB"
	}
	return b, label
}

// return 32 bytes into 36 bytes
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func ConvertToUUIDFormat(uuid string) string {
	if len(uuid) == 36 && strings.Count(uuid, "-") == 4 {
		return uuid
	}
	if len(uuid) != 32 {
		log.Printf("invalid uuid '%s'\n", uuid)
		return uuid
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])
}

func ConvertFromUUIDFormat(uuid string) string {
	return strings.Replace(uuid, "-", "", -1)
}

func ValidUUID(uuid string) bool {
	if uuid == "" {
		return false
	}
	if len(uuid) == 36 {
		sec := strings.Split(uuid, "-")
		if len(sec) != 5 {
			return false
		}
		return len(sec[0]) == 8 && len(sec[1]) == 4 && len(sec[2]) == 4 && len(sec[3]) == 4 && len(sec[4]) == 12
	} else if len(uuid) == 32 {
		if strings.Contains(uuid, "-") {
			return false
		}
		return true
	} else {
		return false
	}
}

// Get environment variable
func GetOsEnv(env string, mandatory bool, defaultValue ...string) (value string) {
	value = os.Getenv(env)
	if value == "" {
		if mandatory {
			log.Panicf(`Can't run without environment variable ${%s} set.`, env)
		} else if len(defaultValue) > 0 {
			value = defaultValue[0]
		}
	}
	return
}

func GetAwsSignature(message, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Multipart Upload
func GetMultipartSignature(headers, awsSecret string) []byte {
	infoMap := map[string]string{
		"signature": GetAwsSignature(headers, awsSecret),
	}

	signature, _ := json.Marshal(infoMap)
	return signature
}

func GetFileExtension(ctype string) string {
	ext := CtypeToExt(ctype)
	if ext == "" {
		return ext
	}
	return strings.Split(ext, ".")[1]
}

func CtypeToExt(ctype string) string {
	ctype = CleanCtype(ctype)
	exts, err := mime.ExtensionsByType(ctype)
	if err != nil {
		return ""
	}
	if len(exts) > 0 {
		return exts[0]
	}
	switch ctype {
	case "application/ttml+xml":
		return ".ttml"
	case "application/x-subrip":
		return ".srt"
	case "application/xml":
		return ".xml"
	case "video/mp4":
		return ".mp4"
	}
	return ""
}

func CleanCtype(ctype string) string {
	return strings.Split(ctype, ";")[0]
}

func ExtToCtype(ext string) string {
	ctype := mime.TypeByExtension(ext)
	if ctype != "" && !strings.Contains(ctype, "text/plain") {
		return CleanCtype(ctype)
	}
	switch ext {
	case ".ttml":
		return "application/ttml+xml"
	case ".srt":
		return "application/x-subrip"
	case ".mp4":
		return "video/mp4"
	case ".xml":
		return "application/xml"
	}
	return ""
}

func FindString(list []string, find string) int {
	for idx, item := range list {
		if item == find {
			return idx
		}
	}
	return -1
}

func EmptyJson(val json.RawMessage) bool {
	if val == nil || len(val) == 0 {
		return true
	}
	if string(val) == "null" {
		return true
	}
	return false
}

// Check if v2 token is expired or not
func ValidV2Token(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 { //header, payload, signature
		log.Println("Invalid JWT structure")
		return false
	}

	// get payload value, which contains expiry information
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println("Error decoding token: ", err.Error())
		return false
	}

	data := struct {
		Expiry int64 `json:"exp"`
	}{}
	json.Unmarshal(payload, &data)

	currentTime := time.Now().Unix()
	if data.Expiry < currentTime {
		return false
	}

	return true
}

func GetDir(u string) string {
	uri, _ := url.Parse(u)
	dir := filepath.Dir(uri.Path)
	return dir
}

func GetTypeByExt(ext string) string {
	switch ext {
	case "mpd":
		return "dash"
	case "ism":
		return "smooth"
	default:
		return "hls"
	}
}

func GetExtByType(assetType string) string {
	switch assetType {
	case "dash", "trailer-dash":
		return "mpd"
	case "smooth":
		return "ism"
	default:
		return "m3u8"
	}
}

func GenerateMD5(file io.Reader, fileSize, maxSize int64) (hash string, hashSize int64) {
	var err error
	h := md5.New()
	hashSize = fileSize
	if fileSize > maxSize {
		_, err = io.CopyN(h, file, maxSize)
		hashSize = maxSize
	} else {
		_, err = io.Copy(h, file)
	}
	if err != nil && err != io.EOF {
		log.Println(err.Error())
		return hash, hashSize
	}
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return hash, hashSize
}

// Parses the metadata and returns the value of the specified field
func GetMetadataStringVal(metadata json.RawMessage, field string) string {
	valueBytes, _, _, _ := jsonparser.Get(metadata, field)
	return string(valueBytes)
}
