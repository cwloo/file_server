package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cwloo/gonet/logs"
)

func setResponseHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id,C-Token,cz-sdk-key,cz-sdk-sign")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,HEAD,TRACE,OPTIONS,DELETE,PUT")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Request-Headers,Access-Control-Request-Method,Content-Type,New-Token,New-Expires-At,New-C-Token,New-C-Expires-At")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Host", r.Header.Get("Host"))
	w.Header().Set("X-Real-IP", r.Header.Get("X-Real-IP"))
	w.Header().Set("X-Forwarded-For", r.Header.Get("X-Forwarded-For"))
	w.Header().Set("X-Forwarded-Proto", r.Header.Get("X-Forwarded-Proto"))
	w.Header().Set("Remote-Host", r.Header.Get("Remote-Host"))
	w.Header().Set("User-Agent", r.Header.Get("User-Agent"))
	w.Header().Set("Referer", r.Header.Get("Referer"))
	w.Header().Set("Access-Control-Request-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Request-Method", r.Header.Get("Access-Control-Request-Method"))
	w.Header().Set("Origin", r.Header.Get("Origin"))
	w.Header().Set("Sec-Fetch-Dest", r.Header.Get("Sec-Fetch-Dest"))
	w.Header().Set("Sec-Fetch-Mode", r.Header.Get("Sec-Fetch-Mode"))
	w.Header().Set("Sec-Fetch-Site", r.Header.Get("Sec-Fetch-Site"))
	// w.Header().Set("Accept-Encoding", r.Header.Get("Accept-Encoding"))
	// w.Header().Set("Accept-Language", r.Header.Get("Accept-Language"))
}

func writeResponse(w http.ResponseWriter, r *http.Request, v any) {
	j, _ := json.Marshal(v)
	setResponseHeader(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(j)))
	_, err := w.Write(j)
	if err != nil {
		logs.Errorf(err.Error())
		return
	}
	logs.Debugf("%v", string(j))
}
