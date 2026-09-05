package stats_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type CashierStatsSuite struct {
	tests.BaseTestSuite
	handler http.Handler
}

func (s *CashierStatsSuite) SetupSuite() {
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

// ─── Total Sales ─────────────────────────────────────────────────────────────

func (s *CashierStatsSuite) TestFindMonthlyTotalSales() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthlyTotalSales(input: { year: 2026, month: 1 }) {
			status
			message
			data { year month total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthlyTotalSales"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearlyTotalSales() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearlyTotalSales(input: { year: 2026 }) {
			status
			message
			data { year total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearlyTotalSales"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindMonthlyTotalSalesById() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthlyTotalSalesById(input: { year: 2026, month: 1, cashier_id: 1 }) {
			status
			message
			data { year month total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthlyTotalSalesById"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearlyTotalSalesById() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearlyTotalSalesById(input: { year: 2026, cashier_id: 1 }) {
			status
			message
			data { year total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearlyTotalSalesById"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindMonthlyTotalSalesByMerchant() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthlyTotalSalesByMerchant(input: { year: 2026, month: 1, merchant_id: 1 }) {
			status
			message
			data { year month total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthlyTotalSalesByMerchant"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearlyTotalSalesByMerchant() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearlyTotalSalesByMerchant(input: { year: 2026, merchant_id: 1 }) {
			status
			message
			data { year total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearlyTotalSalesByMerchant"].(map[string]interface{})["status"])
}

// ─── Sales ───────────────────────────────────────────────────────────────────

func (s *CashierStatsSuite) TestFindMonthSales() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthSales(input: { year: 2026 }) {
			status
			message
			data { month cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthSales"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearSales() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearSales(input: { year: 2026 }) {
			status
			message
			data { year cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearSales"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindMonthSalesByMerchant() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthSalesByMerchant(input: { year: 2026, merchant_id: 1 }) {
			status
			message
			data { month cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthSalesByMerchant"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearSalesByMerchant() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearSalesByMerchant(input: { year: 2026, merchant_id: 1 }) {
			status
			message
			data { year cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearSalesByMerchant"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindMonthSalesById() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findMonthSalesById(input: { year: 2026, cashier_id: 1 }) {
			status
			message
			data { month cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findMonthSalesById"].(map[string]interface{})["status"])
}

func (s *CashierStatsSuite) TestFindYearSalesById() {
	resp, err := graphtest.ExecuteGraphQL(s.handler, `{
		findYearSalesById(input: { year: 2026, cashier_id: 1 }) {
			status
			message
			data { year cashier_id order_count total_sales }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findYearSalesById"].(map[string]interface{})["status"])
}

func TestCashierStatsGraphqlHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CashierStatsSuite))
}
