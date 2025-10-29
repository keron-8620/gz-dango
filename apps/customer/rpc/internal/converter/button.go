package converter

import (
	"go-dango/apps/customer/rpc/internal/models"
	"go-dango/apps/customer/rpc/pb"
)

func ButtonModelToOutBase(
	m models.ButtonModel,
) *pb.ButtonOutBase {
	return &pb.ButtonOutBase{
		Id:           m.Id,
		CreatedAt:    m.CreatedAt.String(),
		UpdatedAt:    m.UpdatedAt.String(),
		Name:         m.Name,
		ArrangeOrder: m.ArrangeOrder,
		IsActive:     m.IsActive,
		Descr:        m.Descr,
	}
}

func ListButtonModelToOutBase(
	ms []models.ButtonModel,
) []*pb.ButtonOutBase {
	mso := make([]*pb.ButtonOutBase, 0, len(ms))
	if len(ms) > 0 {
		for _, m := range ms {
			mo := ButtonModelToOutBase(m)
			mso = append(mso, mo)
		}
	}
	return mso
}

func ButtonModelToOut(
	m models.ButtonModel,
) *pb.ButtonOut {
	return &pb.ButtonOut{
		Id:           m.Id,
		CreatedAt:    m.CreatedAt.String(),
		UpdatedAt:    m.UpdatedAt.String(),
		Name:         m.Name,
		ArrangeOrder: m.ArrangeOrder,
		IsActive:     m.IsActive,
		Descr:        m.Descr,
		Menu:         MenuModelToOutBase(m.Menu),
		Permissions:  ListPermModelToOut(m.Permissions),
	}
}
