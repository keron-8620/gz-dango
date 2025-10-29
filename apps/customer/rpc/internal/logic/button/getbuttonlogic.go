package buttonlogic

import (
	"context"

	"gz-dango/apps/customer/rpc/internal/converter"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"
	"gz-dango/pkg/database"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetButtonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetButtonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetButtonLogic {
	return &GetButtonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetButtonLogic) GetButton(in *pb.GetButtonRequest) (*pb.ButtonOut, error) {
	// todo: add your logic here and delete this line
	m, err := l.svcCtx.Button.FindModel(l.ctx, []string{"Menu", "Permissons"}, in.Pk)
	if err != nil {
		return nil, database.NewGormError(err, nil)
	}
	return converter.ButtonModelToOut(*m), nil
}
