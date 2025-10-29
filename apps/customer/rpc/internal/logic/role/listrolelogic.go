package rolelogic

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRoleLogic) ListRole(in *pb.ListRoleRequest) (*pb.PagRoleOutBase, error) {
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
	query := make(map[string]any, 8)
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
	if in.Name != "" {
		query["name like ?"] = "%" + in.Name + "%"
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
	count, ms, err := l.svcCtx.Role.ListModel(l.ctx, qp)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	mso := converter.ListRoleModelToOutBase(ms)
	return &pb.PagRoleOutBase{
		Items: mso,
		Page:  int64(page),
		Pages: database.CountPages(count, int64(size)),
		Size:  int64(size),
		Total: count,
	}, nil
}
