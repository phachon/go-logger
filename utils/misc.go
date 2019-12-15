package utils

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Misc struct {
}

func NewMisc() *Misc {
	return &Misc{}
}

//格式化 unix 时间戳
func (misc *Misc) FormatUnixTime(unixTime int64) string {
	tm := time.Unix(unixTime, 0)
	return tm.Format("2006-01-02 15:04:05")
}

//map Intersect
func (misc *Misc) MapIntersect(defaultMap map[string]interface{}, inputMap map[string]interface{}) map[string]interface{} {
	for key, _ := range defaultMap {
		inputValue, ok := inputMap[key]
		if !ok {
			continue
		}
		defaultMap[key] = inputValue
	}
	return defaultMap
}

//http get request
func (misc *Misc) HttpGet(queryUrl string, queryValues map[string]string, headerValues map[string]string, timeout int) (body string, code int, err error) {
	if !strings.Contains(queryUrl, "?") {
		queryUrl += "?"
	}

	queryString := ""
	for queryKey, queryValue := range queryValues {
		queryString = queryString + "&" + queryKey + "=" + url.QueryEscape(queryValue)
	}
	queryString = strings.Replace(queryString, "&", "", 1)
	queryUrl += queryString

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return
	}
	if (headerValues != nil) && (len(headerValues) > 0) {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	code = resp.StatusCode
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(bodyByte), code, nil
}

//http post request
func (misc *Misc) HttpPost(queryUrl string, queryValues map[string]string, headerValues map[string]string, timeout int) (body string, code int, err error) {
	if !strings.Contains(queryUrl, "?") {
		queryUrl += "?"
	}
	queryString := ""
	for queryKey, queryValue := range queryValues {
		queryString = queryString + "&" + queryKey + "=" + url.QueryEscape(queryValue)
	}
	queryString = strings.Replace(queryString, "&", "", 1)
	queryUrl += queryString

	req, err := http.NewRequest("POST", queryUrl, strings.NewReader(queryString))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if (headerValues != nil) && (len(headerValues) > 0) {
		for key, value := range headerValues {
			req.Header.Set(key, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(bodyByte), code, nil
}

// rand string
func (m *Misc) RandString(strlen int) string {
	codes := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeLen := len(codes)
	data := make([]byte, strlen)
	rand.Seed(time.Now().UnixNano() + rand.Int63() + rand.Int63() + rand.Int63() + rand.Int63())
	for i := 0; i < strlen; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	return string(data)
}
