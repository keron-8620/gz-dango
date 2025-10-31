package userlogic

import (
	"time"

	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
)

func UserModelToClaims(m *models.UserModel, exp time.Duration) *auth.UserClaims {
	return &auth.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   m.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		UserId:  m.Id,
		IsStaff: m.IsStaff,
		Role:    m.Role.Name,
	}
}

