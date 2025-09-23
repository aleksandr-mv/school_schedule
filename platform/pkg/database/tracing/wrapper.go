package tracing

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingWrapper - generic обертка для любого типа
type TracingWrapper[T any] struct {
	target T
	tracer trace.Tracer
	name   string
}

// NewTracingWrapper создает новую обертку для трейсинга
func NewTracingWrapper[T any](target T, name string) *TracingWrapper[T] {
	return &TracingWrapper[T]{
		target: target,
		tracer: otel.Tracer(name),
		name:   name,
	}
}

// Execute для методов без возвращаемого значения
func (tw *TracingWrapper[T]) Execute(ctx context.Context, methodName string, fn func(T, context.Context) error) error {
	_, err := tw.executeWithTracing(ctx, methodName, func() (interface{}, error) {
		return nil, fn(tw.target, ctx)
	})
	return err
}

// ExecuteWithResult для методов с возвращаемым значением
func ExecuteWithResult[T any, R any](tw *TracingWrapper[T], ctx context.Context, methodName string, fn func(T, context.Context) (R, error)) (R, error) {
	result, err := tw.executeWithTracing(ctx, methodName, func() (interface{}, error) {
		return fn(tw.target, ctx)
	})

	if err != nil {
		var zero R
		return zero, err
	}

	return result.(R), nil
}

// executeWithTracing центральная логика трейсинга
func (tw *TracingWrapper[T]) executeWithTracing(ctx context.Context, methodName string, fn func() (interface{}, error)) (interface{}, error) {
	operationName := tw.formatOperationName(methodName)

	ctx, span := tw.tracer.Start(ctx, operationName)
	defer span.End()

	// Добавляем базовые атрибуты
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("component", tw.name),
		attribute.String("method", methodName),
		attribute.String("target_type", reflect.TypeOf(tw.target).String()),
	)

	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	// Записываем метрики
	span.SetAttributes(
		attribute.Float64("db.duration_ms", float64(duration.Nanoseconds())/1e6),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("db.status", "error"))
	} else {
		span.SetStatus(codes.Ok, "success")
		span.SetAttributes(attribute.String("db.status", "success"))
	}

	return result, err
}

// TraceCall автоматическое определение имени операции из caller'а
func (tw *TracingWrapper[T]) TraceCall(ctx context.Context, fn func(T, context.Context) error) error {
	methodName := tw.getCallerMethodName()
	return tw.Execute(ctx, methodName, fn)
}

// TraceCallWithResult автоматическое определение имени операции для методов с результатом
func TraceCallWithResult[T any, R any](tw *TracingWrapper[T], ctx context.Context, fn func(T, context.Context) (R, error)) (R, error) {
	methodName := tw.getCallerMethodName()
	return ExecuteWithResult(tw, ctx, methodName, fn)
}

// getCallerMethodName получает имя вызывающего метода
func (tw *TracingWrapper[T]) getCallerMethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "unknown"
}

// formatOperationName конвертирует имя метода в snake_case
func (tw *TracingWrapper[T]) formatOperationName(methodName string) string {
	// Конвертируем в snake_case
	var result strings.Builder
	for i, r := range methodName {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('.')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// Unwrap получает доступ к оригинальному объекту
func (tw *TracingWrapper[T]) Unwrap() T {
	return tw.target
}
