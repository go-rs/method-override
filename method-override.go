/*!
 * go-rs/method-override
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package methodoverride

import (
	"strings"

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
 * Method Override
 */
func Load() rest.Handler {
	return func(ctx *rest.Context) {
		method := strings.ToUpper(ctx.Request.Header.Get(header))
		if method != "" && isValidMethod(method) {
			ctx.Set("OriginalMethod", ctx.Request.Method)
			ctx.Request.Method = method
		}
	}
}
