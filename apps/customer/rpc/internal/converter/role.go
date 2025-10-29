package converter

import (
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/pb"
)

func RoleModelToOutBase(
	m models.RoleModel,
) *pb.RoleOutBase {
	return &pb.RoleOutBase{
		Id:        m.Id,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		Name:      m.Name,
		Descr:     m.Descr,
	}
}

func ListRoleModelToOutBase(
	ms []models.RoleModel,
) []*pb.RoleOutBase {
	mso := make([]*pb.RoleOutBase, 0, len(ms))
	if len(ms) > 0 {
		for _, m := range ms {
			mo := RoleModelToOutBase(m)
			mso = append(mso, mo)
		}
	}
	return mso
}

func RoleModelToOut(
	m models.RoleModel,
) *pb.RoleOut {
	return &pb.RoleOut{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt.String(),
		UpdatedAt:   m.UpdatedAt.String(),
		Name:        m.Name,
		Descr:       m.Descr,
		Permissions: ListPermModelToOut(m.Permissions),
		Menus:       ListMenuModelToOutBase(m.Menus),
		Buttons:     ListButtonModelToOutBase(m.Buttons),
	}
}
