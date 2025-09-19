package v1

import (
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (api *API) allowRequest(whoami *model.WhoAMI) *authv3.CheckResponse {
	roleNames := make([]string, len(whoami.RolesWithPermissions))
	permissionsSet := make(map[string]bool)

	for i, role := range whoami.RolesWithPermissions {
		roleNames[i] = role.Role.Name
		for _, perm := range role.Permissions {
			permissionsSet[perm.Resource+":"+perm.Action] = true
		}
	}

	permissions := make([]string, 0, len(permissionsSet))
	for perm := range permissionsSet {
		permissions = append(permissions, perm)
	}

	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   model.HeaderUserUUID,
				Value: whoami.User.ID.String(),
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   model.HeaderUserLogin,
				Value: whoami.User.Login,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   model.HeaderUserRoles,
				Value: strings.Join(roleNames, ","),
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   model.HeaderUserPermissions,
				Value: strings.Join(permissions, ","),
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{model.HeaderCookie, model.HeaderAuthorization},
			},
		},
	}
}
