package enrichstorage

import (
	grpccontroller "enrichstorage/internal/controllers/grpc"
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/enrichstorage/create"
	"enrichstorage/internal/service/enrichstorage/delete"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/internal/service/outboxsender"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
)

var Module = fx.Module("enrichstorage",
	fx.Provide(
		create.NewCreator,
		delete.NewDeleter,
		fx.Annotate(
			update.NewUpdater,
			fx.As(new(grpccontroller.Updater)),
		),
		fx.Annotate(
			get.NewGetter,
			fx.As(new(grpccontroller.PresenceChecker)),
		),
		getDB,
		fx.Annotate(
			repository.NewOutbox,
			fx.As(new(create.Outbox)),
			fx.As(new(outboxsender.Outbox)),
		),
		fx.Annotate(
			repository.NewRecords,
			fx.As(new(update.Repository)),
			fx.As(new(create.Repository)),
			fx.As(new(delete.Repository)),
			fx.As(new(get.Repository)),
		),
		repository.NewTxManager,
		grpccontroller.NewServer,
		gRPCServer,
	),
	fx.Invoke(func(_ *grpc.Server) {}),
)

func getDB(lc fx.Lifecycle, config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Storage.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	lc.Append(fx.StopHook(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}))
	return db, nil
}

func gRPCServer(lc fx.Lifecycle, serv *grpccontroller.Server, config *Config) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	v1.RegisterEnrichStorageServer(grpcServer, serv)
	serverErr := make(chan error, 1)
	lc.Append(fx.StartStopHook(
		func() error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPCPort))
			if err != nil {
				return err
			}
			go func() {
				defer close(serverErr)
				serverErr <- grpcServer.Serve(lis)
			}()
			return nil
		},
		func() error {
			grpcServer.GracefulStop()
			return <-serverErr
		},
	))
	return grpcServer, nil
}
