package apps

import (
	
	"context"
	"os"

	mencache "github.com/MamangRust/microservice-point-of-sale-merchant/cache"
	"github.com/MamangRust/microservice-point-of-sale-merchant/handler"
	"github.com/MamangRust/microservice-point-of-sale-merchant/repository"
	"github.com/MamangRust/microservice-point-of-sale-merchant/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/kafka"
	"github.com/MamangRust/microservice-point-of-sale-pkg/outbox"
	"github.com/MamangRust/microservice-point-of-sale-pkg/server"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := os.Getenv("GRPC_USER_ADDR")
	if userAddr == "" {
		userAddr = "localhost:50053"
	}

	srv.Logger.Info("Connecting to User service via gRPC", zap.String("addr", userAddr))
	userConn, err := server.NewGRPCClient(userAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing merchant service remote gRPC connections")
		userConn.Close()
	}()

	userClient := pbuser.NewUserServiceClient(userConn)

	repos := repository.NewRepositories(srv.GormDB, userClient)
	// Kafka bersifat opsional: tanpa KAFKA_BROKERS (mis. E2E lokal tanpa kafka)
	// service tetap jalan dan event email di-skip (guard s.kafka != nil).
	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.GormDB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Kafka:         myKafka,
		Repositories:  repos,
		Outbox:        outboxService,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterMerchantServiceServer(gs, handlers.Merchant)
		pb.RegisterMerchantDocumentServiceServer(gs, handlers.MerchantDocument)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
