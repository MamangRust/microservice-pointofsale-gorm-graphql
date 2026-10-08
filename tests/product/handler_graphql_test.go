package product_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type ProductGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	productID  int
	merchantID int
	categoryID int
}

func (s *ProductGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	ctx := context.Background()
	userID := s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, userID)

	s.handler = s.NewGraphQLHandler()
}

const createProductQuery = `mutation CreateProduct($input: CreateProductInput!) {
  createProduct(input: $input) { status message data { id name merchantId categoryId } }
}`

const updateProductQuery = `mutation UpdateProduct($input: UpdateProductInput!) {
  updateProduct(input: $input) { status message data { id name } }
}`

const findProductByIdQuery = `query FindProductById($input: FindByIdProductInput!) {
  findByIdProduct(input: $input) { status message data { id name } }
}`

const findAllProductsQuery = `query FindAllProducts($input: FindAllProductInput) {
  findAllProduct(input: $input) { status message data { id } pagination { total_records } }
}`

const findActiveProductsQuery = `query FindActiveProducts($input: FindAllProductInput) {
  findByActiveProduct(input: $input) { status message data { id } }
}`

const findTrashedProductsQuery = `query FindTrashedProducts($input: FindAllProductInput) {
  findByTrashedProduct(input: $input) { status message data { id } }
}`

const findProductsByMerchantQuery = `query FindProductsByMerchant($input: FindAllProductMerchantInput!) {
  findByMerchantProduct(input: $input) { status message data { id } }
}`

const findProductsByCategoryQuery = `query FindProductsByCategory($input: FindAllProductCategoryInput!) {
  findByCategoryProduct(input: $input) { status message data { id } }
}`

const trashProductQuery = `mutation TrashProduct($input: FindByIdProductInput!) {
  trashedProduct(input: $input) { status message data { id deletedAt } }
}`

const restoreProductQuery = `mutation RestoreProduct($input: FindByIdProductInput!) {
  restoreProduct(input: $input) { status message data { id deletedAt } }
}`

const deleteProductPermanentQuery = `mutation DeleteProductPermanent($input: FindByIdProductInput!) {
  deleteProductPermanent(input: $input) { status message }
}`

const restoreAllProductsQuery = `mutation RestoreAllProducts { restoreAllProduct { status message } }`

const deleteAllProductsPermanentQuery = `mutation DeleteAllProductsPermanent { deleteAllProductPermanent { status message } }`

func (s *ProductGraphQLTestSuite) TestProductGraphQLLifecycle() {
	// 1. Create (image is a required Upload scalar)
	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"merchantId":   s.merchantID,
			"categoryId":   s.categoryID,
			"name":         "Test Product",
			"description":  "Test Description",
			"price":        1000,
			"countInStock": 10,
			"brand":        "Test Brand",
			"weight":       1,
			"image":        nil,
		},
	}
	created := s.GraphQLUploadOp(s.handler, "createProduct", createProductQuery, createVars, "variables.input.image", "product.jpg")
	s.Equal("success", created["status"])
	s.Equal("Test Product", tests.GQLData(created)["name"])
	s.productID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.productID)

	// 2. FindById
	found := s.GraphQLOp(s.handler, "findByIdProduct", findProductByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.productID}})
	s.Equal("success", found["status"])
	s.Equal(s.productID, tests.GQLID(tests.GQLData(found), "id"))

	// 3. FindAll
	all := s.GraphQLOp(s.handler, "findAllProduct", findAllProductsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", all["status"])

	// 4. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveProduct", findActiveProductsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", active["status"])

	// 5. FindByMerchant
	byMerchant := s.GraphQLOp(s.handler, "findByMerchantProduct", findProductsByMerchantQuery,
		map[string]interface{}{"input": map[string]interface{}{"merchantId": s.merchantID, "page": 1, "pageSize": 10}})
	s.Equal("success", byMerchant["status"])

	// 6. FindByCategory
	byCategory := s.GraphQLOp(s.handler, "findByCategoryProduct", findProductsByCategoryQuery,
		map[string]interface{}{"input": map[string]interface{}{"categoryName": "Seed Category", "page": 1, "pageSize": 10}})
	s.Equal("success", byCategory["status"])

	// 7. Update
	updateVars := map[string]interface{}{
		"input": map[string]interface{}{
			"productId":    s.productID,
			"merchantId":   s.merchantID,
			"categoryId":   s.categoryID,
			"name":         "Updated Product",
			"description":  "Updated Description",
			"price":        2000,
			"countInStock": 20,
			"brand":        "Updated Brand",
			"weight":       2,
			"image":        nil,
		},
	}
	updated := s.GraphQLUploadOp(s.handler, "updateProduct", updateProductQuery, updateVars, "variables.input.image", "updated.jpg")
	s.Equal("success", updated["status"])
	s.Equal("Updated Product", tests.GQLData(updated)["name"])

	// 8. Trash
	trashed := s.GraphQLOp(s.handler, "trashedProduct", trashProductQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.productID}})
	s.Equal("success", trashed["status"])

	// 9. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedProduct", findTrashedProductsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "pageSize": 10}})
	s.Equal("success", trashedList["status"])

	// 10. Restore
	restored := s.GraphQLOp(s.handler, "restoreProduct", restoreProductQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.productID}})
	s.Equal("success", restored["status"])

	// 11. DeletePermanent
	s.GraphQLRaw(s.handler, trashProductQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.productID}})
	deleted := s.GraphQLOp(s.handler, "deleteProductPermanent", deleteProductPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.productID}})
	s.Equal("success", deleted["status"])

	// 12. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllProduct", restoreAllProductsQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 13. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllProductPermanent", deleteAllProductsPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestProductGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGraphQLTestSuite))
}
