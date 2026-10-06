package graph

import (
	errorstd "errors"
	"fmt"
	"time"

	authgraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/auth"
	cashiergraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/cashier"
	categorygraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/category"
	merchantgraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/merchant"
	merchantdocumentgraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/merchant_document"
	ordergraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/order"
	orderitemgraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/order_item"
	productgraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/product"
	rolegraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/role"
	transactiongraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/transaction"
	usergraphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper/user"

	merchantpermission "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/permission/merchant"
	rolepermission "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/permission/role"
	mencache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis"

	auth_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/auth"
	cashier_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/cashier"
	category_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/category"
	merchant_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/merchant"
	merchant_document_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/merchant_document"
	order_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/order"
	orderitem_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/order_item"
	product_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/product"
	role_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/role"
	transaction_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/transaction"
	user_cache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis/api/user"

	authpb "github.com/MamangRust/microservice-point-of-sale-pb/auth"
	cashierpb "github.com/MamangRust/microservice-point-of-sale-pb/cashier"
	categorypb "github.com/MamangRust/microservice-point-of-sale-pb/category"
	merchantpb "github.com/MamangRust/microservice-point-of-sale-pb/merchant"
	merchantdocumentpb "github.com/MamangRust/microservice-point-of-sale-pb/merchant_document"
	orderpb "github.com/MamangRust/microservice-point-of-sale-pb/order"
	orderitempb "github.com/MamangRust/microservice-point-of-sale-pb/order_item"
	productpb "github.com/MamangRust/microservice-point-of-sale-pb/product"
	rolepb "github.com/MamangRust/microservice-point-of-sale-pb/role"
	statspb "github.com/MamangRust/microservice-point-of-sale-pb/stats"
	transactionpb "github.com/MamangRust/microservice-point-of-sale-pb/transaction"
	userpb "github.com/MamangRust/microservice-point-of-sale-pb/user"
	userrolepb "github.com/MamangRust/microservice-point-of-sale-pb/user_role"

	"github.com/MamangRust/microservice-point-of-sale-pkg/kafka"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-pkg/upload_image"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	sharedErrors "github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthGraphql             AuthHandleGraphql
	RoleGraphql             RoleHandleGraphql
	UserGraphql             UserHandleGraphql
	CashierGraphql          CashierHandleGraphql
	CategoryGraphql         CategoryHandleGraphql
	MerchantGraphql         MerchantHandleGraphql
	MerchantDocumentGraphql MerchantDocumentHandleGraphql
	OrderGraphql            OrderHandleGraphql
	OrderItemGraphql        OrderItemHandleGraphql
	ProductGraphql          ProductHandleGraphql
	TransactionGraphql      TransactionHandleGraphql
	ResolverHandle          *resolverHandler
}

type UserClient struct {
	userpb.UserQueryServiceClient
	userpb.UserCommandServiceClient
}

type RoleClient struct {
	rolepb.RoleQueryServiceClient
	rolepb.RoleCommandServiceClient
	UserRole userrolepb.UserRoleServiceClient
}

// The domain services still declare the legacy stats rpcs, so the stats-reader
// clients (main/ByMerchant/ById) live in named fields instead of being embedded
// — embedding both would make every stats method an ambiguous selector.
type CashierClient struct {
	cashierpb.CashierQueryServiceClient
	cashierpb.CashierCommandServiceClient
	Stats           statspb.CashierStatsServiceClient
	StatsByMerchant statspb.CashierStatsByMerchantServiceClient
	StatsById       statspb.CashierStatsByIdServiceClient
}

type CategoryClient struct {
	categorypb.CategoryQueryServiceClient
	categorypb.CategoryCommandServiceClient
	Stats           statspb.CategoryStatsServiceClient
	StatsByMerchant statspb.CategoryStatsByMerchantServiceClient
	StatsById       statspb.CategoryStatsByIdServiceClient
}

type MerchantClient struct {
	merchantpb.MerchantQueryServiceClient
	merchantpb.MerchantCommandServiceClient
}

type MerchantDocumentClient struct {
	merchantdocumentpb.MerchantDocumentServiceClient
}

type OrderClient struct {
	orderpb.OrderQueryServiceClient
	orderpb.OrderCommandServiceClient
	Stats           statspb.OrderStatsServiceClient
	StatsByMerchant statspb.OrderStatsByMerchantServiceClient
	StatsById       statspb.OrderStatsByIdServiceClient
}

