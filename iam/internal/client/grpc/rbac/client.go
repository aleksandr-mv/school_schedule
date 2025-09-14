package rbac

import (
	def "github.com/aleksandr-mv/school_schedule/iam/internal/client/grpc"
	rbacV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

var _ def.RBACClient = (*client)(nil)

type client struct {
	generatedClient rbacV1.UserRoleServiceClient
}

func NewClient(generatedClient rbacV1.UserRoleServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
