/*!
 * rest-api-framework
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package methodoverride

import (
	"github.com/go-rs/rest-api-framework"
	"strings"
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
		return v == method
	}
	return false
}

/**
 * Method Override
 */
func MethodOverride() rest.Handler {
	return func(ctx *rest.Context) {
		method := strings.ToUpper(ctx.Request.Header.Get(header))
		if method != "" && isValidMethod(method) {
			ctx.Set("original_method", ctx.Request.Method)
			ctx.Request.Method = method
		}
	}
}
