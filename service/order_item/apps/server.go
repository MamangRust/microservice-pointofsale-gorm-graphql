package apps

import (
	
	"context"

	"github.com/MamangRust/microservice-point-of-sale-order-item/handler"
	mencache "github.com/MamangRust/microservice-point-of-sale-order-item/cache"
	"github.com/MamangRust/microservice-point-of-sale-order-item/repository"
	"github.com/MamangRust/microservice-point-of-sale-order-item/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order-item")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterOrderItemServiceServer(gs, handlers.OrderItem)
	}

	return srv, nil
}
