package httpUtils

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	_http "github.com/deeptest-com/deeptest-next/pkg/libs/http"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Verbose = true
)

func Get(url string, headers map[string]string) (ret []byte, err error) {
	return gets(url, "GET", headers)
}

func Post(url string, data interface{}, headers map[string]string) (ret []byte, err error) {
	return posts(url, "POST", data, headers)
}
func Put(url string, data interface{}, headers map[string]string) (ret []byte, err error) {
	return posts(url, "PUT", data, headers)
}
func Delete(url string, headers map[string]string) (ret []byte, err error) {
	return gets(url, "DELETE", headers)
}

func gets(url, method string, headers map[string]string) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	if !_http.IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("read response failed, StatusCode: %d.", resp.StatusCode))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := ioutil.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}
func posts(url string, method string, data interface{}, headers map[string]string) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	dataBytes, err := json.Marshal(data)
	if Verbose {
		_logUtils.Infof("===DEBUG===     data: %s", string(dataBytes))
	}

	if err != nil {
		_logUtils.Infof(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)

	req, err := http.NewRequest(method, url, strings.NewReader(dataStr))
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	//req.Header.SetVariable("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !_http.IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := ioutil.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}
