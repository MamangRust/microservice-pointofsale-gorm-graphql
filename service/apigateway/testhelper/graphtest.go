package testhelper

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	mycontext "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/context"
	graph "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/handler"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/middlewares"
	mencache "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/redis"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceConnections mirrors graph.ServiceConnections for test use.
type ServiceConnections = graph.ServiceConnections

// CreateDummyConn creates a lazy gRPC connection that will never actually connect.
func CreateDummyConn() *grpc.ClientConn {
	conn, err := grpc.NewClient("localhost:1", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	return conn
}

// NewResolver creates a Resolver with the provided service connections.
func NewResolver(conns *ServiceConnections, log logger.LoggerInterface) *graph.Resolver {
	myMencache := mencache.NewCacheApiGateway(&mencache.Deps{
		Redis:  nil,
		Logger: log,
	})

	return graph.NewResolver(&graph.Deps{
		Clients:  conns,
		Logger:   log,
		Kafka:    nil,
		Mencache: myMencache,
	})
}

// NewResolverWithRedis creates a Resolver with the provided service connections and Redis client.
func NewResolverWithRedis(conns *ServiceConnections, log logger.LoggerInterface, redisClient *redis.Client) *graph.Resolver {
	myMencache := mencache.NewCacheApiGateway(&mencache.Deps{
		Redis:  redisClient,
		Logger: log,
	})

	return graph.NewResolver(&graph.Deps{
		Clients:  conns,
		Logger:   log,
		Kafka:    nil,
		Mencache: myMencache,
	})
}

// permissiveRoleChecker stands in for the RBAC checker of the running gateway.
//
// The harness mounts the schema without AuthMiddleware, and that middleware is
// what puts the authenticated user in the request context, so there is no user
// for the @hasRole directive to authorise and it passes every field through.
// The checker is wired anyway so the directive is never left nil, and so that a
// test which injects a user (see WithUser) gets a predictable answer instead of
// an "rbac: role checker is not configured" error.
type permissiveRoleChecker struct{}

func (permissiveRoleChecker) CheckRole(context.Context, int, ...string) error { return nil }

// NewGraphQLHTTPHandler creates an http.Handler from a gqlgen Resolver.
func NewGraphQLHTTPHandler(resolver *graph.Resolver) http.Handler {
	return NewGraphQLHTTPHandlerWithRoleChecker(resolver, permissiveRoleChecker{})
}

// NewGraphQLHTTPHandlerWithRoleChecker builds the handler with an explicit RBAC
// checker, so a test can exercise the @hasRole directive (pair it with WithUser
// to simulate an authenticated caller).
func NewGraphQLHTTPHandlerWithRoleChecker(resolver *graph.Resolver, checker middlewares.RoleChecker) http.Handler {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers:  resolver,
		Directives: graph.DirectiveRoot{HasRole: middlewares.HasRole(checker)},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})

	return srv
}

// NewTestHandler creates a full GraphQL handler from test service connections.
func NewTestHandler(conns *ServiceConnections, redisClient *redis.Client, log logger.LoggerInterface) http.Handler {
	resolver := NewResolverWithRedis(conns, log, redisClient)
	return NewGraphQLHTTPHandler(resolver)
}

// SetupServiceConnections creates ServiceConnections from test map.
func SetupServiceConnections(conns map[string]*grpc.ClientConn) *ServiceConnections {
	sc := &ServiceConnections{}

	if c, ok := conns["auth"]; ok {
		sc.AuthClient = c
	}
	if c, ok := conns["role"]; ok {
		sc.RoleClient = c
	}
	if c, ok := conns["user"]; ok {
		sc.UserClient = c
	}
	if c, ok := conns["category"]; ok {
		sc.CategoryClient = c
	}
	if c, ok := conns["cashier"]; ok {
		sc.CashierClient = c
	}
	if c, ok := conns["merchant"]; ok {
		sc.MerchantClient = c
	}
	if c, ok := conns["order"]; ok {
		sc.OrderClient = c
	}
	if c, ok := conns["order-item"]; ok {
		sc.OrderItemClient = c
	}
	if c, ok := conns["product"]; ok {
		sc.ProductClient = c
	}
	if c, ok := conns["transaction"]; ok {
		sc.TransactionClient = c
	}
	if c, ok := conns["stats_reader"]; ok {
		sc.StatsReaderClient = c
	}

	return sc
}

// WithUser injects the given user id into the request context, emulating what
// AuthMiddleware does for authenticated traffic.
func WithUser(next http.Handler, userID int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(mycontext.WithUserID(r.Context(), userID)))
	})
}

// SeedMerchantCache writes a merchant ID-to-API-key mapping into Redis
func SeedMerchantCache(redisClient *redis.Client, merchantID string, apiKey string) error {
	key := "merchant_api_key:" + merchantID
	return redisClient.Set(context.Background(), key, apiKey, 0).Err()
}

// SeedRoleCache writes a user role mapping into Redis
func SeedRoleCache(redisClient *redis.Client, userID string, roles []string) error {
	key := "user_roles:" + userID
	return redisClient.Set(context.Background(), key, roles, 0).Err()
}
