package order_item_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	orderID int
}

func (s *OrderItemApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	resolver := graphtest.NewResolver(&graphtest.ServiceConnections{
		AuthClient:        s.Conns["auth"],
		CashierClient:     s.Conns["cashier"],
		CategoryClient:    s.Conns["category"],
		MerchantClient:    s.Conns["merchant"],
		OrderClient:       s.Conns["order"],
		OrderItemClient:   s.Conns["order-item"],
		ProductClient:     s.Conns["product"],
		RoleClient:        s.Conns["role"],
		StatsReaderClient: s.Conns["stats-reader"],
		TransactionClient: s.Conns["transaction"],
		UserClient:        s.Conns["user"],
	}, s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *OrderItemApiTestSuite) TestOrderItemApiLifecycle() {
	// 1. FindAll
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findAllOrderItem(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllOrderItem"].(map[string]interface{})["status"])

	// 2. FindOrderItemByOrder
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findOrderItemByOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findOrderItemByOrder"].(map[string]interface{})["status"])

	// 3. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveOrderItem(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveOrderItem"].(map[string]interface{})["status"])

	// 4. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedOrderItem(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedOrderItem"].(map[string]interface{})["status"])
}

func TestOrderItemApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemApiTestSuite))
}
