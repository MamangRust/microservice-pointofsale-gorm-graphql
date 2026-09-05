package transaction_test

import (
	"context"
		"testing"

	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	trans_cache "github.com/MamangRust/microservice-point-of-sale-transacton/cache"
	trans_handler "github.com/MamangRust/microservice-point-of-sale-transacton/handler"
	trans_repo "github.com/MamangRust/microservice-point-of-sale-transacton/repository"
	trans_service "github.com/MamangRust/microservice-point-of-sale-transacton/service"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionGapiTestSuite struct {
	tests.BaseTestSuite
	client pb.TransactionServiceClient
}

func (s *TransactionGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	gormDB := s.GormDB()

	cashierClient := pbcashier.NewCashierServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchant.NewMerchantServiceClient(s.Conns["merchant"])
	orderClient := pborder.NewOrderServiceClient(s.Conns["order"])
	orderItemClient := pborderitem.NewOrderItemServiceClient(s.Conns["order-item"])

	// Transaction dependencies
	mencache := trans_cache.NewMencache(cacheStore)
	repos := trans_repo.NewRepositories(gormDB, cashierClient, merchantClient, orderClient, orderItemClient)
	svc := trans_service.NewService(&trans_service.Deps{
		Kafka:         nil,
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := trans_handler.NewHandler(&trans_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterTransactionServiceServer(server, handler.Transaction)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pb.NewTransactionServiceClient(conn)
}

func (s *TransactionGapiTestSuite) TestTransactionGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	var cashierID int
	err := s.GormDB().WithContext(ctx).Raw(
		`SELECT cashier_id FROM cashiers WHERE user_id = ? AND merchant_id = ? AND deleted_at IS NULL LIMIT 1`,
		userID, merchID,
	).Scan(&cashierID).Error

	// 2. Create
	createRes, err := s.client.Create(ctx, &pb.CreateTransactionRequest{
		OrderId:       int32(orderID),
		CashierId:     int32(cashierID),
		PaymentMethod: "E-Wallet",
		PaymentStatus: "pending",
		Amount:        100000,
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	transID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.client.FindById(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)
	s.Equal("E-Wallet", getRes.Data.PaymentMethod)

	// 4. FindAll
	allRes, err := s.client.FindAll(ctx, &pb.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pb.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.client.Update(ctx, &pb.UpdateTransactionRequest{
		TransactionId: transID,
		OrderId:       int32(orderID),
		CashierId:     int32(cashierID),
		PaymentMethod: "Credit Card",
		Amount:        150000,
	})
	s.Require().NoError(err)
	s.Equal("Credit Card", updateRes.Data.PaymentMethod)

	// 7. Trash
	_, err = s.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pb.FindAllTransactionRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	_, err = s.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{Id: transID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.client.RestoreAllTransaction(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestTransactionGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGapiTestSuite))
}
