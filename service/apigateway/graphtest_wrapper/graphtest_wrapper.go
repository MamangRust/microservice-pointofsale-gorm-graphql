package graphtest_wrapper

// Re-export graphtest so tests in external modules can access it.
// This works because graphtest_wrapper lives inside the apigateway module.

import (
	"context"
	"net/http"

	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type ServiceConnections = graphtest.ServiceConnections
type Resolver = graphtest.Resolver

func NewResolver(conns *ServiceConnections, log logger.LoggerInterface, redisClient *goredis.Client) *Resolver {
	return graphtest.NewResolver(conns, log, redisClient)
}

func NewHandler(resolver *Resolver) http.Handler {
	return graphtest.NewHandler(resolver)
}

func WithUserID(ctx context.Context, userID int) context.Context {
	return graphtest.WithUserID(ctx, userID)
}

func ExecuteGraphQL(srv http.Handler, query string, variables map[string]interface{}, authToken string) (*graphtest.GraphQLResponse, error) {
	return graphtest.ExecuteGraphQL(srv, query, variables, authToken)
}

// ConnMap is a helper to build ServiceConnections from a map of *grpc.ClientConn.
// Keys must match: auth, user, role, category, merchant, cashier, order, order-item, product, transaction.
func ConnMap(conns map[string]*grpc.ClientConn) *ServiceConnections {
	return &ServiceConnections{
		AuthClient:        conns["auth"],
		UserClient:        conns["user"],
		RoleClient:        conns["role"],
		CategoryClient:    conns["category"],
		MerchantClient:    conns["merchant"],
		CashierClient:     conns["cashier"],
		OrderClient:       conns["order"],
		OrderItemClient:   conns["order-item"],
		ProductClient:     conns["product"],
		TransactionClient: conns["transaction"],
		StatsReaderClient: conns["stats-reader"],
	}
}
