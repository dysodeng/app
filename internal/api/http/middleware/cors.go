package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CrossDomain cors middleware
func CrossDomain(ctx *gin.Context) {
	method := ctx.Request.Method

	origin := ctx.Request.Header.Get("Origin")

	var headerKeys []string
	for k := range ctx.Request.Header {
		headerKeys = append(headerKeys, k)
	}

	headerString := strings.Join(headerKeys, ",")
	if headerString == "" {
		headerString = "access-control-allow-origin, access-control-allow-headers"
	} else {
		headerString = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerString)
	}

	if origin != "" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", headerString)
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-Control-Allow-Headers", "Authorization,Content-Length,X-CSRF-Token,Token,session,X_Requested_With,Accept,Origin,Host,Connection,Accept-Encoding,Accept-Language,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Pragma,Platform-Type,X-Trace-Id,X-Span-Id")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,Platform-Type,X-Trace-Id,X-Span-Id")
		ctx.Header("Access-Control-Allow-Credentials", "false")
		ctx.Set("Content-Type", "application/json")
	}

	if method == "OPTIONS" {
		ctx.JSON(http.StatusOK, "")
	}

	ctx.Next()
}
