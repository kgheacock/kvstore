package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/config"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s", r.RequestURI, config.Config.ForwardAddress)
		next.ServeHTTP(w, r)
	})
}

func forwardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.Config.IsFollower {
			next.ServeHTTP(w, r)
			return
		}

		//Rather compute this once beforehand than in every case.
		errResp := &ErrorForwardResponse{Error: "Main instance is down", Message: fmt.Sprintf("Error in %s", r.Method)}
		jsonErrResp, err := json.Marshal(errResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(jsonErrResp)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		url := r.URL
		url.Scheme = "http"
		url.Host = config.Config.ForwardAddress

		proxyReq, err := http.NewRequest(r.Method, url.String(), bytes.NewReader(body))
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(jsonErrResp)
			return
		}

		proxyReq.Header = r.Header
		proxyReq.Header.Set("Host", r.Host)
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

		client := &http.Client{}
		proxyResp, err := client.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(jsonErrResp)
			return
		}
		defer proxyResp.Body.Close()

		proxyBody, err := ioutil.ReadAll(proxyResp.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(jsonErrResp)
			return
		}

		w.WriteHeader(proxyResp.StatusCode)
		w.Write(proxyBody)
	})
}
