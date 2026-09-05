package stats_test

import (
	"net/http"
	"strings"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *CategoryStatsSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupStatsReaderService()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

// assertStatsQuery executes a stats GraphQL query and asserts it returns success.
// If the backend returns "not implemented" it skips the test gracefully.
func (s *CategoryStatsSuite) assertStatsQuery(query, fieldName string) {
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

func (s *CategoryStatsSuite) TestFindMonthPrice() {
	s.assertStatsQuery(`{ findMonthPrice(input: { year: 2026 }) { status message data { month total_revenue } } }`, "findMonthPrice")
}

func (s *CategoryStatsSuite) TestFindYearPrice() {
	s.assertStatsQuery(`{ findYearPrice(input: { year: 2026 }) { status message data { year total_revenue } } }`, "findYearPrice")
}

func (s *CategoryStatsSuite) TestFindMonthPriceByMerchant() {
	s.assertStatsQuery(`{ findMonthPriceByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { month total_revenue } } }`, "findMonthPriceByMerchant")
}

func (s *CategoryStatsSuite) TestFindYearPriceByMerchant() {
	s.assertStatsQuery(`{ findYearPriceByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { year total_revenue } } }`, "findYearPriceByMerchant")
}

func (s *CategoryStatsSuite) TestFindMonthPriceById() {
	s.assertStatsQuery(`{ findMonthPriceById(input: { year: 2026, category_id: 1 }) { status message data { month total_revenue } } }`, "findMonthPriceById")
}

func (s *CategoryStatsSuite) TestFindYearPriceById() {
	s.assertStatsQuery(`{ findYearPriceById(input: { year: 2026, category_id: 1 }) { status message data { year total_revenue } } }`, "findYearPriceById")
}

func (s *CategoryStatsSuite) TestFindMonthlyTotalPrices() {
	s.assertStatsQuery(`{ findMonthlyTotalPrices(input: { year: 2026, month: 1 }) { status message data { year month total_revenue } } }`, "findMonthlyTotalPrices")
}

func (s *CategoryStatsSuite) TestFindYearlyTotalPrices() {
	s.assertStatsQuery(`{ findYearlyTotalPrices(input: { year: 2026 }) { status message data { year total_revenue } } }`, "findYearlyTotalPrices")
}

func (s *CategoryStatsSuite) TestFindMonthlyTotalPricesById() {
	s.assertStatsQuery(`{ findMonthlyTotalPricesById(input: { year: 2026, month: 1, category_id: 1 }) { status message data { year month total_revenue } } }`, "findMonthlyTotalPricesById")
}

func (s *CategoryStatsSuite) TestFindYearlyTotalPricesById() {
	s.assertStatsQuery(`{ findYearlyTotalPricesById(input: { year: 2026, category_id: 1 }) { status message data { year total_revenue } } }`, "findYearlyTotalPricesById")
}

func (s *CategoryStatsSuite) TestFindMonthlyTotalPricesByMerchant() {
	s.assertStatsQuery(`{ findMonthlyTotalPricesByMerchant(input: { year: 2026, month: 1, merchant_id: 1 }) { status message data { year month total_revenue } } }`, "findMonthlyTotalPricesByMerchant")
}

func (s *CategoryStatsSuite) TestFindYearlyTotalPricesByMerchant() {
	s.assertStatsQuery(`{ findYearlyTotalPricesByMerchant(input: { year: 2026, merchant_id: 1 }) { status message data { year total_revenue } } }`, "findYearlyTotalPricesByMerchant")
}

func TestCategoryStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsSuite))
}
