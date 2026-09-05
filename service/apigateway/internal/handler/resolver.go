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

	pbauth "github.com/MamangRust/microservice-pointofsale-grpc/pb"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	pbcategory "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	pbproduct "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	pbrole "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	pbtransaction "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
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
	StatsRead               *StatsReadHandleGraphql
	ResolverHandle          *resolverHandler}

type UserClient struct {
	pbuser.UserServiceClient
}

type RoleClient struct {
	pbrole.RoleServiceClient
}

type CashierClient struct {
	pbcashier.CashierServiceClient
}

type CategoryClient struct {
	pbcategory.CategoryServiceClient
}

type MerchantClient struct {
	pbmerchant.MerchantServiceClient
}

type MerchantDocumentClient struct {
	pbmerchant.MerchantDocumentServiceClient
}

type OrderClient struct {
	pborder.OrderServiceClient
}

type OrderItemClient struct {
	pborderitem.OrderItemServiceClient
}

type ProductClient struct {
	pbproduct.ProductServiceClient
}

type TransactionClient struct {
	pbtransaction.TransactionServiceClient
}

type AuthHandleGraphql struct {
	AuthClient pbauth.AuthServiceClient
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

	newAuth := func(c *grpc.ClientConn) pbauth.AuthServiceClient {
		if c == nil {
			return nil
		}
		return pbauth.NewAuthServiceClient(c)
	}
	newUser := func(c *grpc.ClientConn) pbuser.UserServiceClient {
		if c == nil {
			return nil
		}
		return pbuser.NewUserServiceClient(c)
	}
	newRole := func(c *grpc.ClientConn) pbrole.RoleServiceClient {
		if c == nil {
			return nil
		}
		return pbrole.NewRoleServiceClient(c)
	}
	newMerchant := func(c *grpc.ClientConn) pbmerchant.MerchantServiceClient {
		if c == nil {
			return nil
		}
		return pbmerchant.NewMerchantServiceClient(c)
	}
	newMerchantDoc := func(c *grpc.ClientConn) pbmerchant.MerchantDocumentServiceClient {
		if c == nil {
			return nil
		}
		return pbmerchant.NewMerchantDocumentServiceClient(c)
	}
	newCategory := func(c *grpc.ClientConn) pbcategory.CategoryServiceClient {
		if c == nil {
			return nil
		}
		return pbcategory.NewCategoryServiceClient(c)
	}
	newCashier := func(c *grpc.ClientConn) pbcashier.CashierServiceClient {
		if c == nil {
			return nil
		}
		return pbcashier.NewCashierServiceClient(c)
	}
	newOrder := func(c *grpc.ClientConn) pborder.OrderServiceClient {
		if c == nil {
			return nil
		}
		return pborder.NewOrderServiceClient(c)
	}
	newOrderItem := func(c *grpc.ClientConn) pborderitem.OrderItemServiceClient {
		if c == nil {
			return nil
		}
		return pborderitem.NewOrderItemServiceClient(c)
	}
	newProduct := func(c *grpc.ClientConn) pbproduct.ProductServiceClient {
		if c == nil {
			return nil
		}
		return pbproduct.NewProductServiceClient(c)
	}
	newTransaction := func(c *grpc.ClientConn) pbtransaction.TransactionServiceClient {
		if c == nil {
			return nil
		}
		return pbtransaction.NewTransactionServiceClient(c)
	}
	newCategoryStats := func(c *grpc.ClientConn) pbstats.CategoryStatsServiceClient {
		if c == nil {
			return nil
		}
		return pbstats.NewCategoryStatsServiceClient(c)
	}
	newOrderStats := func(c *grpc.ClientConn) pbstats.OrderStatsServiceClient {
		if c == nil {
			return nil
		}
		return pbstats.NewOrderStatsServiceClient(c)
	}
	newTransactionStats := func(c *grpc.ClientConn) pbstats.TransactionStatsServiceClient {
		if c == nil {
			return nil
		}
		return pbstats.NewTransactionStatsServiceClient(c)
	}
	newCashierStats := func(c *grpc.ClientConn) pbstats.CashierStatsServiceClient {
		if c == nil {
			return nil
		}
		return pbstats.NewCashierStatsServiceClient(c)
	}
	newProductStats := func(c *grpc.ClientConn) pbstats.ProductStatsServiceClient {
		if c == nil {
			return nil
		}
		return pbstats.NewProductStatsServiceClient(c)
	}

	return &Resolver{
		ResolverHandle: resolverHandle,
		AuthGraphql: AuthHandleGraphql{
			AuthClient: newAuth(deps.Clients.AuthClient),
			Logger:     deps.Logger,
			Mapping:    authgraphqlmapper.NewAuthGraphqlMapper(),
			Cache:      cacheAuth,
		},
		RoleGraphql: RoleHandleGraphql{
			RoleClient: RoleClient{newRole(deps.Clients.RoleClient)},
			Kafka:      deps.Kafka,
			Logger:     deps.Logger,
			Mapping:    rolegraphqlmapper.NewRoleGraphqlMapper(),
			Permission: rolepermission.NewRolePermission(deps.Kafka, "request-role", "response-role", 5*time.Second, deps.Logger, deps.Mencache),
			Cache:      cacheRole,
		},
		UserGraphql: UserHandleGraphql{
			UserClient: UserClient{newUser(deps.Clients.UserClient)},
			Logger:     deps.Logger,
			Mapping:    usergraphqlmapper.NewUserGraphqlMapper(),
			Cache:      cacheUser,
		},
		CashierGraphql: CashierHandleGraphql{
			CashierClient: CashierClient{newCashier(deps.Clients.CashierClient)},
			Logger:        deps.Logger,
			Mapping:       cashiergraphqlmapper.NewCashierGraphqlMapper(),
			Cache:         cacheCashier,
		},
		CategoryGraphql: CategoryHandleGraphql{
			CategoryClient: CategoryClient{newCategory(deps.Clients.CategoryClient)},
			Logger:         deps.Logger,
			Mapping:        categorygraphqlmapper.NewCategoryGraphqlMapper(),
			Cache:          cacheCategory,
		},
		MerchantGraphql: MerchantHandleGraphql{
			MerchantClient: MerchantClient{newMerchant(deps.Clients.MerchantClient)},
			Logger:         deps.Logger,
			Mapping:        merchantgraphqlmapper.NewMerchantGraphqlMapper(),
			Cache:          cacheMerchant,
		},
		MerchantDocumentGraphql: MerchantDocumentHandleGraphql{
			MerchantClient: MerchantDocumentClient{newMerchantDoc(deps.Clients.MerchantClient)},
			Logger:         deps.Logger,
			Mapping:        merchantdocumentgraphqlmapper.NewMerchantDocumentGraphqlMapper(),
			Cache:          cacheMerchantDocument,
		},
		OrderGraphql: OrderHandleGraphql{
			OrderClient: OrderClient{newOrder(deps.Clients.OrderClient)},
			Logger:      deps.Logger,
			Mapping:     ordergraphqlmapper.NewOrderGraphqlMapper(),
			Cache:       cacheOrder,
		},
		OrderItemGraphql: OrderItemHandleGraphql{
			OrderItemClient: OrderItemClient{newOrderItem(deps.Clients.OrderItemClient)},
			Logger:          deps.Logger,
			Mapping:         orderitemgraphqlmapper.NewOrderItemGraphqlMapper(),
			Cache:           cacheOrderItem,
		},
		ProductGraphql: ProductHandleGraphql{
			ProductClient: ProductClient{newProduct(deps.Clients.ProductClient)},
			Logger:        deps.Logger,
			Mapping:       productgraphqlmapper.NewProductGraphqlMapper(),
			Cache:         cacheProduct,
			ImageUpload:   upload_image.NewImageUpload(deps.Logger),
		},
		TransactionGraphql: TransactionHandleGraphql{
			TransactionClient: TransactionClient{newTransaction(deps.Clients.TransactionClient)},
			Logger:            deps.Logger,
			Mapping:           transactiongraphqlmapper.NewTransactionGraphqlMapper(),
			Permission:        merchantpermission.NewMerchantPermission(deps.Kafka, "request-transaction", "response-transaction", 5*time.Second, deps.Logger),
			Cache:             cacheTransaction,
		},
		StatsRead: &StatsReadHandleGraphql{
			CategoryStats:   newCategoryStats(deps.Clients.StatsReaderClient),
			OrderStats:      newOrderStats(deps.Clients.StatsReaderClient),
			TransactionStats: newTransactionStats(deps.Clients.StatsReaderClient),
			CashierStats:    newCashierStats(deps.Clients.StatsReaderClient),
			ProductStats:    newProductStats(deps.Clients.StatsReaderClient),
		},
	}
}

type StatsReadHandleGraphql struct {
	CategoryStats    pbstats.CategoryStatsServiceClient
	OrderStats       pbstats.OrderStatsServiceClient
	TransactionStats pbstats.TransactionStatsServiceClient
	CashierStats     pbstats.CashierStatsServiceClient
	ProductStats     pbstats.ProductStatsServiceClient
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
