package category_test

import (
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type CategoryGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	categoryID int
}

func (s *CategoryGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupCategoryService()

	s.handler = s.NewGraphQLHandler()
}

const createCategoryQuery = `mutation CreateCategory($input: CreateCategoryRequest!) {
  createCategory(input: $input) { status message data { id name description } }
}`

const updateCategoryQuery = `mutation UpdateCategory($input: UpdateCategoryRequest!) {
  updateCategory(input: $input) { status message data { id name description } }
}`

const findCategoryByIdQuery = `query FindCategoryById($input: FindByIdCategoryRequest!) {
  findByIdCategory(input: $input) { status message data { id name } }
}`

const findAllCategoriesQuery = `query FindAllCategories($input: FindAllCategoryRequest) {
  findAllCategory(input: $input) { status message data { id name } pagination { total_records } }
}`

const findActiveCategoriesQuery = `query FindActiveCategories($input: FindAllCategoryRequest) {
  findByActiveCategory(input: $input) { status message data { id name } pagination { total_records } }
}`

const findTrashedCategoriesQuery = `query FindTrashedCategories($input: FindAllCategoryRequest) {
  findByTrashedCategory(input: $input) { status message data { id name } pagination { total_records } }
}`

const trashCategoryQuery = `mutation TrashCategory($input: FindByIdCategoryRequest!) {
  trashedCategory(input: $input) { status message data { id name deleted_at } }
}`

const restoreCategoryQuery = `mutation RestoreCategory($input: FindByIdCategoryRequest!) {
  restoreCategory(input: $input) { status message data { id name deleted_at } }
}`

const deleteCategoryPermanentQuery = `mutation DeleteCategoryPermanent($input: FindByIdCategoryRequest!) {
  deleteCategoryPermanent(input: $input) { status message }
}`

const restoreAllCategoriesQuery = `mutation RestoreAllCategories { restoreAllCategory { status message } }`

const deleteAllCategoriesPermanentQuery = `mutation DeleteAllCategoriesPermanent { deleteAllCategoryPermanent { status message } }`

func (s *CategoryGraphQLTestSuite) TestCategoryGraphQLLifecycle() {
	// 1. Create
	created := s.GraphQLOp(s.handler, "createCategory", createCategoryQuery, map[string]interface{}{
		"input": map[string]interface{}{"name": "Test Category", "description": "Test Description"},
	})
	s.Equal("success", created["status"])
	s.Equal("Test Category", tests.GQLData(created)["name"])
	s.categoryID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.categoryID)

	// 2. FindById
	found := s.GraphQLOp(s.handler, "findByIdCategory", findCategoryByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.categoryID}})
	s.Equal("success", found["status"])
	s.Equal(s.categoryID, tests.GQLID(tests.GQLData(found), "id"))

	// 3. FindAll
	all := s.GraphQLOp(s.handler, "findAllCategory", findAllCategoriesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 4. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveCategory", findActiveCategoriesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 5. Update
	updated := s.GraphQLOp(s.handler, "updateCategory", updateCategoryQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"category_id": s.categoryID,
			"name":        "Updated Category",
			"description": "Updated Description",
		},
	})
	s.Equal("success", updated["status"])
	s.Equal("Updated Category", tests.GQLData(updated)["name"])

	// 6. Trash
	trashed := s.GraphQLOp(s.handler, "trashedCategory", trashCategoryQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.categoryID}})
	s.Equal("success", trashed["status"])

	// 7. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedCategory", findTrashedCategoriesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashedList["status"])

	// 8. Restore
	restored := s.GraphQLOp(s.handler, "restoreCategory", restoreCategoryQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.categoryID}})
	s.Equal("success", restored["status"])

	// 9. DeletePermanent
	s.GraphQLRaw(s.handler, trashCategoryQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.categoryID}})
	deleted := s.GraphQLOp(s.handler, "deleteCategoryPermanent", deleteCategoryPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.categoryID}})
	s.Equal("success", deleted["status"])

	// 10. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllCategory", restoreAllCategoriesQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 11. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllCategoryPermanent", deleteAllCategoriesPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestCategoryGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGraphQLTestSuite))
}
