package permissionlogic

import (
	"context"
	"time"

	"go-dango/apps/customer/rpc/internal/converter"
	"go-dango/apps/customer/rpc/internal/svc"
	"go-dango/apps/customer/rpc/pb"
	"go-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionLogic {
	return &ListPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPermissionLogic) ListPermission(in *pb.ListPermissionRequest) (*pb.PagPermissionOutBase, error) {
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
	query := make(map[string]any, 10)
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
	if in.Url != "" {
		query["url like ?"] = "%" + in.Url + "%"
	}
	if in.Method != "" {
		query["method = ?"] = in.Method
	}
	if in.Label != "" {
		query["label = ?"] = in.Label
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
	count, ms, err := l.svcCtx.Perm.ListModel(l.ctx, qp)
	if err != nil {
		rErr := database.NewGormError(err, nil)
		return nil, rErr
	}
	mso := converter.ListPermModelToOut(ms)
	return &pb.PagPermissionOutBase{
		Items: mso,
		Page:  int64(page),
		Pages: database.CountPages(count, int64(size)),
		Size:  int64(size),
		Total: count,
	}, nil
}
