package converter

import (
	"gz-dango/apps/customer/rpc/internal/models"
	"gz-dango/apps/customer/rpc/pb"
)

func PermModelToOutBase(
	m models.PermissionModel,
) *pb.PermissionOutBase {
	return &pb.PermissionOutBase{
		Id:        m.Id,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		Url:       m.Url,
		Method:    m.Method,
		Descr:     m.Descr,
	}
}

func ListPermModelToOut(
	ms []models.PermissionModel,
) []*pb.PermissionOutBase {
	mso := make([]*pb.PermissionOutBase, 0, len(ms))
	if len(ms) > 0 {
		for _, m := range ms {
			mo := PermModelToOutBase(m)
			mso = append(mso, mo)
		}
	}
	return mso
}
