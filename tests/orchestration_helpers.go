package tests

import (
	"bytes"
	"context"
	pb "github.com/MamangRust/microservice-point-of-sale-pb/auth"
	cashier "github.com/MamangRust/microservice-point-of-sale-pb/cashier"
	category "github.com/MamangRust/microservice-point-of-sale-pb/category"
	merchant "github.com/MamangRust/microservice-point-of-sale-pb/merchant"
	merchant_document "github.com/MamangRust/microservice-point-of-sale-pb/merchant_document"
	order "github.com/MamangRust/microservice-point-of-sale-pb/order"
	order_item "github.com/MamangRust/microservice-point-of-sale-pb/order_item"
	product "github.com/MamangRust/microservice-point-of-sale-pb/product"
	role "github.com/MamangRust/microservice-point-of-sale-pb/role"
	statspb "github.com/MamangRust/microservice-point-of-sale-pb/stats"
	transaction "github.com/MamangRust/microservice-point-of-sale-pb/transaction"
	user "github.com/MamangRust/microservice-point-of-sale-pb/user"
	user_role "github.com/MamangRust/microservice-point-of-sale-pb/user_role"
	"mime/multipart"

	"fmt"
	"github.com/MamangRust/microservice-point-of-sale-pkg/auth"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"

	clickhouseDB "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/MamangRust/microservice-point-of-sale-pkg/clickhouse"
	stats_handler "github.com/MamangRust/microservice-point-of-sale-stats-reader/handler"
	stats_repo "github.com/MamangRust/microservice-point-of-sale-stats-reader/repository"
	clickhouseTC "github.com/testcontainers/testcontainers-go/modules/clickhouse"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/microservice-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/microservice-point-of-sale-role/handler"
	role_repo "github.com/MamangRust/microservice-point-of-sale-role/repository"
	role_service "github.com/MamangRust/microservice-point-of-sale-role/service"

	// User
	user_cache "github.com/MamangRust/microservice-point-of-sale-user/cache"
	user_handler "github.com/MamangRust/microservice-point-of-sale-user/handler"
	user_repo "github.com/MamangRust/microservice-point-of-sale-user/repository"
	user_service "github.com/MamangRust/microservice-point-of-sale-user/service"

	// Auth
	auth_cache "github.com/MamangRust/microservice-point-of-sale-auth/cache"
	auth_handler "github.com/MamangRust/microservice-point-of-sale-auth/handler"
	auth_repo "github.com/MamangRust/microservice-point-of-sale-auth/repository"
	auth_service "github.com/MamangRust/microservice-point-of-sale-auth/service"

	// Category
	category_cache "github.com/MamangRust/microservice-point-of-sale-category/cache"
	category_handler "github.com/MamangRust/microservice-point-of-sale-category/handler"
	category_repo "github.com/MamangRust/microservice-point-of-sale-category/repository"
	category_service "github.com/MamangRust/microservice-point-of-sale-category/service"

	// Product
	product_cache "github.com/MamangRust/microservice-point-of-sale-product/cache"
	product_handler "github.com/MamangRust/microservice-point-of-sale-product/handler"
	product_repo "github.com/MamangRust/microservice-point-of-sale-product/repository"
	product_service "github.com/MamangRust/microservice-point-of-sale-product/service"

	// Merchant
	merchant_cache "github.com/MamangRust/microservice-point-of-sale-merchant/cache"
	merchant_handler "github.com/MamangRust/microservice-point-of-sale-merchant/handler"
	merchant_repo "github.com/MamangRust/microservice-point-of-sale-merchant/repository"
	merchant_service "github.com/MamangRust/microservice-point-of-sale-merchant/service"

	// Order
	order_cache "github.com/MamangRust/microservice-point-of-sale-order/cache"
	order_handler "github.com/MamangRust/microservice-point-of-sale-order/handler"
	order_repo "github.com/MamangRust/microservice-point-of-sale-order/repository"
	order_service "github.com/MamangRust/microservice-point-of-sale-order/service"

	// Transaction
	transaction_cache "github.com/MamangRust/microservice-point-of-sale-transacton/cache"
	transaction_handler "github.com/MamangRust/microservice-point-of-sale-transacton/handler"
	transaction_repo "github.com/MamangRust/microservice-point-of-sale-transacton/repository"
	transaction_service "github.com/MamangRust/microservice-point-of-sale-transacton/service"

	// Order Item
	order_item_cache "github.com/MamangRust/microservice-point-of-sale-order-item/cache"
	order_item_handler "github.com/MamangRust/microservice-point-of-sale-order-item/handler"
	order_item_repo "github.com/MamangRust/microservice-point-of-sale-order-item/repository"
	order_item_service "github.com/MamangRust/microservice-point-of-sale-order-item/service"

	// Cashier
	mencache "github.com/MamangRust/microservice-point-of-sale-cashier/cache"
	cashier_handler "github.com/MamangRust/microservice-point-of-sale-cashier/handler"
	cashier_repo "github.com/MamangRust/microservice-point-of-sale-cashier/repository"
	cashier_service "github.com/MamangRust/microservice-point-of-sale-cashier/service"
)

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(gormDB)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        s.Log,
		Mencache:      roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(roleSvc)
	server := grpc.NewServer()
	role.RegisterRoleQueryServiceServer(server, roleGapi)
	role.RegisterRoleCommandServiceServer(server, roleGapi)
	user_role.RegisterUserRoleServiceServer(server, roleGapi)
	addr, err := RunGRPCServer(server)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupUserService() {
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	roleClient := role.NewRoleQueryServiceClient(s.Conns["role"])
	userRoleClient := user_role.NewUserRoleServiceClient(s.Conns["role"])
	userRepos := user_repo.NewRepositories(gormDB, roleClient, userRoleClient)
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(userSvc)
	server := grpc.NewServer()
	user.RegisterUserQueryServiceServer(server, userGapi)
	user.RegisterUserCommandServiceServer(server, userGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	userQueryClient := user.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := user.NewUserCommandServiceClient(s.Conns["user"])
	roleClient := role.NewRoleQueryServiceClient(s.Conns["role"])
	userRoleClient := user_role.NewUserRoleServiceClient(s.Conns["role"])

	authRepos := auth_repo.NewRepositories(gormDB, userQueryClient, userCommandClient, roleClient, userRoleClient)
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(gormDB)
	catSvc := category_service.NewService(&category_service.Deps{
		Mencache:      catMencache,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(catSvc)
	server := grpc.NewServer()
	category.RegisterCategoryQueryServiceServer(server, catGapi)
	category.RegisterCategoryCommandServiceServer(server, catGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	if _, ok := s.Conns["category"]; !ok {
		s.SetupCategoryService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	prodMencache := product_cache.NewMencache(cacheStore)
	categoryClient := category.NewCategoryQueryServiceClient(s.Conns["category"])
	merchantClient := merchant.NewMerchantQueryServiceClient(s.Conns["merchant"])
	prodRepos := product_repo.NewRepositories(gormDB, categoryClient, merchantClient)
	prodSvc := product_service.NewService(&product_service.Deps{
		Mencache:      prodMencache,
		Repositories:  prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(prodSvc)
	server := grpc.NewServer()
	product.RegisterProductQueryServiceServer(server, prodGapi)
	product.RegisterProductCommandServiceServer(server, prodGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userQueryClient := user.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := user.NewUserCommandServiceClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(gormDB, userQueryClient, userCommandClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi, merchantDocGapi := merchant_handler.NewHandler(merchantSvc)
	server := grpc.NewServer()
	merchant.RegisterMerchantQueryServiceServer(server, merchantGapi)
	merchant.RegisterMerchantCommandServiceServer(server, merchantGapi)
	merchant_document.RegisterMerchantDocumentServiceServer(server, merchantDocGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["product"]; !ok {
		s.SetupProductService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	cashierClient := cashier.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := merchant.NewMerchantQueryServiceClient(s.Conns["merchant"])
	productQueryClient := product.NewProductQueryServiceClient(s.Conns["product"])
	productCommandClient := product.NewProductCommandServiceClient(s.Conns["product"])
	orderItemQueryClient := order_item.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	orderItemCommandClient := order_item.NewOrderItemCommandServiceClient(s.Conns["order-item"])

	orderMencache := order_cache.NewMencache(cacheStore)
	orderRepos := order_repo.NewRepositories(gormDB, cashierClient, merchantClient, productQueryClient, productCommandClient, orderItemQueryClient, orderItemCommandClient)
	orderSvc := order_service.NewService(&order_service.Deps{
		Mencache:      orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(orderSvc)
	server := grpc.NewServer()
	order.RegisterOrderQueryServiceServer(server, orderGapi)
	order.RegisterOrderCommandServiceServer(server, orderGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["order"]; !ok {
		s.SetupOrderService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	cashierClient := cashier.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := merchant.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := order.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemQueryClient := order_item.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	orderItemCommandClient := order_item.NewOrderItemCommandServiceClient(s.Conns["order-item"])

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	transactionRepos := transaction_repo.NewRepositories(gormDB, cashierClient, merchantClient, orderClient, orderItemQueryClient, orderItemCommandClient)
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Mencache:      transactionMencache,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(transactionSvc, s.Log)
	server := grpc.NewServer()
	transaction.RegisterTransactionQueryServiceServer(server, transactionGapi)
	transaction.RegisterTransactionCommandServiceServer(server, transactionGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(gormDB)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Mencache:      itemMencache,
		Repositories:  itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(itemSvc, s.Log)
	server := grpc.NewServer()
	order_item.RegisterOrderItemQueryServiceServer(server, itemGapi)
	order_item.RegisterOrderItemCommandServiceServer(server, itemGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCashierService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	gormDB := s.GormDB()

	cashierMencache := mencache.NewMencache(cacheStore)
	userQueryClient := user.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := user.NewUserCommandServiceClient(s.Conns["user"])
	merchantClient := merchant.NewMerchantQueryServiceClient(s.Conns["merchant"])
	cashierRepos := cashier_repo.NewRepositories(gormDB, userQueryClient, userCommandClient, merchantClient)
	cashierSvc := cashier_service.NewService(&cashier_service.Deps{
		Ctx:           context.Background(),
		Mencache:      cashierMencache,
		Repositories:  cashierRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cashierGapi := cashier_handler.NewHandler(cashierSvc)
	server := grpc.NewServer()
	cashier.RegisterCashierQueryServiceServer(server, cashierGapi)
	cashier.RegisterCashierCommandServiceServer(server, cashierGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cashier"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}

func (s *BaseTestSuite) SetupStatsReaderService() {
	// Start ClickHouse testcontainer
	chCtx := context.Background()
	chContainer, err := clickhouseTC.RunContainer(chCtx,
		clickhouseTC.WithDatabase("pos_stats"),
		clickhouseTC.WithPassword("test_password"),
	)
	if err != nil {
		s.T().Skipf("Skipping stats reader tests: ClickHouse container failed to start: %v", err)
	}

	// Get the mapped host/port
	chHost, err := chContainer.Host(chCtx)
	if err != nil {
		s.T().Skipf("Skipping stats reader tests: failed to get ClickHouse host: %v", err)
	}
	chPort, err := chContainer.MappedPort(chCtx, "9000")
	if err != nil {
		s.T().Skipf("Skipping stats reader tests: failed to get ClickHouse port: %v", err)
	}

	// Connect to ClickHouse
	connStr := fmt.Sprintf("%s:%s", chHost, chPort.Port())
	chConn, err := clickhouseDB.Open(&clickhouseDB.Options{
		Addr: []string{connStr},
		Auth: clickhouseDB.Auth{
			Database: "pos_stats",
			Username: "default",
			Password: "test_password",
		},
	})
	if err != nil {
		s.T().Skipf("Skipping stats reader tests: failed to connect to ClickHouse: %v", err)
	}

	// Apply schema
	if err := clickhouse.ApplySchema(chCtx, chConn, s.Log); err != nil {
		s.T().Skipf("Skipping stats reader tests: failed to apply schema: %v", err)
	}

	// Create stats-reader repository + handlers (single ClickHouse reader
	// implementing the full stats contract)
	repo := stats_repo.NewClickHouseReaderRepository(chConn)
	var statsCache *stats_handler.StatsCache // no cache in tests

	server := grpc.NewServer()
	statspb.RegisterOrderStatsServiceServer(server, stats_handler.NewOrderStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterOrderStatsByMerchantServiceServer(server, stats_handler.NewOrderStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterOrderStatsByIdServiceServer(server, stats_handler.NewOrderStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCashierStatsServiceServer(server, stats_handler.NewCashierStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCashierStatsByMerchantServiceServer(server, stats_handler.NewCashierStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCashierStatsByIdServiceServer(server, stats_handler.NewCashierStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCategoryStatsServiceServer(server, stats_handler.NewCategoryStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCategoryStatsByMerchantServiceServer(server, stats_handler.NewCategoryStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterCategoryStatsByIdServiceServer(server, stats_handler.NewCategoryStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterTransactionStatsStatusServiceServer(server, stats_handler.NewTransactionStatsHandler(repo, statsCache, s.Log))
	statspb.RegisterTransactionStatsMethodServiceServer(server, stats_handler.NewTransactionStatsHandler(repo, statsCache, s.Log))

	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["stats-reader"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}
