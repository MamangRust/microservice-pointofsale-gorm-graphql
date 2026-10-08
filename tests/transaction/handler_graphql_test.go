package transaction_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler       http.Handler
	transactionID int
	userID        int
	cashierID     int
	merchID       int
	orderID       int
}

func (s *TransactionGraphQLTestSuite) SetupSuite() {
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
	s.userID = s.SeedUser(ctx)
	s.merchID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchID, catID)
	s.orderID = s.SeedOrder(ctx, s.userID, s.merchID, prodID)
	s.SeedOrderItem(ctx, s.orderID, prodID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	s.Require().NoError(s.GormDB().Raw(
		`SELECT cashier_id FROM cashiers WHERE user_id = $1 AND merchant_id = $2 AND deleted_at IS NULL LIMIT 1`,
		s.userID, s.merchID,
	).Scan(&s.cashierID).Error)

	s.handler = s.NewGraphQLHandler()
}

const createTransactionQuery = `mutation CreateTransaction($input: CreateTransactionInput!) {
  createTransaction(input: $input) { status message data { id orderId merchantId paymentMethod amount paymentStatus } }
}`

const updateTransactionQuery = `mutation UpdateTransaction($input: UpdateTransactionInput!) {
  updateTransaction(input: $input) { status message data { id paymentMethod amount } }
}`

const findAllTransactionsQuery = `query FindAllTransactions($input: FindAllTransactionInput) {
  findAllTransaction(input: $input) { status message data { id } pagination { total_records } }
}`

const findTransactionByIdQuery = `query FindTransactionById($input: FindByIdTransactionInput!) {
  findByIdTransaction(input: $input) { status message data { id paymentMethod } }
}`

const findTransactionsByMerchantQuery = `query FindTransactionsByMerchant($input: FindAllTransactionMerchantInput!) {
  findByMerchantTransaction(input: $input) { status message data { id } }
}`

const findActiveTransactionsQuery = `query FindActiveTransactions($input: FindAllTransactionInput) {
  findByActiveTransaction(input: $input) { status message data { id } }
}`

const findTrashedTransactionsQuery = `query FindTrashedTransactions($input: FindAllTransactionInput) {
  findByTrashedTransaction(input: $input) { status message data { id } }
}`

const trashTransactionQuery = `mutation TrashTransaction($input: FindByIdTransactionInput!) {
  trashedTransaction(input: $input) { status message data { id deletedAt } }
}`

const restoreTransactionQuery = `mutation RestoreTransaction($input: FindByIdTransactionInput!) {
  restoreTransaction(input: $input) { status message data { id deletedAt } }
}`

const deleteTransactionPermanentQuery = `mutation DeleteTransactionPermanent($input: FindByIdTransactionInput!) {
  deleteTransactionPermanent(input: $input) { status message }
}`

const restoreAllTransactionsQuery = `mutation RestoreAllTransactions { restoreAllTransaction { status message } }`

const deleteAllTransactionsPermanentQuery = `mutation DeleteAllTransactionsPermanent { deleteAllTransactionPermanent { status message } }`

func (s *TransactionGraphQLTestSuite) TestTransactionGraphQLLifecycle() {
	// 1. Create
	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"orderId":       s.orderID,
			"cashierId":     s.cashierID,
			"paymentMethod": "Transfer Bank",
			"amount":        100000,
			"paymentStatus": "pending",
		},
	}
	created := s.GraphQLOp(s.handler, "createTransaction", createTransactionQuery, createVars)
	s.Equal("success", created["status"])
	s.transactionID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.transactionID)

	// 2. FindAll
	all := s.GraphQLOp(s.handler, "findAllTransaction", findAllTransactionsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", all["status"])

	// 3. FindById
	found := s.GraphQLOp(s.handler, "findByIdTransaction", findTransactionByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.transactionID}})
	s.Equal("success", found["status"])
	s.Equal(s.transactionID, tests.GQLID(tests.GQLData(found), "id"))

	// 4. FindByMerchant
	byMerchant := s.GraphQLOp(s.handler, "findByMerchantTransaction", findTransactionsByMerchantQuery,
		map[string]interface{}{"input": map[string]interface{}{"merchantId": s.merchID, "page": 1, "pageSize": 10}})
	s.Equal("success", byMerchant["status"])

	// 5. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveTransaction", findActiveTransactionsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", active["status"])

	// 6. Update
	updateVars := map[string]interface{}{
		"input": map[string]interface{}{
			"transactionId": s.transactionID,
			"orderId":       s.orderID,
			"cashierId":     s.cashierID,
			"paymentMethod": "GOPAY",
			"amount":        100000,
			"paymentStatus": "success",
		},
	}
	updated := s.GraphQLOp(s.handler, "updateTransaction", updateTransactionQuery, updateVars)
	s.Equal("success", updated["status"])

	// 7. Trash
	trashed := s.GraphQLOp(s.handler, "trashedTransaction", trashTransactionQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.transactionID}})
	s.Equal("success", trashed["status"])

	// 8. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedTransaction", findTrashedTransactionsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", trashedList["status"])

	// 9. Restore
	restored := s.GraphQLOp(s.handler, "restoreTransaction", restoreTransactionQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.transactionID}})
	s.Equal("success", restored["status"])

	// 10. DeletePermanent
	s.GraphQLRaw(s.handler, trashTransactionQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.transactionID}})
	deleted := s.GraphQLOp(s.handler, "deleteTransactionPermanent", deleteTransactionPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.transactionID}})
	s.Equal("success", deleted["status"])

	// 11. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllTransaction", restoreAllTransactionsQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 12. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllTransactionPermanent", deleteAllTransactionsPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestTransactionGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionGraphQLTestSuite))
}
