package tracing

import (
	"context"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingWrapper - простая обертка для трейсинга
type TracingWrapper[T any] struct {
	target T
	tracer trace.Tracer
	name   string
}

// NewTracingWrapper создает новую обертку
func NewTracingWrapper[T any](target T, name string) *TracingWrapper[T] {
	return &TracingWrapper[T]{
		target: target,
		tracer: otel.Tracer(name),
		name:   name,
	}
}

// Trace выполняет функцию с трейсингом
func (tw *TracingWrapper[T]) Trace(ctx context.Context, operation string, fn func(T, context.Context) error) error {
	return tw.execute(ctx, operation, func() (interface{}, error) {
		return nil, fn(tw.target, ctx)
	})
}

// TraceWithResult выполняет функцию с результатом и трейсингом
func TraceWithResult[T any, R any](tw *TracingWrapper[T], ctx context.Context, operation string, fn func(T, context.Context) (R, error)) (R, error) {
	result, err := tw.execute(ctx, operation, func() (interface{}, error) {
		return fn(tw.target, ctx)
	})

	if err != nil {
		var zero R
		return zero, err
	}

	return result.(R), nil
}

// execute центральная логика трейсинга
func (tw *TracingWrapper[T]) execute(ctx context.Context, operation string, fn func() (interface{}, error)) (interface{}, error) {
	ctx, span := tw.tracer.Start(ctx, operation)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("component", tw.name),
		attribute.String("operation", operation),
		attribute.String("target_type", reflect.TypeOf(tw.target).String()),
	)

	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Nanoseconds())/1e6),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return result, err
}
