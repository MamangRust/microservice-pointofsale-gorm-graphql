package apps

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"github.com/MamangRust/microservice-point-of-sale-user/handler"
	mencache "github.com/MamangRust/microservice-point-of-sale-user/cache"
	"github.com/MamangRust/microservice-point-of-sale-user/repository"
	"github.com/MamangRust/microservice-point-of-sale-user/service"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)
	hash := hash.NewHashingPassword()
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("user")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterUserServiceServer(gs, handlers.User)
	}

	return srv, nil
}
