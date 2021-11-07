package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetStringOrFile return the content of the file if it is a file otherwise return the string
func GetStringOrFile(arg string) string {
	var file []byte
	var err error
	if file, err = ioutil.ReadFile(arg); os.IsNotExist(err) {
		return arg
	}
	// Remove last \n or \r
	if file[len(file)-1] == byte(10) || file[len(file)-1] == byte(13) {
		file = file[:len(file)-1]
	}
	return string(file)
}

// RandomString return a string of length n
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// StructToMap return a url.Values from a struct
func StructToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		values.Set(typ.Field(i).Tag.Get("form"), fmt.Sprint(iVal.Field(i)))
	}
	return
}

// NewUUID generate an UUID
func NewUUID() (string, error) {
	var uuid = make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	uuidString := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return uuidString, nil

}

// GetUserAgent return an agent among popular user agent
func GetUserAgent() string {
	var userAgents = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0"}
	return userAgents[rand.Intn(len(userAgents))]
}

// Credits: https://stackoverflow.com/a/46202939/7245054
// ReSubMatchMap will applied a regex with named capture and return a map
func ReSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	matches := r.FindAllStringSubmatch(str, -1)

	subMatchMap := make(map[string]string)
	for _, match := range matches {
		for i, name := range r.SubexpNames() {
			if i != 0 && len(match) >= i {
				subMatchMap[name] = match[i]
			}
		}
	}

	return subMatchMap
}

// GetBodyInWebsite return the body of the website
func GetBodyInWebsite(url string, proxy func(*http.Request) (*url.URL, error), headers map[string]string) (string, int, error) {
	// Get random user agent
	userAgent := GetUserAgent()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", userAgent)
	// Add the headers to the request
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}

	client := &http.Client{

		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           proxy,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", -1, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp.StatusCode, nil
}
