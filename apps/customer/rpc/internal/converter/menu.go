package converter

import (
	"go-dango/apps/customer/rpc/internal/models"
	"go-dango/apps/customer/rpc/pb"
)

func MenuModelToOutBase(
	m models.MenuModel,
) *pb.MenuOutBase {
	return &pb.MenuOutBase{
		Id:        m.Id,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		Path:      m.Path,
		Component: m.Component,
		Name:      m.Name,
		Label:     m.Label,
		Meta: &pb.MetaSchemas{
			Title: m.Meta.Title,
			Icon:  m.Meta.Icon,
		},
		ArrangeOrder: m.ArrangeOrder,
		IsActive:     m.IsActive,
		Descr:        m.Descr,
	}
}

func ListMenuModelToOutBase(
	ms []models.MenuModel,
) []*pb.MenuOutBase {
	mso := make([]*pb.MenuOutBase, 0, len(ms))
	if len(ms) > 0 {
		for _, m := range ms {
			mo := MenuModelToOutBase(m)
			mso = append(mso, mo)
		}
	}
	return mso
}

func MenuModelToOut(
	m models.MenuModel,
) *pb.MenuOut {
	var parent *pb.MenuOutBase
	if m.Parent != nil {
		parent = MenuModelToOutBase(*m.Parent)
	}
	return &pb.MenuOut{
		Id:        m.Id,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
		Path:      m.Path,
		Component: m.Component,
		Name:      m.Name,
		Label:     m.Label,
		Meta: &pb.MetaSchemas{
			Title: m.Meta.Title,
			Icon:  m.Meta.Icon,
		},
		ArrangeOrder: m.ArrangeOrder,
		IsActive:     m.IsActive,
		Descr:        m.Descr,
		Parent:       parent,
		Permissions:  ListPermModelToOut(m.Permissions),
	}
}
