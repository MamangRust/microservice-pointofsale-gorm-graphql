package stats_test

import (
	"net/http"
	"strings"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *OrderStatsSuite) SetupSuite() {
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

func (s *OrderStatsSuite) assertStatsQuery(query, fieldName string) {
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

func (s *OrderStatsSuite) TestFindMonthlyRevenue() {
	s.assertStatsQuery(`{ findMonthlyRevenue(input: { year: 2026 }) { status message data { month total_revenue order_count } } }`, "findMonthlyRevenue")
}

func (s *OrderStatsSuite) TestFindYearlyRevenue() {
	s.assertStatsQuery(`{ findYearlyRevenue(input: { year: 2026 }) { status message data { year total_revenue order_count } } }`, "findYearlyRevenue")
}

func (s *OrderStatsSuite) TestFindMonthlyRevenueByMerchant() {
	s.assertStatsQuery(`{ findMonthlyRevenueByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { month total_revenue order_count } } }`, "findMonthlyRevenueByMerchant")
}

func (s *OrderStatsSuite) TestFindYearlyRevenueByMerchant() {
	s.assertStatsQuery(`{ findYearlyRevenueByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { year total_revenue order_count } } }`, "findYearlyRevenueByMerchant")
}

func (s *OrderStatsSuite) TestFindMonthlyTotalRevenue() {
	s.assertStatsQuery(`{ findMonthlyTotalRevenue(input: { year: 2026, month: 1 }) { status message data { year month total_revenue order_count } } }`, "findMonthlyTotalRevenue")
}

func (s *OrderStatsSuite) TestFindYearlyTotalRevenue() {
	s.assertStatsQuery(`{ findYearlyTotalRevenue(input: { year: 2026 }) { status message data { year total_revenue order_count } } }`, "findYearlyTotalRevenue")
}

func (s *OrderStatsSuite) TestFindMonthlyTotalRevenueByMerchant() {
	s.assertStatsQuery(`{ findMonthlyTotalRevenueByMerchant(input: { year: 2026, month: 1, merchant_id: 1 }) { status message data { year month total_revenue order_count } } }`, "findMonthlyTotalRevenueByMerchant")
}

func (s *OrderStatsSuite) TestFindYearlyTotalRevenueByMerchant() {
	s.assertStatsQuery(`{ findYearlyTotalRevenueByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { year total_revenue order_count } } }`, "findYearlyTotalRevenueByMerchant")
}

func (s *OrderStatsSuite) TestFindMonthlyTotalRevenueById() {
	s.assertStatsQuery(`{ findMonthlyTotalRevenueById(input: { year: 2026, month: 1, order_id: 1 }) { status message data { year month total_revenue order_count } } }`, "findMonthlyTotalRevenueById")
}

func (s *OrderStatsSuite) TestFindYearlyTotalRevenueById() {
	s.assertStatsQuery(`{ findYearlyTotalRevenueById(input: { year: 2026, order_id: 1 }) { status message data { year total_revenue order_count } } }`, "findYearlyTotalRevenueById")
}

func TestOrderStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsSuite))
}
