package v1

import (
	"fmt"
	"net/http"
	"net/url"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (api *API) extractSessionID(req *authv3.CheckRequest) (uuid.UUID, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return uuid.Nil, fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if sessionID, ok := headers[model.HeaderSessionID]; ok && sessionID != "" {
		return uuid.Parse(sessionID)
	}

	if cookieHeader, ok := headers[model.HeaderCookie]; ok && cookieHeader != "" {
		req := &http.Request{Header: make(http.Header)}
		req.Header.Add(model.HeaderCookie, cookieHeader)

		if cookie, err := req.Cookie(model.SessionCookieName); err == nil {
			sessionID, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return uuid.Parse(cookie.Value)
			}
			return uuid.Parse(sessionID)
		}
	}

	return uuid.Nil, fmt.Errorf("session ID not found in headers or cookies")
}
