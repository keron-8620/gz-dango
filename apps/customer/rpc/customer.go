package main

import (
	"flag"
	"fmt"

	"gz-dango/apps/customer/rpc/internal/config"
	buttonServer "gz-dango/apps/customer/rpc/internal/server/button"
	menuServer "gz-dango/apps/customer/rpc/internal/server/menu"
	permissionServer "gz-dango/apps/customer/rpc/internal/server/permission"
	roleServer "gz-dango/apps/customer/rpc/internal/server/role"
	userServer "gz-dango/apps/customer/rpc/internal/server/user"
	"gz-dango/apps/customer/rpc/internal/svc"
	"gz-dango/apps/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/customer.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.ServiceLog)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterPermissionServer(grpcServer, permissionServer.NewPermissionServer(ctx))
		pb.RegisterMenuServer(grpcServer, menuServer.NewMenuServer(ctx))
		pb.RegisterButtonServer(grpcServer, buttonServer.NewButtonServer(ctx))
		pb.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		pb.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer func() {
		ctx.Close()
		s.Stop()
	}()

	if err := ctx.RefreshPolicies(); err != nil {
		panic(err)
	}

	go ctx.WatchCasbinPolicies()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
