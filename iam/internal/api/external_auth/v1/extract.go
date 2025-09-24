package v1

import (
	"fmt"
	"net/http"
	"net/url"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/grpc/interceptor"
)

func (api *API) extractSessionID(req *authv3.CheckRequest) (uuid.UUID, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return uuid.Nil, fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers
	if cookieHeader, ok := headers[interceptor.HeaderCookie]; ok && cookieHeader != "" {
		req := &http.Request{Header: make(http.Header)}
		req.Header.Add(interceptor.HeaderCookie, cookieHeader)

		if cookie, err := req.Cookie(interceptor.SessionCookieName); err == nil {
			sessionID, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return uuid.Parse(cookie.Value)
			}
			return uuid.Parse(sessionID)
		}
	}

	return uuid.Nil, fmt.Errorf("session ID not found in cookies")
}
