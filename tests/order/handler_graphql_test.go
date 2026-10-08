package order_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler   http.Handler
	orderID   int
	userID    int
	cashierID int
	merchID   int
	prodID    int
}

func (s *OrderGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupCashierService()
	s.SetupTransactionService()
	s.SetupOrderService()

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	s.merchID = s.SeedMerchant(ctx, s.userID)
	s.prodID = s.SeedProduct(ctx, s.merchID, catID)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	s.Require().NoError(s.GormDB().Raw(
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Order GraphQL Cashier') RETURNING cashier_id`,
		s.merchID, s.userID,
	).Scan(&s.cashierID).Error)

	s.handler = s.NewGraphQLHandler()
}

const createOrderQuery = `mutation CreateOrder($input: CreateOrderInput!) {
  createOrder(input: $input) { status message data { id merchant_id cashier_id total_price } }
}`

const findOrderByIdQuery = `query FindOrderById($input: FindByIdOrderInput!) {
  findByIdOrder(input: $input) { status message data { id merchant_id cashier_id } }
}`

const findAllOrdersQuery = `query FindAllOrders($input: FindAllOrderInput!) {
  findAllOrder(input: $input) { status message data { id } pagination { total_records } }
}`

const findActiveOrdersQuery = `query FindActiveOrders($input: FindAllOrderInput!) {
  findByActiveOrder(input: $input) { status message data { id } }
}`

const findTrashedOrdersQuery = `query FindTrashedOrders($input: FindAllOrderInput!) {
  findByTrashedOrder(input: $input) { status message data { id } }
}`

const findOrderItemsByOrderQuery = `query FindOrderItemsByOrder($input: FindByIdOrderItemInput!) {
  findOrderItemByOrder(input: $input) { status message data { id order_id product_id } }
}`

const updateOrderQuery = `mutation UpdateOrder($input: UpdateOrderInput!) {
  updateOrder(input: $input) { status message data { id total_price } }
}`

const trashOrderQuery = `mutation TrashOrder($input: FindByIdOrderInput!) {
  trashedOrder(input: $input) { status message data { id deleted_at } }
}`

const restoreOrderQuery = `mutation RestoreOrder($input: FindByIdOrderInput!) {
  restoreOrder(input: $input) { status message data { id deleted_at } }
}`

const deleteOrderPermanentQuery = `mutation DeleteOrderPermanent($input: FindByIdOrderInput!) {
  deleteOrderPermanent(input: $input) { status message }
}`

const restoreAllOrdersQuery = `mutation RestoreAllOrders { restoreAllOrder { status message } }`

const deleteAllOrdersPermanentQuery = `mutation DeleteAllOrdersPermanent { deleteAllOrderPermanent { status message } }`

func (s *OrderGraphQLTestSuite) TestOrderGraphQLLifecycle() {
	// 1. Create
	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"merchant_id": s.merchID,
			"cashier_id":  s.cashierID,
			"items": []map[string]interface{}{
				{"product_id": s.prodID, "quantity": 1},
			},
		},
	}
	created := s.GraphQLOp(s.handler, "createOrder", createOrderQuery, createVars)
	s.Equal("success", created["status"])
	s.orderID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.orderID)

	// 2. FindById
	found := s.GraphQLOp(s.handler, "findByIdOrder", findOrderByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	s.Equal("success", found["status"])
	s.Equal(s.orderID, tests.GQLID(tests.GQLData(found), "id"))

	// 3. FindAll
	all := s.GraphQLOp(s.handler, "findAllOrder", findAllOrdersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 4. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveOrder", findActiveOrdersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 5. Update
	items := s.GraphQLOp(s.handler, "findOrderItemByOrder", findOrderItemsByOrderQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	itemList := tests.GQLFieldList(items, "data")
	s.Require().NotEmpty(itemList)
	orderItemID := tests.GQLID(itemList[0].(map[string]interface{}), "id")

	updateVars := map[string]interface{}{
		"input": map[string]interface{}{
			"order_id": s.orderID,
			"items": []map[string]interface{}{
				{"order_item_id": orderItemID, "product_id": s.prodID, "quantity": 1},
			},
		},
	}
	updated := s.GraphQLOp(s.handler, "updateOrder", updateOrderQuery, updateVars)
	s.Equal("success", updated["status"])

	// 6. Trash
	trashed := s.GraphQLOp(s.handler, "trashedOrder", trashOrderQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	s.Equal("success", trashed["status"])

	// 7. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedOrder", findTrashedOrdersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashedList["status"])

	// 8. Restore
	restored := s.GraphQLOp(s.handler, "restoreOrder", restoreOrderQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	s.Equal("success", restored["status"])

	// 9. DeletePermanent
	s.GraphQLRaw(s.handler, trashOrderQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	deleted := s.GraphQLOp(s.handler, "deleteOrderPermanent", deleteOrderPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.orderID}})
	s.Equal("success", deleted["status"])

	// 10. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllOrder", restoreAllOrdersQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 11. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllOrderPermanent", deleteAllOrdersPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestOrderGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGraphQLTestSuite))
}
