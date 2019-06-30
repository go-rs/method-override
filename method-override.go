/*!
 * go-rs/method-override
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package methodoverride

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-rs/ordered-json"
	"github.com/go-rs/rest-api-framework"
)

var (
	header  = "X-HTTP-Method-Override"
	methods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"}
)

/**
 * Validate method
 */
func isValidMethod(method string) bool {
	for _, v := range methods {
		if v == method {
			return true
		}
	}
	return false
}

/**
 * Convert json to string
 */
func jsonToString(val interface{}) string {
	b, _ := json.Marshal(val)
	return string(b)
}

/**
 * High-level merge of request json into query
 */
func mergeQuery(query url.Values, body map[string]interface{}) url.Values {
	for key, val := range body {
		str, ok := val.(string)
		if ok {
			query.Add(key, str)
		} else {
			if reflect.TypeOf(val).String() == "[]interface {}" {
				for _, val := range val.([]interface{}) {
					str, ok = val.(string)
					if ok {
						query.Add(key, str)
					} else {
						query.Add(key, jsonToString(val))
					}
				}
			} else {
				query.Add(key, jsonToString(val))
			}
		}
	}
	return query
}

/**
 * Method Override
 */
func Load() rest.Handler {
	return func(ctx *rest.Context) {
		method := strings.ToUpper(ctx.Request.Header.Get(header))
		if method != "" && isValidMethod(method) {
			ctx.Set("OriginalMethod", ctx.Request.Method)
			if ctx.Request.Method == "POST" && method == "GET" {

				contentType := strings.ToLower(ctx.Request.Header.Get("content-type"))

				if contentType == "application/json" {
					if ctx.Body != nil && reflect.TypeOf(ctx.Body).String() == "*orderedjson.OrderedMap" {
						body := ctx.Body.(*orderedjson.OrderedMap)
						ctx.Query = mergeQuery(ctx.Query, body.GetMap())
					} else {
						var body map[string]interface{}
						err := json.NewDecoder(ctx.Request.Body).Decode(&body)
						if err != nil {
							ctx.Status(400).Throw(errors.New("MALFORMED_BODY"))
							return
						}
						ctx.Query = mergeQuery(ctx.Query, body)
					}
				} else if len(ctx.Request.Form) > 0 {
					ctx.Query = ctx.Request.Form
				}
			}
			ctx.Request.Method = method
		}
	}
}
