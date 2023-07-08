package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// 字符串签名
func StringSign(str, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(str))
	sign := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString(sign)
}
