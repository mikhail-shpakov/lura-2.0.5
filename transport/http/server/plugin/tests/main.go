// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
)

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer("krakend-server-example")

type registerer string

var logger Logger = nil

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (registerer) registerHandlers(_ context.Context, _ map[string]interface{}, _ http.Handler) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	// possible config example:
	/*
	   "extra_config":{
	       "plugin/http-server":{
	           "name":["krakend-server-example"],
	           "krakend-server-example":{
	               "A":"foo",
	               "B":42
	           }
	       }
	   }
	*/

	if logger == nil {
		// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
		}), nil
	}

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
		logger.Debug("request:", html.EscapeString(req.URL.Path))
	}), nil
}

func main() {}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}
