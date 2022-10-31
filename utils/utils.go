package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var randSource = rand.New(rand.NewSource(time.Now().Unix()))
var etcpasswdReg = regexp.MustCompile("((root|bin|daemon|sys|sync|games|man|mail|news|www-data|uucp|backup|list|proxy|gnats|nobody|syslog|mysql|bind|ftp|sshd|postfix):[\\d\\w\\-\\s,]+:\\d+:\\d+:[\\w\\-_\\s,]*:[\\w\\-_\\s,\\/]*:[\\w\\-_,\\/]*[\\r\\n])")

func RandomLowercase(n int) string {
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	return randomStr(randSource, lowercase, n)
}

func RandomUppercase(n int) string {
	lowercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return randomStr(randSource, lowercase, n)
}

func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(data)
	return string(decode), err
}

func UrlEncode(data string) string {
	return url.QueryEscape(data)
}

func UrlDecode(data string) (string, error) {
	return url.QueryUnescape(data)
}

func randomStr(randSource *rand.Rand, letterBytes string, n int) string {
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
		//letterBytes   = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	)
	randBytes := make([]byte, n)
	for i, cache, remain := n-1, randSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			randBytes[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(randBytes)
}

//func BytesContains(b[]byte,subslice[]byte) bool {
//	return bytes.Contains(b,subslice)
//}

func MD5(data interface{}) (result string) {
	switch a := data.(type) {
	case []byte:
		result = fmt.Sprintf("%x", md5.Sum(a))
	case string:
		result = fmt.Sprintf("%x", md5.Sum([]byte(a)))
	case int:
		result = fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(a))))
	}
	return
}

func UrlPath(u *url.URL) string {
	p := u.Path
	if len(p) == 0 {
		p = "/"
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, strings.TrimRight(path.Dir(p), "/"))
}

func UrlString(u *url.URL) string {
	return u.String()
}

func UrlRoot(u *url.URL) string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

func MatchEtcPasswd(content string) bool {
	return etcpasswdReg.MatchString(content)
}
