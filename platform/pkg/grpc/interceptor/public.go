package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

const (
	// publicMethodKey ключ для пометки публичных методов в контексте
	publicMethodKey contextKey = "public-method"
)

// PublicFilter фильтрует публичные и системные методы
// Пропускает их без обработки, остальные передает дальше по цепочке
type PublicFilter struct {
	// Кеш публичных методов (из proto аннотаций)
	publicMethodsCache map[string]bool
	// Системные методы (не требуют аутентификации)
	systemMethods map[string]bool
}

// NewPublicFilter создает новый фильтр публичных методов
func NewPublicFilter() *PublicFilter {
	filter := &PublicFilter{
		publicMethodsCache: make(map[string]bool),
		systemMethods: map[string]bool{
			"/envoy.service.auth.v3.Authorization/Check":                     true,
			"/grpc.health.v1.Health/Check":                                   true,
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
		},
	}

	filter.buildPublicMethodsCache()

	return filter
}

// Unary возвращает interceptor для фильтрации публичных методов
func (f *PublicFilter) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if f.systemMethods[info.FullMethod] {
			ctx = context.WithValue(ctx, publicMethodKey, true)
			return handler(ctx, req)
		}

		if isPublic, exists := f.publicMethodsCache[info.FullMethod]; exists && isPublic {
			ctx = context.WithValue(ctx, publicMethodKey, true)
			return handler(ctx, req)
		}

		return handler(ctx, req)
	}
}

// buildPublicMethodsCache заполняет кеш публичными методами из proto аннотаций
func (f *PublicFilter) buildPublicMethodsCache() {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for j := 0; j < services.Len(); j++ {
			service := services.Get(j)

			methods := service.Methods()
			for k := 0; k < methods.Len(); k++ {
				method := methods.Get(k)

				options := method.Options().(*descriptorpb.MethodOptions)
				if options == nil {
					continue
				}

				ext := proto.GetExtension(options, commonV1.E_Public)
				if isPublic, ok := ext.(bool); ok && isPublic {
					fullMethod := "/" + string(service.FullName()) + "/" + string(method.Name())
					f.publicMethodsCache[fullMethod] = true
				}
			}
		}
		return true
	})
}

// IsPublicMethod проверяет, помечен ли метод как публичный в контексте
func IsPublicMethod(ctx context.Context) bool {
	isPublic, ok := ctx.Value(publicMethodKey).(bool)
	return ok && isPublic
}
