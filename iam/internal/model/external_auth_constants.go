package model

// External Auth константы для работы с Envoy
const (
	// Cookies
	SessionCookieName = "X-Session-Uuid"

	// Заголовки запроса
	HeaderSessionID = "sid"

	// Заголовки пользователя (добавляются Envoy для внутренних сервисов)
	HeaderUserUUID        = "X-User-UUID"
	HeaderUserLogin       = "X-User-Login"
	HeaderUserRoles       = "X-User-Roles"
	HeaderUserPermissions = "X-User-Permissions"

	// HTTP заголовки
	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"
	HeaderContentType   = "content-type"
	HeaderAuthStatus    = "X-Auth-Status"

	// Значения
	ContentTypeJSON  = "application/json"
	AuthStatusDenied = "denied"
)
