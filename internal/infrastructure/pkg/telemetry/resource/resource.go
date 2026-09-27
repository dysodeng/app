package resource

import (
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

var globalResource *resource.Resource
var once sync.Once

func Resource() *resource.Resource {
	once.Do(func() {
		var err error
		cfg := config.Current()
		globalResource, err = resource.Merge(
			resource.Default(),
			resource.NewSchemaless(
				semconv.ServiceName(serviceName(cfg)),
				semconv.ServiceVersion(cfg.Monitor.ServiceVersion),
				attribute.String("env", cfg.App.Environment),
			),
		)
		if err != nil {
			panic(err)
		}
	})
	return globalResource
}
