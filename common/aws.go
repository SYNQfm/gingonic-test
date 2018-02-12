package common

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
)

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
