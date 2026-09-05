package order_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	handler   http.Handler
	orderID   int
	userID    int
	cashierID int
	merchID   int
	prodID    int
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()
	s.SetupCashierService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	s.merchID = s.SeedMerchant(ctx, s.userID)
	s.prodID = s.SeedProduct(ctx, s.merchID, catID)
	s.orderID = s.SeedOrder(ctx, s.userID, s.merchID, s.prodID)
	_ = s.GormDB().WithContext(ctx).Raw(
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES (?, ?, 'Order Api Cashier') RETURNING cashier_id`,
		s.merchID, s.userID,
	).Scan(&s.cashierID)

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findByIdOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllOrder(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveOrder(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedOrder(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedOrder(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteOrderPermanent(input: { id: `+strconv.Itoa(s.orderID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteOrderPermanent"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllOrder { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllOrder"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllOrderPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllOrderPermanent"].(map[string]interface{})["status"])
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
