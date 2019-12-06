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
	"github.com/colbyleiske/cse138_assignment2/shard"
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
		incClock, ok := r.Context().Value(ctx.ContextCausalContextKey).(shard.CausalContext)
		if !ok {
			log.Println("Could not get context from incoming request")
			return
		}

		if !ok {
			resp := struct {
				kvstore.ResponseMessage
				Exists bool `json:"doesExist"`
			}{
				kvstore.ResponseMessage{"No key", fmt.Sprintf("Error in %s", r.Method), "", addr, incClock}, false,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(key) > 50 {
			resp := kvstore.ResponseMessage{"Key is too long", fmt.Sprintf("Error in %s", r.Method), "", addr, incClock}
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
		cc := struct {
			CausalContext shard.CausalContext `json:"causal-context"`
		}{}

		if err := json.Unmarshal(bodyBytes, &cc); err != nil {
			log.Println("", err)
			return
		}

		if len(cc.CausalContext.Context) == 0 {
			cc.CausalContext.Context = make(map[string]vectorclock.VectorClock)
		}

		incContext, ok := cc.CausalContext.Context[config.Config.CurrentShardID]
		if !ok {
			log.Println("no clocks included - assuming a 0 clock")
			cc.CausalContext.Context[config.Config.CurrentShardID] = *vectorclock.NewVectorClock(config.Config.CurrentShard().Nodes, config.Config.Address)
		}
		//Given a map of shards -> Vector clcoks, check if our clock is populated... If it is not present or empty, we set it to all 0
		if _, ok := incContext.Clocks[config.Config.Address]; !ok || len(incContext.Clocks) == 0 {
			log.Println("no clocks included - assuming a 0 clock")
			incContext = *vectorclock.NewVectorClock(config.Config.CurrentShard().Nodes, config.Config.Address)
		}

		if !incContext.HappenedBefore(config.Config.CurrentShard().VectorClock) {
			//we received a request from a node that has seen more in the future than us
			log.Println("doesn't work - too much context")
			resp := kvstore.ResponseMessage{"Unable to satisfy request", fmt.Sprintf("Error in %s", r.Method), "", "", cc.CausalContext}
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(resp)
			return
		}

		ctx := context.WithValue(r.Context(), ctx.ContextCausalContextKey, cc.CausalContext)
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
