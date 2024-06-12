package enrichstorage

import (
	"context"
	grpccontroller "enrichstorage/internal/controllers/grpc"
	v1 "enrichstorage/internal/controllers/grpc/protoc"
	createhandler "enrichstorage/internal/controllers/http/create"
	deletehandler "enrichstorage/internal/controllers/http/delete"
	gethandler "enrichstorage/internal/controllers/http/get"
	updatehandler "enrichstorage/internal/controllers/http/update"
	"enrichstorage/internal/providers/repository"
	"enrichstorage/internal/service/enrichstorage/create"
	"enrichstorage/internal/service/enrichstorage/delete"
	"enrichstorage/internal/service/enrichstorage/get"
	"enrichstorage/internal/service/enrichstorage/update"
	"enrichstorage/internal/service/outboxsender"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
			fx.As(new(gethandler.GetterService)),
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
		fx.Annotate(
			func() *repository.FioComparator { return &repository.FioComparator{} },
			fx.As(new(get.KeysComparator)),
		),
		fx.Annotate(
			repository.NewTxManager,
			fx.As(new(create.TxManager)),
		),
		grpccontroller.NewServer,
		gRPCServer,
		createhandler.NewHandler,
		deletehandler.NewHandler,
		updatehandler.NewHandler,
		gethandler.NewHandler,
		ginHandler,
		httpServer,

		fx.Annotate(
			createSwaggerURL,
			fx.ResultTags(`name:"swaggerURL"`),
		),
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
	CreateHdlr *createhandler.Handler
	DeleteHdlr *deletehandler.Handler
	UpdateHdlr *updatehandler.Handler
	GetHdlr    *gethandler.Handler
	SwaggerURL string `name:"swaggerURL"`
}

func ginHandler(params ginHandlerParams) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	records := router.Group("/api/v1/records")
	{
		records.POST("/create", params.CreateHdlr.Handle)
		records.POST("/delete", params.DeleteHdlr.Handle)
		records.POST("/update", params.UpdateHdlr.Handle)
		records.GET("/get", params.GetHdlr.Handle)
	}
	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.NewHandler()),
	)
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
