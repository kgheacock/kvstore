package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/ctx"
	"github.com/colbyleiske/cse138_assignment2/kvstore"
	"github.com/colbyleiske/cse138_assignment2/vectorclock"
	"github.com/gorilla/mux"
)

var (
	ErrIPNotFound error = errors.New("ip not found")
)

func (s *Store) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s", config.Config.Address, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *Store) validateParametersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, ok := vars["key"]
		var addr string
		if len(r.Header.Get("X-Real-Ip")) != 0 {
			addr = config.Config.Address
		}

		if !ok {
			resp := struct {
				kvstore.ResponseMessage
				Exists bool `json:"doesExist"`
			}{
				kvstore.ResponseMessage{"No key", fmt.Sprintf("Error in %s", r.Method), "", addr, config.Config.CurrentShard().VectorClock}, false,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(key) > 50 {
			resp := kvstore.ResponseMessage{"Key is too long", fmt.Sprintf("Error in %s", r.Method), "", addr, config.Config.CurrentShard().VectorClock}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Store) bufferRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source, ok := r.Context().Value(ctx.ContextSourceKey).(string)
		if ok && source == ctx.INTERNAL {
			next.ServeHTTP(w, r)
			return
		}
		if s.kvstore.State() == kvstore.NORMAL {
			next.ServeHTTP(w, r)
			return
		}

		w.Write([]byte("Resharding in progress"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	})
}

/*
To get the value from the context
source, ok := r.Context().Value(ctx.ContextSourceKey).(string)
ctx.INTERNAL or ctx.EXTERNAL
*/
func (s *Store) checkSourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source := ctx.INTERNAL
		if !config.IsIPInternal(r.Header.Get("X-Real-Ip")) || len(r.Header.Get("X-Forwarded-For")) != 0 {
			source = ctx.EXTERNAL
		}

		ctx := context.WithValue(r.Context(), ctx.ContextSourceKey, source)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Store) checkVectorClock(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		causalContext := struct {
			Context vectorclock.VectorClock `json:"causal-context"`
		}{}

		if err := json.Unmarshal(bodyBytes, &causalContext); err != nil {
			log.Println("some error here about failing to get json read")
			return
		}

		ctx := context.WithValue(r.Context(), ctx.ContextCausalContextKey, causalContext.Context)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Store) forwardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		keyIPLocation, err := s.hasher.DAL().GetServerByKey(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if keyIPLocation == config.Config.Address {
			next.ServeHTTP(w, r)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		url := r.URL
		url.Scheme = "http"
		url.Host = keyIPLocation

		proxyReq, err := http.NewRequest(r.Method, url.String(), bytes.NewReader(body))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		proxyReq.Header = r.Header
		proxyReq.Header.Set("X-Real-Ip", config.Config.Address)
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

		client := &http.Client{}
		proxyResp, err := client.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer proxyResp.Body.Close()

		proxyBody, err := ioutil.ReadAll(proxyResp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(proxyResp.StatusCode)
		w.Write(proxyBody)
	})
}
