package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

func CopyCtx(ctx context.Context) context.Context {
	return opentracing.ContextWithSpan(context.TODO(), opentracing.SpanFromContext(ctx))
}

func NewCtx(operation string) context.Context {
	return opentracing.ContextWithSpan(context.TODO(), opentracing.GlobalTracer().StartSpan(operation))
}
