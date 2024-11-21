package trace

import (
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Error sets the error status on the span and adds an event with the error details.
func Error(err error, span trace.Span) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)

		span.SetAttributes(
			attribute.String("error.message", err.Error()),
			attribute.String("error.file", file),
			attribute.Int64("error.line", int64(line)),
		)

		span.AddEvent("Error occurred", trace.WithAttributes(
			attribute.String("error.message", err.Error()),
			attribute.String("error.file", file),
			attribute.Int64("error.line", int64(line)),
		))

		span.SetStatus(codes.Error, err.Error())
	}
}
