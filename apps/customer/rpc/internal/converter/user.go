package converter

import (
	"go-dango/apps/customer/rpc/internal/models"
	"go-dango/apps/customer/rpc/pb"
)

func UserModelToOut(
	m models.UserModel,
) *pb.UserOut {
	return &pb.UserOut{
		Id:        m.Id,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		Username:  m.Username,
		IsActive:  m.IsActive,
		IsStaff:   m.IsStaff,
		Role:      RoleModelToOutBase(m.Role),
	}
}

func ListUserModelToOut(
	ms []models.UserModel,
) []*pb.UserOut {
	mso := make([]*pb.UserOut, 0, len(ms))
	if len(ms) > 0 {
		for _, m := range ms {
			mo := UserModelToOut(m)
			mso = append(mso, mo)
		}
	}
	return mso
}
