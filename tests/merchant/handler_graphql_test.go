package merchant_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type MerchantGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	merchantID int
	userID     int
}

func (s *MerchantGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()

	s.userID = s.SeedUser(context.Background())
	s.handler = s.NewGraphQLHandler()
}

const createMerchantQuery = `mutation CreateMerchant($input: CreateMerchantInput!) {
  createMerchant(input: $input) { status message data { id user_id name status } }
}`

const findMerchantByIdQuery = `query FindMerchantById($input: FindByIdMerchantInput!) {
  findByIdMerchant(input: $input) { status message data { id name } }
}`

const findAllMerchantsQuery = `query FindAllMerchants($input: FindAllMerchantInput!) {
  findAllMerchant(input: $input) { status message data { id } pagination { total_records } }
}`

const findActiveMerchantsQuery = `query FindActiveMerchants($input: FindAllMerchantInput!) {
  findByActiveMerchant(input: $input) { status message data { id } }
}`

const findTrashedMerchantsQuery = `query FindTrashedMerchants($input: FindAllMerchantInput!) {
  findByTrashedMerchant(input: $input) { status message data { id } }
}`

const updateMerchantQuery = `mutation UpdateMerchant($input: UpdateMerchantInput!) {
  updateMerchant(input: $input) { status message data { id name } }
}`

const trashMerchantQuery = `mutation TrashMerchant($input: FindByIdMerchantInput!) {
  trashedMerchant(input: $input) { status message data { id deleted_at } }
}`

const restoreMerchantQuery = `mutation RestoreMerchant($input: FindByIdMerchantInput!) {
  restoreMerchant(input: $input) { status message data { id deleted_at } }
}`

const deleteMerchantPermanentQuery = `mutation DeleteMerchantPermanent($input: FindByIdMerchantInput!) {
  deleteMerchantPermanent(input: $input) { status message }
}`

const restoreAllMerchantsQuery = `mutation RestoreAllMerchants { restoreAllMerchant { status message } }`

const deleteAllMerchantsPermanentQuery = `mutation DeleteAllMerchantsPermanent { deleteAllMerchantPermanent { status message } }`

func (s *MerchantGraphQLTestSuite) TestMerchantGraphQLLifecycle() {
	// 1. Create
	createVars := map[string]interface{}{
		"input": map[string]interface{}{
			"user_id":       s.userID,
			"name":          "Test Merchant",
			"description":   "Test Description",
			"address":       "Test Address",
			"contact_email": "merchant@example.com",
			"contact_phone": "123456789",
			"status":        "active",
		},
	}
	created := s.GraphQLOp(s.handler, "createMerchant", createMerchantQuery, createVars)
	s.Equal("success", created["status"])
	s.Equal("Test Merchant", tests.GQLData(created)["name"])
	s.merchantID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.merchantID)

	// 2. FindById
	found := s.GraphQLOp(s.handler, "findByIdMerchant", findMerchantByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.merchantID}})
	s.Equal("success", found["status"])
	s.Equal(s.merchantID, tests.GQLID(tests.GQLData(found), "id"))

	// 3. FindAll
	all := s.GraphQLOp(s.handler, "findAllMerchant", findAllMerchantsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 4. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveMerchant", findActiveMerchantsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 5. Update
	updateVars := map[string]interface{}{
		"input": map[string]interface{}{
			"merchant_id":   s.merchantID,
			"user_id":       s.userID,
			"name":          "Updated Merchant",
			"description":   "Updated Description",
			"address":       "Updated Address",
			"contact_email": "updated@example.com",
			"contact_phone": "987654321",
			"status":        "active",
		},
	}
	updated := s.GraphQLOp(s.handler, "updateMerchant", updateMerchantQuery, updateVars)
	s.Equal("success", updated["status"])
	s.Equal("Updated Merchant", tests.GQLData(updated)["name"])

	// 6. Trash
	trashed := s.GraphQLOp(s.handler, "trashedMerchant", trashMerchantQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.merchantID}})
	s.Equal("success", trashed["status"])

	// 7. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedMerchant", findTrashedMerchantsQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashedList["status"])

	// 8. Restore
	restored := s.GraphQLOp(s.handler, "restoreMerchant", restoreMerchantQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.merchantID}})
	s.Equal("success", restored["status"])

	// 9. DeletePermanent
	s.GraphQLRaw(s.handler, trashMerchantQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.merchantID}})
	deleted := s.GraphQLOp(s.handler, "deleteMerchantPermanent", deleteMerchantPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.merchantID}})
	s.Equal("success", deleted["status"])

	// 10. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllMerchant", restoreAllMerchantsQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 11. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllMerchantPermanent", deleteAllMerchantsPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestMerchantGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantGraphQLTestSuite))
}
