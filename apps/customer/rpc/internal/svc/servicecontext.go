package svc

import (
	"gz-dango/apps/customer/rpc/internal/config"
	"gz-dango/pkg/auth"
	"gz-dango/pkg/crypto"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Auth   *auth.AuthEnforcer
	Hasher crypto.Hasher

	Perm   *PermissionService
	Menu   *MenuService
	Button *ButtonService
	Role   *RoleService
	User   *UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	hasher := crypto.NewBcryptHasher(12)
	return &ServiceContext{
		Config: c,
		Hasher: hasher,
	}
}
