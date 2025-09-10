package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

// PermissionInterceptor interceptor для проверки прав доступа
type PermissionInterceptor struct {
	// Статический кеш аннотаций для методов (заполняется при инициализации)
	permissionCache map[string]string
}

// NewPermissionInterceptor создает новый interceptor проверки прав доступа
func NewPermissionInterceptor() *PermissionInterceptor {
	interceptor := &PermissionInterceptor{
		permissionCache: make(map[string]string),
	}

	// Предварительно заполняем кеш всеми аннотациями
	interceptor.buildPermissionCache()

	return interceptor
}

// UnaryServerInterceptor возвращает unary server interceptor для проверки прав доступа
func (i *PermissionInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Извлекаем аннотацию permission из кеша
		permission, exists := i.permissionCache[info.FullMethod]
		if !exists || permission == "" {
			// Аннотация не указана, пропускаем проверку
			return handler(ctx, req)
		}

		// Проверяем права доступа
		if err := i.checkPermission(ctx, permission); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// checkPermission проверяет права доступа пользователя
func (i *PermissionInterceptor) checkPermission(ctx context.Context, permission string) error {
	// Получаем права пользователя из контекста
	permissions, ok := GetPermissionsFromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "Права доступа не найдены в контексте")
	}

	// Простая проверка прав (для большинства случаев достаточно)
	for _, userPermission := range permissions {
		if userPermission != nil && userPermission.Resource+":"+userPermission.Action == permission {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "Недостаточно прав доступа")
}

// buildPermissionCache предварительно заполняет кеш всеми аннотациями из protobuf
func (i *PermissionInterceptor) buildPermissionCache() {
	// Проходим по всем зарегистрированным файлам
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		// Проходим по всем сервисам в файле
		services := fd.Services()
		for j := 0; j < services.Len(); j++ {
			service := services.Get(j)

			// Проходим по всем методам сервиса
			methods := service.Methods()
			for k := 0; k < methods.Len(); k++ {
				method := methods.Get(k)

				// Извлекаем опции метода
				options := method.Options().(*descriptorpb.MethodOptions)
				if options == nil {
					continue
				}

				// Ищем нашу кастомную аннотацию permission
				ext := proto.GetExtension(options, commonV1.E_Permission)
				if permission, ok := ext.(string); ok && permission != "" {
					// Формируем полное имя метода
					fullMethod := "/" + string(service.FullName()) + "/" + string(method.Name())
					i.permissionCache[fullMethod] = permission
				}
			}
		}
		return true
	})
}
