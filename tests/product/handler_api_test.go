package product_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type ProductApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	productID  int
	merchantID int
	categoryID int
}

func (s *ProductApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *ProductApiTestSuite) TestProductApiLifecycle() {
	// Product uses Upload scalar for images — GraphQL multipart is complex.
	// We test queries and mutations that don't require image upload.
	// Create is skipped because it requires file upload via GraphQL multipart.

	// Seed a product via DB directly for query testing
	var prodID int
	err := s.GormDB().WithContext(s.Ctx).Raw(
		`INSERT INTO products (merchant_id, category_id, name, description, price, count_in_stock, brand, weight, slug_product, image_product)
		 VALUES (?, ?, 'Seed Product', 'Seed Desc', 10000, 10, 'Brand', 1000, 'seed-product', 'seed.jpg')
		 RETURNING product_id`, s.merchantID, s.categoryID,
	).Scan(&prodID).Error
	s.Require().NoError(err)
	s.productID = prodID

	// 1. FindById
	resp, err := graphtest.ExecuteGraphQL(s.handler, `query { findByIdProduct(input: { id: `+strconv.Itoa(s.productID)+` }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdProduct"].(map[string]interface{})["status"])

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllProduct(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllProduct"].(map[string]interface{})["status"])

	// 3. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveProduct(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveProduct"].(map[string]interface{})["status"])

	// 4. FindByMerchant
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByMerchantProduct(input: { merchantId: `+strconv.Itoa(s.merchantID)+`, page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByMerchantProduct"].(map[string]interface{})["status"])

	// 5. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedProduct(input: { id: `+strconv.Itoa(s.productID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedProduct"].(map[string]interface{})["status"])

	// 6. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedProduct(input: { page: 1, pageSize: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedProduct"].(map[string]interface{})["status"])

	// 7. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreProduct(input: { id: `+strconv.Itoa(s.productID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreProduct"].(map[string]interface{})["status"])

	// 8. Re-trash for permanent delete
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedProduct(input: { id: `+strconv.Itoa(s.productID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedProduct"].(map[string]interface{})["status"])

	// 9. DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteProductPermanent(input: { id: `+strconv.Itoa(s.productID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteProductPermanent"].(map[string]interface{})["status"])

	// 9. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllProduct { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllProduct"].(map[string]interface{})["status"])

	// 10. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllProductPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllProductPermanent"].(map[string]interface{})["status"])
}

func TestProductApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductApiTestSuite))
}
