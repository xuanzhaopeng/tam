package main

import (
	"net/http"
	"encoding/json"
	"tam/account"
	"fmt"
	"io/ioutil"
)

var paths = struct {
	Fetch, Release string
}{
	Fetch:      "/fetch",
	Release:    "/release",
}

func reply(w http.ResponseWriter, msg map[string]interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func httpMethodOnly(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			reply(w, errMsg(fmt.Sprintf("%s method not allowed", r.Method)), http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func errMsg(msg string) map[string]interface{} {
	return map[string]interface{}{"message": msg}
}

func fetchHandler(writer http.ResponseWriter, req *http.Request) {
	var filter account.Filter
	if req.Body == nil {
		http.Error(writer, "Please send a request body", http.StatusBadRequest)
		return
	}

	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		reply(writer, errMsg(err.Error()), http.StatusNotFound)
		return
	}

	if len(bodyData) > 0 {
		err := json.Unmarshal(bodyData, &filter)
		if err != nil {
			reply(writer, errMsg(err.Error()), http.StatusNotFound)
			return
		}
	} else {
		filter = account.BuildEmptyFilter()
	}

	data, err := accounts.Fetch(filter)
	if err != nil {
		reply(writer, errMsg(err.Error()), http.StatusNotFound)
		return
	}
	reply(writer, data, http.StatusOK)
}

func releaseHandler(writer http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get(accountKey)
	if key == "" {
		reply(writer, errMsg(fmt.Sprintf("the url parameter '%s' is required", accountKey)), http.StatusNotFound)
		return
	}
	err := accounts.Release(key)
	if err != nil {
		reply(writer, errMsg(err.Error()), http.StatusNotFound)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(paths.Fetch, httpMethodOnly(http.MethodPost, fetchHandler))
	mux.HandleFunc(paths.Release, httpMethodOnly(http.MethodGet, releaseHandler))
	return mux
}
