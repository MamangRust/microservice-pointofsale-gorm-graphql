package category_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type CategoryApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	categoryID int
}

func (s *CategoryApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupUserService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *CategoryApiTestSuite) TestCategoryApiLifecycle() {
	// 1. Create
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		createCategory(input: { name: "Test Category", description: "Test Description" }) {
			status message data { id name description }
		}
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createCategory"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.Equal("Test Category", data["name"])
	s.categoryID = int(data["id"].(float64))

	// 2. FindById
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByIdCategory(input: { id: `+strconv.Itoa(s.categoryID)+` }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdCategory"].(map[string]interface{})["status"])

	// 3. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllCategory(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllCategory"].(map[string]interface{})["status"])

	// 4. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveCategory(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveCategory"].(map[string]interface{})["status"])

	// 5. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		updateCategory(input: { category_id: `+strconv.Itoa(s.categoryID)+`, name: "Updated Category", description: "Updated Description" }) {
			status message data { id name }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateCategory"].(map[string]interface{})["status"])

	// 6. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedCategory(input: { id: `+strconv.Itoa(s.categoryID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedCategory"].(map[string]interface{})["status"])

	// 7. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedCategory(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedCategory"].(map[string]interface{})["status"])

	// 8. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreCategory(input: { id: `+strconv.Itoa(s.categoryID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreCategory"].(map[string]interface{})["status"])

	// 9. DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteCategoryPermanent(input: { id: `+strconv.Itoa(s.categoryID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteCategoryPermanent"].(map[string]interface{})["status"])

	// 10. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllCategory { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllCategory"].(map[string]interface{})["status"])

	// 11. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllCategoryPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllCategoryPermanent"].(map[string]interface{})["status"])
}

func TestCategoryApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryApiTestSuite))
}
