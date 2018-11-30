package utils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
)

// DoPost issues a POST to the specified URL.
func DoPost(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Annotatef(err, "new get request %s", url)
	}
	req.Header.Set("Content-Type", "application/json")
	return request(req)
}

// DoGet issues a GET to the specified URL.
func DoGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Annotatef(err, "new get request %s", url)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	return request(req)
}

func request(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Annotate(err, "request error")
	}
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Annotate(err, "fail to read body")
	}

	if resp.StatusCode/200 != 1 {
		return nil, errors.Errorf("request status code %d data %s", resp.StatusCode, string(bodyByte))
	}

	return bodyByte, nil
}
