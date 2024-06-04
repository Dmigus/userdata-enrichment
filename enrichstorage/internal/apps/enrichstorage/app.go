package enrichstorage

import (
	"context"
	grpccontroller "enrichstorage/internal/controllers/grpc"
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	createhandler "enrichstorage/internal/controllers/http/create"
	deletehandler "enrichstorage/internal/controllers/http/delete"
	updatehandler "enrichstorage/internal/controllers/http/update"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/enrichstorage/create"
	"enrichstorage/internal/service/enrichstorage/delete"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/internal/service/outboxsender"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"net/http"
)

var Module = fx.Module("enrichstorage",
	fx.Provide(
		fx.Annotate(
			create.NewCreator,
			fx.As(new(createhandler.CreatorService)),
		),
		fx.Annotate(
			delete.NewDeleter,
			fx.As(new(deletehandler.DeleteService)),
		),

		fx.Annotate(
			update.NewUpdater,
			fx.As(new(grpccontroller.Updater)),
			fx.As(new(updatehandler.UpdateService)),
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
		createhandler.NewHandler,
		deletehandler.NewHandler,
		updatehandler.NewHandler,
		ginHandler,
		httpServer,
	),
	fx.Invoke(func(*grpc.Server) {}),
	fx.Invoke(func(*http.Server) {}),
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

type ginHandlerParams struct {
	fx.In
	createHdlr *createhandler.Handler
	deleteHdlr *deletehandler.Handler
	updateHdlr *updatehandler.Handler
}

func ginHandler(params ginHandlerParams) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/records/create", params.createHdlr.Handle)
	router.POST("/records/delete", params.deleteHdlr.Handle)
	router.POST("/records/update", params.updateHdlr.Handle)
	return router.Handler()
}

func httpServer(lc fx.Lifecycle, config *Config, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTPPort),
		Handler: handler,
	}
	lc.Append(fx.StartStopHook(
		func() {
			go func() {
				_ = server.ListenAndServe()
			}()
		},
		func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	))
	return server
}
