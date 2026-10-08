package order_item_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	orderID int
}

func (s *OrderItemGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)
	s.orderID = s.SeedOrder(ctx, userID, merchID, prodID)

	s.handler = s.NewGraphQLHandler()
}

const findAllOrderItemsQuery = `query FindAllOrderItems($input: FindAllOrderItemInput!) {
  findAllOrderItem(input: $input) { status message data { id order_id product_id } pagination { total_records } }
}`

const findOrderItemsByOrderQuery = `query FindOrderItemsByOrder($input: FindByIdOrderItemInput!) {
  findOrderItemByOrder(input: $input) { status message data { id order_id product_id quantity } }
}`

const findActiveOrderItemsQuery = `query FindActiveOrderItems($input: FindAllOrderItemInput!) {
  findByActiveOrderItem(input: $input) { status message data { id } }
}`

const findTrashedOrderItemsQuery = `query FindTrashedOrderItems($input: FindAllOrderItemInput!) {
  findByTrashedOrderItem(input: $input) { status message data { id } }
}`

func (s *OrderItemGraphQLTestSuite) TestOrderItemGraphQLLifecycle() {
	// 1. FindAll
	all := s.GraphQLOp(s.handler, "findAllOrderItem", findAllOrderItemsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 2. FindOrderItemByOrder
	s.Require().NotZero(s.orderID)
	byOrder := s.GraphQLOp(s.handler, "findOrderItemByOrder", findOrderItemsByOrderQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	s.Equal("success", byOrder["status"])
	s.NotEmpty(tests.GQLFieldList(byOrder, "data"))

	// 3. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveOrderItem", findActiveOrderItemsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 4. FindByTrashed
	trashed := s.GraphQLOp(s.handler, "findByTrashedOrderItem", findTrashedOrderItemsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashed["status"])
}

func TestOrderItemGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemGraphQLTestSuite))
}