type OrderItemClient struct {
	orderitempb.OrderItemQueryServiceClient
	orderitempb.OrderItemCommandServiceClient
}

type ProductClient struct {
	productpb.ProductQueryServiceClient
	productpb.ProductCommandServiceClient
}

type TransactionClient struct {
	transactionpb.TransactionQueryServiceClient
	transactionpb.TransactionCommandServiceClient
	StatsStatus statspb.TransactionStatsStatusServiceClient
	StatsMethod statspb.TransactionStatsMethodServiceClient
}

type AuthHandleGraphql struct {
	AuthClient authpb.AuthServiceClient
	Logger     logger.LoggerInterface
	Mapping    authgraphqlmapper.AuthGraphqlMapper
	Cache      auth_cache.AuthMencache
}

type RoleHandleGraphql struct {
	RoleClient RoleClient
	Logger     logger.LoggerInterface
	Mapping    rolegraphqlmapper.RoleGraphqlMapper
	Kafka      *kafka.Kafka
	Permission rolepermission.RolePermission
	Cache      role_cache.RoleMencache
}

type UserHandleGraphql struct {
	UserClient UserClient
	Logger     logger.LoggerInterface
	Mapping    usergraphqlmapper.UserGraphqlMapper
	Cache      user_cache.UserMencache
}

type CashierHandleGraphql struct {
	CashierClient CashierClient
	Logger        logger.LoggerInterface
	Mapping       cashiergraphqlmapper.CashierGraphqlMapper
	Cache         cashier_cache.CashierMencache
}

type CategoryHandleGraphql struct {
	CategoryClient CategoryClient
	Logger         logger.LoggerInterface
	Mapping        categorygraphqlmapper.CategoryGraphqlMapper
	Cache          category_cache.CategoryMencache
}

type MerchantHandleGraphql struct {
	MerchantClient MerchantClient
	Logger         logger.LoggerInterface
	Mapping        merchantgraphqlmapper.MerchantGraphqlMapper
	Cache          merchant_cache.MerchantMenCache
}

type MerchantDocumentHandleGraphql struct {
	MerchantClient MerchantDocumentClient
	Logger         logger.LoggerInterface
	Mapping        merchantdocumentgraphqlmapper.MerchantDocumentGraphqlMapper
	Cache          merchant_document_cache.MerchantDocumentMencache
}

type OrderHandleGraphql struct {
	OrderClient OrderClient
	Logger      logger.LoggerInterface
	Mapping     ordergraphqlmapper.OrderGraphqlMapper
	Cache       order_cache.OrderMencache
}

type OrderItemHandleGraphql struct {
	OrderItemClient OrderItemClient
	Logger          logger.LoggerInterface
	Mapping         orderitemgraphqlmapper.OrderItemGraphqlMapper
	Cache           orderitem_cache.OrderItemCache
}

type ProductHandleGraphql struct {
	ProductClient ProductClient
	Logger        logger.LoggerInterface
	Mapping       productgraphqlmapper.ProductGraphqlMapper
	Cache         product_cache.ProductMencache
	ImageUpload   upload_image.ImageUploads
}

type TransactionHandleGraphql struct {
	TransactionClient TransactionClient
	Logger            logger.LoggerInterface
	Mapping           transactiongraphqlmapper.TransactionGraphqlMapper
	Permission        merchantpermission.MerchantPermission
	Cache             transaction_cache.TransactionMencache
}

type ServiceConnections struct {
	AuthClient        *grpc.ClientConn
	CashierClient     *grpc.ClientConn
	CategoryClient    *grpc.ClientConn
	MerchantClient    *grpc.ClientConn
	OrderClient       *grpc.ClientConn
	OrderItemClient   *grpc.ClientConn
	ProductClient     *grpc.ClientConn
	RoleClient        *grpc.ClientConn
	StatsReaderClient *grpc.ClientConn
	TransactionClient *grpc.ClientConn
	UserClient        *grpc.ClientConn
}

type Deps struct {
	Clients  *ServiceConnections
	Logger   logger.LoggerInterface
	Kafka    *kafka.Kafka
	Mencache mencache.CacheApiGateway
}

