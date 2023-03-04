package impl

import (
	"database/sql"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/playmood/restful-demo/apps"
	"github.com/playmood/restful-demo/apps/book"
	"github.com/playmood/restful-demo/conf"
	"google.golang.org/grpc"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	db   *sql.DB
	log  logger.Logger
	book book.ServiceServer
	book.UnimplementedServiceServer
}

func (s *service) Config() {
	db := conf.C().MySQL.GetDB()
	s.log = zap.L().Named(s.Name())
	s.db = db
	s.book = apps.GetImpl(book.AppName).(book.ServiceServer)
}

func (s *service) Name() string {
	return book.AppName
}

func (s *service) Registry(server *grpc.Server) {
	book.RegisterServiceServer(server, svr)
}

func init() {
	// app.RegistryGrpcApp(svr)
	apps.RegistryImpl(svr)
	apps.RegistryGrpc(svr)
}
