package svc

import (
	"go-dango/apps/customer/rpc/internal/config"
	"go-dango/pkg/auth"
	"go-dango/pkg/crypto"

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