func NewResolver(
	deps *Deps,
) *Resolver {
	observability, _ := observability.NewObservability(
		"graphql-client",
		deps.Logger,
	)

	resolverHandle := NewResolverHandler(observability, deps.Logger)

	store := deps.Mencache.GetStore()
	cacheAuth := auth_cache.NewMencache(store)
	cacheUser := user_cache.NewUserMencache(store)
	cacheRole := role_cache.NewRoleMencache(store)
	cacheMerchant := merchant_cache.NewMerchantMencache(store)
	cacheMerchantDocument := merchant_document_cache.NewMerchantDocumentMencache(store)
	cacheCashier := cashier_cache.NewCashierMencache(store)
	cacheCategory := category_cache.NewCategoryMencache(store)
	cacheOrder := order_cache.NewOrderMencache(store)
	cacheOrderItem := orderitem_cache.NewOrderItemCache(store)
	cacheProduct := product_cache.NewProductMencache(store)
	cacheTransaction := transaction_cache.NewTransactionMencache(store)

	// RBAC for the @hasRole directive: roles are resolved through the Kafka
	// request-role/response-role handshake against the role service's consumer
	// and cached in Redis, matching the payment gateway's RolePermission.
	rolePermission := rolepermission.NewRolePermission(
		deps.Kafka,
		"request-role",
		"response-role",
		5*time.Second,
		deps.Logger,
		deps.Mencache,
	)

	return &Resolver{
		ResolverHandle: resolverHandle,
		AuthGraphql: AuthHandleGraphql{
			AuthClient: authpb.NewAuthServiceClient(deps.Clients.AuthClient),
			Logger:     deps.Logger,
			Mapping:    authgraphqlmapper.NewAuthGraphqlMapper(),
			Cache:      cacheAuth,
		},
		RoleGraphql: RoleHandleGraphql{RoleClient: RoleClient{
			RoleQueryServiceClient:   rolepb.NewRoleQueryServiceClient(deps.Clients.RoleClient),
			RoleCommandServiceClient: rolepb.NewRoleCommandServiceClient(deps.Clients.RoleClient),
			UserRole:                 userrolepb.NewUserRoleServiceClient(deps.Clients.RoleClient),
		},
			Kafka:      deps.Kafka,
			Logger:     deps.Logger,
			Mapping:    rolegraphqlmapper.NewRoleGraphqlMapper(),
			Permission: rolePermission,
			Cache:      cacheRole,
		},
		UserGraphql: UserHandleGraphql{UserClient: UserClient{
			UserQueryServiceClient:   userpb.NewUserQueryServiceClient(deps.Clients.UserClient),
			UserCommandServiceClient: userpb.NewUserCommandServiceClient(deps.Clients.UserClient),
		},
			Logger:  deps.Logger,
			Mapping: usergraphqlmapper.NewUserGraphqlMapper(),
			Cache:   cacheUser,
		},
		CashierGraphql: CashierHandleGraphql{CashierClient: CashierClient{
			CashierQueryServiceClient:   cashierpb.NewCashierQueryServiceClient(deps.Clients.CashierClient),
			CashierCommandServiceClient: cashierpb.NewCashierCommandServiceClient(deps.Clients.CashierClient),
			Stats:                       statspb.NewCashierStatsServiceClient(deps.Clients.StatsReaderClient),
			StatsByMerchant:             statspb.NewCashierStatsByMerchantServiceClient(deps.Clients.StatsReaderClient),
			StatsById:                   statspb.NewCashierStatsByIdServiceClient(deps.Clients.StatsReaderClient),
		},
			Logger:  deps.Logger,
			Mapping: cashiergraphqlmapper.NewCashierGraphqlMapper(),
			Cache:   cacheCashier,
		},
		CategoryGraphql: CategoryHandleGraphql{CategoryClient: CategoryClient{
			CategoryQueryServiceClient:   categorypb.NewCategoryQueryServiceClient(deps.Clients.CategoryClient),
			CategoryCommandServiceClient: categorypb.NewCategoryCommandServiceClient(deps.Clients.CategoryClient),
			Stats:                        statspb.NewCategoryStatsServiceClient(deps.Clients.StatsReaderClient),
			StatsByMerchant:              statspb.NewCategoryStatsByMerchantServiceClient(deps.Clients.StatsReaderClient),
			StatsById:                    statspb.NewCategoryStatsByIdServiceClient(deps.Clients.StatsReaderClient),
		},
			Logger:  deps.Logger,
			Mapping: categorygraphqlmapper.NewCategoryGraphqlMapper(),
			Cache:   cacheCategory,
		},
		MerchantGraphql: MerchantHandleGraphql{MerchantClient: MerchantClient{
			MerchantQueryServiceClient:   merchantpb.NewMerchantQueryServiceClient(deps.Clients.MerchantClient),
			MerchantCommandServiceClient: merchantpb.NewMerchantCommandServiceClient(deps.Clients.MerchantClient),
		},
			Logger:  deps.Logger,
			Mapping: merchantgraphqlmapper.NewMerchantGraphqlMapper(),
			Cache:   cacheMerchant,
		},
		MerchantDocumentGraphql: MerchantDocumentHandleGraphql{MerchantClient: MerchantDocumentClient{merchantdocumentpb.NewMerchantDocumentServiceClient(deps.Clients.MerchantClient)},
			Logger:  deps.Logger,
			Mapping: merchantdocumentgraphqlmapper.NewMerchantDocumentGraphqlMapper(),
			Cache:   cacheMerchantDocument,
		},
		OrderGraphql: OrderHandleGraphql{OrderClient: OrderClient{
			OrderQueryServiceClient:   orderpb.NewOrderQueryServiceClient(deps.Clients.OrderClient),
			OrderCommandServiceClient: orderpb.NewOrderCommandServiceClient(deps.Clients.OrderClient),
			Stats:                     statspb.NewOrderStatsServiceClient(deps.Clients.StatsReaderClient),
			StatsByMerchant:           statspb.NewOrderStatsByMerchantServiceClient(deps.Clients.StatsReaderClient),
			StatsById:                 statspb.NewOrderStatsByIdServiceClient(deps.Clients.StatsReaderClient),
		},
			Logger:  deps.Logger,
			Mapping: ordergraphqlmapper.NewOrderGraphqlMapper(),
			Cache:   cacheOrder,
		},
		OrderItemGraphql: OrderItemHandleGraphql{OrderItemClient: OrderItemClient{
			OrderItemQueryServiceClient:   orderitempb.NewOrderItemQueryServiceClient(deps.Clients.OrderItemClient),
			OrderItemCommandServiceClient: orderitempb.NewOrderItemCommandServiceClient(deps.Clients.OrderItemClient),
		},
			Logger:  deps.Logger,
			Mapping: orderitemgraphqlmapper.NewOrderItemGraphqlMapper(),
			Cache:   cacheOrderItem,
		},
		ProductGraphql: ProductHandleGraphql{ProductClient: ProductClient{
			ProductQueryServiceClient:   productpb.NewProductQueryServiceClient(deps.Clients.ProductClient),
			ProductCommandServiceClient: productpb.NewProductCommandServiceClient(deps.Clients.ProductClient),
		},
			Logger:      deps.Logger,
			Mapping:     productgraphqlmapper.NewProductGraphqlMapper(),
			Cache:       cacheProduct,
			ImageUpload: upload_image.NewImageUpload(deps.Logger),
		},
		TransactionGraphql: TransactionHandleGraphql{TransactionClient: TransactionClient{
			TransactionQueryServiceClient:   transactionpb.NewTransactionQueryServiceClient(deps.Clients.TransactionClient),
			TransactionCommandServiceClient: transactionpb.NewTransactionCommandServiceClient(deps.Clients.TransactionClient),
			StatsStatus:                     statspb.NewTransactionStatsStatusServiceClient(deps.Clients.StatsReaderClient),
			StatsMethod:                     statspb.NewTransactionStatsMethodServiceClient(deps.Clients.StatsReaderClient),
		},
			Logger:     deps.Logger,
			Mapping:    transactiongraphqlmapper.NewTransactionGraphqlMapper(),
			Permission: merchantpermission.NewMerchantPermission(deps.Kafka, "request-transaction", "response-transaction", 5*time.Second, deps.Logger),
			Cache:      cacheTransaction,
		},
	}
}

func (h *Resolver) handleGraphQLError(err error, operation string) *errors.AppError {
	if err == nil {
		return nil
	}

	var appErr *errors.AppError
	if errorstd.As(err, &appErr) {
		return appErr
	}

	return errors.NewInternalError(err).WithMessage("Failed to " + operation)
}

func (r *Resolver) parseValidationErrors(err error) []sharedErrors.ValidationError {
	var validationErrs []sharedErrors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, sharedErrors.ValidationError{
				Field:   fe.Field(),
				Message: r.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []sharedErrors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (r *Resolver) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
