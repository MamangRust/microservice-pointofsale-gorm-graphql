package transaction_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionApiTestSuite struct {
	tests.BaseTestSuite
	handler       http.Handler
	transactionID int
	userID        int
	cashierID     int
	merchID       int
	orderID       int
}

func (s *TransactionApiTestSuite) SetupSuite() {
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
	s.userID = s.SeedUser(ctx)
	s.merchID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchID, catID)
	s.orderID = s.SeedOrder(ctx, s.userID, s.merchID, prodID)
	s.SeedOrderItem(ctx, s.orderID, prodID)

	err := s.GormDB().WithContext(ctx).Raw(
		`INSERT INTO cashiers (merchant_id, user_id, name) VALUES (?, ?, 'Txn Api Cashier') RETURNING cashier_id`,
		s.merchID, s.userID,
	).Scan(&s.cashierID).Error
	s.Require().NoError(err)

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *TransactionApiTestSuite) TestTransactionApiLifecycle() {
	// 1. Create
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		createTransaction(input: {
			orderId: `+strconv.Itoa(s.orderID)+`, cashierId: `+strconv.Itoa(s.cashierID)+`,
			paymentMethod: "Transfer Bank", amount: 100000, paymentStatus: "success"
		}) {
			status message data { id }
		}
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createTransaction"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.transactionID = int(data["id"].(float64))

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllTransaction(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllTransaction"].(map[string]interface{})["status"])

	// 3. FindById
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByIdTransaction(input: { id: `+strconv.Itoa(s.transactionID)+` }) { status message data { id } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdTransaction"].(map[string]interface{})["status"])

	// 4. FindByMerchant
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByMerchantTransaction(input: { merchantId: `+strconv.Itoa(s.merchID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByMerchantTransaction"].(map[string]interface{})["status"])

	// 5. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveTransaction(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveTransaction"].(map[string]interface{})["status"])

	// 6. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		updateTransaction(input: {
			transactionId: `+strconv.Itoa(s.transactionID)+`, orderId: `+strconv.Itoa(s.orderID)+`,
			cashierId: `+strconv.Itoa(s.cashierID)+`, paymentMethod: "GOPAY", amount: 100000, paymentStatus: "success"
		}) {
			status message
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateTransaction"].(map[string]interface{})["status"])

	// 7. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedTransaction(input: { id: `+strconv.Itoa(s.transactionID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedTransaction"].(map[string]interface{})["status"])

	// 8. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedTransaction(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedTransaction"].(map[string]interface{})["status"])

	// 9. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreTransaction(input: { id: `+strconv.Itoa(s.transactionID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreTransaction"].(map[string]interface{})["status"])

	// 10. Re-trash for permanent delete
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedTransaction(input: { id: `+strconv.Itoa(s.transactionID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedTransaction"].(map[string]interface{})["status"])

	// 11. DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteTransactionPermanent(input: { id: `+strconv.Itoa(s.transactionID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteTransactionPermanent"].(map[string]interface{})["status"])

	// 11. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllTransaction { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllTransaction"].(map[string]interface{})["status"])

	// 12. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllTransactionPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllTransactionPermanent"].(map[string]interface{})["status"])
}

func TestTransactionApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionApiTestSuite))
}
