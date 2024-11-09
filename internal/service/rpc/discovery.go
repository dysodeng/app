package rpc

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/trace"
	"github.com/dysodeng/rpc"
	"github.com/dysodeng/rpc/naming/etcd"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var builder resolver.Builder
var discovery rpc.ServiceDiscovery

func init() {
	builder = etcd.NewEtcdBuilder(config.Etcd.Grpc.Addr, etcd.WithBuilderNamespace(config.App.Name))
	resolver.Register(builder)

	discovery = rpc.NewServiceDiscovery(config.App.Name, builder)
}

func ServiceDiscovery() rpc.ServiceDiscovery {
	return discovery
}

func Error(err error) (error, int32) {
	grpcStatus := status.Convert(err)
	return errors.New(grpcStatus.Message()), int32(grpcStatus.Code())
}

func Ctx(ctx context.Context) context.Context {
	t := trace.New()
	return metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		"trace_id":         t.TraceId(ctx),
		"span_id":          t.SpanId(ctx),
		"span_name":        t.SpanName(ctx),
		"parent_span_id":   t.ParentSpanId(ctx),
		"parent_span_name": t.SpanName(ctx),
	}))
}

func FromCtx(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	var traceId, spanId, parentSpanId, spanName, parentSpanName string
	if vals, ok := md["trace_id"]; ok {
		traceId = vals[0]
	}
	if vals, ok := md["span_id"]; ok {
		spanId = vals[0]
	}
	if vals, ok := md["span_name"]; ok {
		spanName = vals[0]
	}
	if vals, ok := md["parent_span_id"]; ok {
		parentSpanId = vals[0]
	}
	if vals, ok := md["parent_span_name"]; ok {
		parentSpanName = vals[0]
	}

	valCtx := context.WithValue(ctx, "traceId", traceId)
	valCtx = context.WithValue(valCtx, "spanId", spanId)
	valCtx = context.WithValue(valCtx, "spanName", spanName)
	valCtx = context.WithValue(valCtx, "parentSpanId", parentSpanId)
	valCtx = context.WithValue(valCtx, "parentSpanName", parentSpanName)

	return valCtx
}
