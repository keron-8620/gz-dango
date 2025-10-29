package menulogic

import (
	"context"
	"time"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMenuLogic) ListMenu(in *pb.ListMenuRequest) (*pb.PagMenuOutBase, error) {
	// todo: add your logic here and delete this line
	var (
		page int = database.DefaultPage
		size int = database.DefaultSize
	)
	if in.Page > 1 {
		page = int(in.Page)
	}
	if in.Size > 0 {
		size = int(in.Size)
	}
	query := make(map[string]any, 12)
	if in.Pk > 0 {
		query["id = ?"] = in.Pk
	}
	if in.Pks != "" {
		pks := database.StringToListUint(in.Pks)
		if len(pks) > 1 {
			query["id in ?"] = pks
		}
	}
	if in.BeforeCreatedAt != "" {
		bft, err := time.Parse(time.RFC3339, in.BeforeCreatedAt)
		if err == nil {
			query["created_at < ?"] = bft
		}
	}
	if in.AfterCreatedAt != "" {
		act, err := time.Parse(time.RFC3339, in.AfterCreatedAt)
		if err == nil {
			query["created_at > ?"] = act
		}
	}
	if in.BeforeUpdatedAt != "" {
		but, err := time.Parse(time.RFC3339, in.BeforeUpdatedAt)
		if err == nil {
			query["updated_at < ?"] = but
		}
	}
	if in.AfterUpdatedAt != "" {
		aut, err := time.Parse(time.RFC3339, in.AfterUpdatedAt)
		if err == nil {
			query["updated_at > ?"] = aut
		}
	}
	if in.ParentId != nil {
		pid := in.ParentId.GetValue()
		if pid == 0 {
			query["parent_id is null"] = nil
		} else {
			query["parent_id = ?"] = pid
		}
	}
	if in.Path != "" {
		query["path like ?"] = "%" + in.Path + "%"
	}
	if in.Component != "" {
		query["component like ?"] = "%" + in.Component + "%"
	}
	if in.Name != "" {
		query["name like ?"] = "%" + in.Name + "%"
	}
	if in.Label != "" {
		query["label = ?"] = in.Label
	}
	if in.IsActive != nil {
		query["is_active = ?"] = in.IsActive
	}
	if in.Descr != "" {
		query["descr like ?"] = "%" + in.Descr + "%"
	}
	qp := database.QueryParams{
		Preloads: []string{},
		Query:    query,
		OrderBy:  []string{"id"},
		Limit:    max(size, 0),
		Offset:   max(page-1, 0),
		IsCount:  true,
	}
	count, ms, err := l.svcCtx.Menu.ListModel(l.ctx, qp)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	mso := converter.ListMenuModelToOutBase(ms)
	return &pb.PagMenuOutBase{
		Items: mso,
		Page:  int64(page),
		Pages: database.CountPages(count, int64(size)),
		Size:  int64(size),
		Total: count,
	}, nil
}
