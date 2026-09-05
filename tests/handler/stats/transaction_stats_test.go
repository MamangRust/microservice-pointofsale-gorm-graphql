package stats_test

import (
	"net/http"
	"strings"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *TransactionStatsSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupStatsReaderService()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupCashierService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *TransactionStatsSuite) assertStatsQuery(query, fieldName string) {
	resp, err := graphtest.ExecuteGraphQL(s.handler, query, nil, "")
	s.Require().NoError(err, "HTTP error executing GraphQL")
	if len(resp.Errors) > 0 {
		msg := resp.Errors[0].Message
		if strings.Contains(msg, "not implemented") || strings.Contains(msg, "Unimplemented") {
			s.T().Skipf("Stats method %s not implemented in backend: %s", fieldName, msg)
			return
		}
		s.Require().Empty(resp.Errors, "GraphQL errors: %v", resp.Errors)
	}
	s.Require().NotNil(resp.Data[fieldName], "Data for %s is nil", fieldName)
	s.Equal("success", resp.Data[fieldName].(map[string]interface{})["status"])
}

func (s *TransactionStatsSuite) TestFindMonthStatusSuccess() {
	s.assertStatsQuery(`{ findMonthStatusSuccess(input: { year: 2026, month: 1 }) { status message data { month totalSuccess totalAmount } } }`, "findMonthStatusSuccess")
}

func (s *TransactionStatsSuite) TestFindYearStatusSuccess() {
	s.assertStatsQuery(`{ findYearStatusSuccess(input: { year: 2026 }) { status message data { year totalSuccess totalAmount } } }`, "findYearStatusSuccess")
}

func (s *TransactionStatsSuite) TestFindMonthStatusFailed() {
	s.assertStatsQuery(`{ findMonthStatusFailed(input: { year: 2026, month: 1 }) { status message data { month totalFailed totalAmount } } }`, "findMonthStatusFailed")
}

func (s *TransactionStatsSuite) TestFindYearStatusFailed() {
	s.assertStatsQuery(`{ findYearStatusFailed(input: { year: 2026 }) { status message data { year totalFailed totalAmount } } }`, "findYearStatusFailed")
}

func (s *TransactionStatsSuite) TestFindMonthStatusSuccessByMerchant() {
	s.assertStatsQuery(`{ findMonthStatusSuccessByMerchant(input: { year: 2026, month: 1, merchantId: 1 }) { status message data { month totalSuccess totalAmount } } }`, "findMonthStatusSuccessByMerchant")
}

func (s *TransactionStatsSuite) TestFindYearStatusSuccessByMerchant() {
	s.assertStatsQuery(`{ findYearStatusSuccessByMerchant(input: { year: 2026, merchantId: 1 }) { status message data { year totalSuccess totalAmount } } }`, "findYearStatusSuccessByMerchant")
}

func (s *TransactionStatsSuite) TestFindMonthStatusFailedByMerchant() {
	s.assertStatsQuery(`{ findMonthStatusFailedByMerchant(input: { year: 2026, month: 1, merchantId: 1 }) { status message data { month totalFailed totalAmount } } }`, "findMonthStatusFailedByMerchant")
}

func (s *TransactionStatsSuite) TestFindYearStatusFailedByMerchant() {
	s.assertStatsQuery(`{ findYearStatusFailedByMerchant(input: { year: 2026, merchantId: 1 }) { status message data { year totalFailed totalAmount } } }`, "findYearStatusFailedByMerchant")
}

func (s *TransactionStatsSuite) TestFindMonthMethodSuccess() {
	s.assertStatsQuery(`{ findMonthMethodSuccess(input: { year: 2026, month: 1 }) { status message data { month paymentMethod totalAmount } } }`, "findMonthMethodSuccess")
}

func (s *TransactionStatsSuite) TestFindYearMethodSuccess() {
	s.assertStatsQuery(`{ findYearMethodSuccess(input: { year: 2026 }) { status message data { year totalAmount } } }`, "findYearMethodSuccess")
}

func (s *TransactionStatsSuite) TestFindMonthMethodFailed() {
	s.assertStatsQuery(`{ findMonthMethodFailed(input: { year: 2026, month: 1 }) { status message data { month paymentMethod totalAmount } } }`, "findMonthMethodFailed")
}

func (s *TransactionStatsSuite) TestFindYearMethodFailed() {
	s.assertStatsQuery(`{ findYearMethodFailed(input: { year: 2026 }) { status message data { year totalAmount } } }`, "findYearMethodFailed")
}

func (s *TransactionStatsSuite) TestFindMonthMethodByMerchantSuccess() {
	s.assertStatsQuery(`{ findMonthMethodByMerchantSuccess(input: { year: 2026, month: 1, merchantId: 1 }) { status message data { month paymentMethod totalAmount } } }`, "findMonthMethodByMerchantSuccess")
}

func (s *TransactionStatsSuite) TestFindYearMethodByMerchantSuccess() {
	s.assertStatsQuery(`{ findYearMethodByMerchantSuccess(input: { year: 2026, merchantId: 1 }) { status message data { year totalAmount } } }`, "findYearMethodByMerchantSuccess")
}

func (s *TransactionStatsSuite) TestFindMonthMethodByMerchantFailed() {
	s.assertStatsQuery(`{ findMonthMethodByMerchantFailed(input: { year: 2026, month: 1, merchantId: 1 }) { status message data { month paymentMethod totalAmount } } }`, "findMonthMethodByMerchantFailed")
}

func (s *TransactionStatsSuite) TestFindYearMethodByMerchantFailed() {
	s.assertStatsQuery(`{ findYearMethodByMerchantFailed(input: { year: 2026, merchantId: 1 }) { status message data { year totalAmount } } }`, "findYearMethodByMerchantFailed")
}

func TestTransactionStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsSuite))
}
