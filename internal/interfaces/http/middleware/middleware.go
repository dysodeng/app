package middleware

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/metrics"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// 记录请求日志
		statusCode := c.Writer.Status()

		// 状态码颜色
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = "\033[97;42m" // 绿色
		case statusCode >= 300 && statusCode < 400:
			statusColor = "\033[90;47m" // 白色
		case statusCode >= 400 && statusCode < 500:
			statusColor = "\033[97;43m" // 黄色
		default:
			statusColor = "\033[97;41m" // 红色
		}

		// 方法颜色
		var methodColor string
		switch method {
		case "GET":
			methodColor = "\033[97;44m" // 蓝色
		case "POST":
			methodColor = "\033[97;42m" // 绿色
		case "PUT":
			methodColor = "\033[97;43m" // 黄色
		case "DELETE":
			methodColor = "\033[97;41m" // 红色
		case "PATCH":
			methodColor = "\033[97;45m" // 紫色
		case "HEAD":
			methodColor = "\033[97;46m" // 青色
		default:
			methodColor = "\033[97;44m" // 蓝色
		}

		// 重置颜色
		resetColor := "\033[0m"

		_, _ = gin.DefaultWriter.Write([]byte(
			"[GIN] " + time.Now().Format("2006/01/02 - 15:04:05") +
				" | " + methodColor + method + resetColor +
				" | " + path +
				" | " + c.ClientIP() +
				" | " + c.Request.UserAgent() +
				" | " + time.Since(start).String() +
				" | " + statusColor + strconv.Itoa(statusCode) + resetColor +
				" | " + c.Errors.String() +
				"\n",
		))
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var (
	// httpRequestCounter HTTP请求计数器
	httpRequestCounter metric.Int64Counter
	// httpRequestDuration HTTP请求耗时直方图
	httpRequestDuration metric.Float64Histogram
	// httpRequestsInflight 当前并发处理的请求数
	httpRequestsInflight metric.Int64UpDownCounter
	httpMetricsOnce      sync.Once
)

func initHTTPMetrics() {
	m := metrics.Meter()
	if m == nil {
		return
	}
	c, _ := m.Int64Counter("http.server.requests_total")
	h, _ := m.Float64Histogram("http.server.duration", metric.WithUnit("s"))
	u, _ := m.Int64UpDownCounter("http.server.inflight")
	httpRequestCounter = c
	httpRequestDuration = h
	httpRequestsInflight = u
}

// Metrics 指标中间件
func Metrics() gin.HandlerFunc {
	httpMetricsOnce.Do(initHTTPMetrics)
	return func(c *gin.Context) {
		if httpRequestCounter == nil || httpRequestDuration == nil || httpRequestsInflight == nil {
			return
		}

		route := c.FullPath()
		method := c.Request.Method

		ctx := context.Background()

		commonAttrs := []attribute.KeyValue{
			attribute.String("http.request.method", method),
			attribute.String("http.route", route),
		}

		httpRequestsInflight.Add(ctx, 1, metric.WithAttributes(commonAttrs...))

		start := time.Now()

		defer func() {
			httpRequestsInflight.Add(ctx, -1, metric.WithAttributes(commonAttrs...))

			duration := time.Since(start).Seconds()
			status := c.Writer.Status()

			finalAttrs := append(commonAttrs, attribute.String("http.response.status_code", strconv.Itoa(status)))

			httpRequestCounter.Add(ctx, 1, metric.WithAttributes(finalAttrs...))
			httpRequestDuration.Record(ctx, duration, metric.WithAttributes(finalAttrs...))
		}()

		c.Next()
	}
}
