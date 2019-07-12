// go-rs/method-override
// Copyright(c) 2019 Roshan Gade. All rights reserved.
// MIT Licensed

package methodoverride

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-rs/rest-api-framework"
)

var (
	header  = "X-HTTP-Method-Override"
	methods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
		http.MethodPatch,
	}
)

// validate method
func isValidMethod(method string) bool {
	for _, v := range methods {
		if strings.EqualFold(v, method) {
			return true
		}
	}
	return false
}

// convert json to string
func jsonToString(val interface{}) string {
	b, _ := json.Marshal(val)
	return string(b)
}

// High-level merge of request json into query
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

// Possible error codes expose to handle errors
const (
	ErrCodeMalformedBody = "MALFORMED_BODY"
	ErrCodeFormParse     = "FORM_PARSE_ERROR"
)

// Load method override middleware/interceptor.
func Load() rest.Handler {
	return func(ctx *rest.Context) {
		method := ctx.Request.Header.Get(header)
		if method != "" && isValidMethod(method) {
			ctx.Set("OriginalMethod", ctx.Request.Method)
			if strings.EqualFold(method, http.MethodGet) && strings.EqualFold(ctx.Request.Method, http.MethodPost) {
				// parsing only json request body and form-data
				contentType := strings.ToLower(ctx.Request.Header.Get("content-type"))
				if strings.Contains(contentType, "application/json") {
					var body map[string]interface{}
					if ctx.Body != nil && reflect.TypeOf(ctx.Body).String() == "map[string]interface {}" {
						body = ctx.Body
					} else {
						err := json.NewDecoder(ctx.Request.Body).Decode(&body)
						if err != nil {
							ctx.Status(400).ThrowWithError(ErrCodeMalformedBody, err)
							return
						}
					}
					ctx.Query = mergeQuery(ctx.Query, body)
				} else if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
					if ctx.Request.Form != nil {
						ctx.Query = ctx.Request.Form
					} else {
						err := ctx.Request.ParseForm()
						if err != nil {
							ctx.Status(400).ThrowWithError(ErrCodeFormParse, err)
							return
						}
						ctx.Query = ctx.Request.Form
					}
				}
			}

			ctx.Request.Method = method
		}
	}
}
