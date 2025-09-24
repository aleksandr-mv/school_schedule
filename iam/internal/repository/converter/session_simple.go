package converter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

func ToRedisHash(whoami *model.WhoAMI, sessionID uuid.UUID, expiresAt time.Time) (map[string]interface{}, error) {
	now := time.Now()
	whoami.Session.ID = sessionID
	whoami.Session.CreatedAt = now
	whoami.Session.UpdatedAt = now
	whoami.Session.ExpiresAt = expiresAt
	rolesJSON := ""
	if len(whoami.RolesWithPermissions) > 0 {
		bytes, err := json.Marshal(whoami.RolesWithPermissions)
		if err != nil {
			return nil, fmt.Errorf("marshal roles: %w", err)
		}
		rolesJSON = string(bytes)
	}

	notifJSON := ""
	if len(whoami.User.NotificationMethods) > 0 {
		bytes, err := json.Marshal(whoami.User.NotificationMethods)
		if err != nil {
			return nil, fmt.Errorf("marshal notifications: %w", err)
		}
		notifJSON = string(bytes)
	}

	hash := map[string]interface{}{
		"session_id":           sessionID.String(),
		"session_created_at":   now.UnixNano(),
		"session_updated_at":   now.UnixNano(),
		"session_expires_at":   expiresAt.UnixNano(),
		"user_id":              whoami.User.ID.String(),
		"user_login":           whoami.User.Login,
		"user_email":           whoami.User.Email,
		"user_created_at":      whoami.User.CreatedAt.UnixNano(),
		"roles":                rolesJSON,
		"notification_methods": notifJSON,
	}

	if whoami.User.UpdatedAt != nil {
		hash["user_updated_at"] = whoami.User.UpdatedAt.UnixNano()
	}

	return hash, nil
}

func FromRedisHash(hash map[string]string) (*model.WhoAMI, error) {
	sessionID, err := uuid.Parse(hash["session_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	userID, err := uuid.Parse(hash["user_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	sessionCreatedAt := time.Unix(0, parseInt64(hash["session_created_at"]))
	sessionUpdatedAt := time.Unix(0, parseInt64(hash["session_updated_at"]))
	sessionExpiresAt := time.Unix(0, parseInt64(hash["session_expires_at"]))
	userCreatedAt := time.Unix(0, parseInt64(hash["user_created_at"]))

	var userUpdatedAt *time.Time
	if updatedAtStr := hash["user_updated_at"]; updatedAtStr != "" {
		t := time.Unix(0, parseInt64(updatedAtStr))
		userUpdatedAt = &t
	}

	var roles []*model.RoleWithPermissions
	if rolesStr := hash["roles"]; rolesStr != "" {
		var tempRoles []model.RoleWithPermissions
		if err := json.Unmarshal([]byte(rolesStr), &tempRoles); err == nil {
			roles = make([]*model.RoleWithPermissions, len(tempRoles))
			for i := range tempRoles {
				roles[i] = &tempRoles[i]
			}
		}
	}

	var notifications []*model.NotificationMethod
	if notifStr := hash["notification_methods"]; notifStr != "" {
		var tempNotif []model.NotificationMethod
		if err := json.Unmarshal([]byte(notifStr), &tempNotif); err == nil {
			notifications = make([]*model.NotificationMethod, len(tempNotif))
			for i := range tempNotif {
				notifications[i] = &tempNotif[i]
			}
		}
	}

	return &model.WhoAMI{
		Session: model.Session{
			ID:        sessionID,
			CreatedAt: sessionCreatedAt,
			UpdatedAt: sessionUpdatedAt,
			ExpiresAt: sessionExpiresAt,
		},
		User: model.User{
			ID:                  userID,
			Login:               hash["user_login"],
			Email:               hash["user_email"],
			NotificationMethods: notifications,
			CreatedAt:           userCreatedAt,
			UpdatedAt:           userUpdatedAt,
		},
		RolesWithPermissions: roles,
	}, nil
}

func parseInt64(s string) int64 {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}
