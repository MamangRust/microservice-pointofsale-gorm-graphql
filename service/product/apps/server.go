package apps

import (
	
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-product/handler"
	mencache "github.com/MamangRust/microservice-point-of-sale-product/cache"
	"github.com/MamangRust/microservice-point-of-sale-product/repository"
	"github.com/MamangRust/microservice-point-of-sale-product/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterProductServiceServer(gs, handlers.Product)
	}

	return srv, nil
}
