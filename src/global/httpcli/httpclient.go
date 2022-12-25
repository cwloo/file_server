package httpcli

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

func HttpClient() (client *http.Client) {
	transport := &http.Transport{
		DisableKeepAlives:     false,
		TLSHandshakeTimeout:   time.Duration(3600) * time.Second,
		IdleConnTimeout:       time.Duration(3600) * time.Second,
		ResponseHeaderTimeout: time.Duration(3600) * time.Second,
		ExpectContinueTimeout: time.Duration(3600) * time.Second,
	}
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar:       jar,
		Timeout:   time.Duration(3600) * time.Second,
		Transport: transport,
	}
	return
}
