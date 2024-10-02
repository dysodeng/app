package trace

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Interface interface {
	TraceId(ctx context.Context) string
	SpanId(ctx context.Context) string
	ParentSpanId(ctx context.Context) string
	SpanName(ctx context.Context) string
	NewSpan(ctx context.Context, span string) context.Context
}

type trace struct{}

var _ Interface = (*trace)(nil)

func New() Interface {
	return &trace{}
}

func (trace *trace) TraceId(ctx context.Context) string {
	traceId := ctx.Value("traceId")
	if traceId == nil {
		return ""
	}
	return traceId.(string)
}

func (trace *trace) SpanId(ctx context.Context) string {
	spanId := ctx.Value("spanId")
	if spanId == nil {
		return ""
	}
	return spanId.(string)
}

func (trace *trace) ParentSpanId(ctx context.Context) string {
	parentSpanId := ctx.Value("parentSpanId")
	if parentSpanId == nil {
		return ""
	}
	return parentSpanId.(string)
}

func (trace *trace) SpanName(ctx context.Context) string {
	spanName := ctx.Value("spanName")
	if spanName == nil {
		return ""
	}
	return spanName.(string)
}

func (trace *trace) NewSpan(ctx context.Context, name string) context.Context {
	traceId := trace.TraceId(ctx)
	parentSpanId := trace.SpanId(ctx)
	spanId := GenerateSpanID()
	if traceId == "" {
		traceId = spanId
	}
	traceCtx := context.WithValue(ctx, "traceId", traceId)
	traceCtx = context.WithValue(traceCtx, "spanId", spanId)
	traceCtx = context.WithValue(traceCtx, "parentSpanId", parentSpanId)
	traceCtx = context.WithValue(traceCtx, "spanName", name)
	return traceCtx
}

func GenerateSpanID() string {
	uuidLong := uuidToLong(uuid.New())
	times := uint64(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixMilli()))
	spanId := ((times ^ uint64(uuidLong)) << 32) | uint64(rand.Int31())
	return strconv.FormatUint(spanId, 16)
}

func uuidToLong(u uuid.UUID) int64 {
	// 假设我们只使用UUID的前8字节
	return int64(u[0]) | int64(u[1])<<8 | int64(u[2])<<16 | int64(u[3])<<24 |
		int64(u[4])<<32 | int64(u[5])<<40 | int64(u[6])<<48 | int64(u[7])<<56
}
