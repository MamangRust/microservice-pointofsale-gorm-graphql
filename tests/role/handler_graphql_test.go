package role_test

import (
	"context"
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type RoleGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	roleID  int
	userID  int
}

func (s *RoleGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()

	s.userID = s.SeedUser(context.Background())

	s.handler = s.NewGraphQLHandler()
}

const createRoleQuery = `mutation CreateRole($input: CreateRoleInput!) {
  createRole(input: $input) { status message data { id name } }
}`

const updateRoleQuery = `mutation UpdateRole($input: UpdateRoleInput!) {
  updateRole(input: $input) { status message data { id name } }
}`

const findRoleByIdQuery = `query FindRoleById($input: FindByIdRoleInput!) {
  findByIdRole(input: $input) { status message data { id name } }
}`

const findAllRolesQuery = `query FindAllRoles($input: FindAllRoleInput) {
  findAllRole(input: $input) { status message data { id name } pagination { total_records } }
}`

const findActiveRolesQuery = `query FindActiveRoles($input: FindAllRoleInput) {
  findByActiveRole(input: $input) { status message data { id name } pagination { total_records } }
}`

const findTrashedRolesQuery = `query FindTrashedRoles($input: FindAllRoleInput) {
  findByTrashedRole(input: $input) { status message data { id name } pagination { total_records } }
}`

const findRolesByUserQuery = `query FindRolesByUser($input: FindByIdUserRoleInput!) {
  findByUserIdRole(input: $input) { status message data { id name } }
}`

const trashRoleQuery = `mutation TrashRole($input: FindByIdRoleInput!) {
  trashedRole(input: $input) { status message data { id name deleted_at } }
}`

const restoreRoleQuery = `mutation RestoreRole($input: FindByIdRoleInput!) {
  restoreRole(input: $input) { status message data { id name deleted_at } }
}`

const deleteRolePermanentQuery = `mutation DeleteRolePermanent($input: FindByIdRoleInput!) {
  deleteRolePermanent(input: $input) { status message }
}`

const restoreAllRolesQuery = `mutation RestoreAllRoles { restoreAllRole { status message } }`

const deleteAllRolesPermanentQuery = `mutation DeleteAllRolesPermanent { deleteAllRolePermanent { status message } }`

func (s *RoleGraphQLTestSuite) TestRoleGraphQLLifecycle() {
	// 1. FindByUserId (seeded user carries the default role)
	s.Require().NotZero(s.userID)
	byUser := s.GraphQLOp(s.handler, "findByUserIdRole", findRolesByUserQuery,
		map[string]interface{}{"input": map[string]interface{}{"user_id": s.userID}})
	s.Equal("success", byUser["status"])

	// 2. Create
	created := s.GraphQLOp(s.handler, "createRole", createRoleQuery,
		map[string]interface{}{"input": map[string]interface{}{"name": "Test Role"}})
	s.Equal("success", created["status"])
	s.Equal("Test Role", tests.GQLData(created)["name"])
	s.roleID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.roleID)

	// 3. FindById
	found := s.GraphQLOp(s.handler, "findByIdRole", findRoleByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"role_id": s.roleID}})
	s.Equal("success", found["status"])
	s.Equal(s.roleID, tests.GQLID(tests.GQLData(found), "id"))

	// 4. FindAll
	all := s.GraphQLOp(s.handler, "findAllRole", findAllRolesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 5. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveRole", findActiveRolesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 6. Update
	updated := s.GraphQLOp(s.handler, "updateRole", updateRoleQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.roleID, "name": "Updated Role"}})
	s.Equal("success", updated["status"])
	s.Equal("Updated Role", tests.GQLData(updated)["name"])

	// 7. Trash
	trashed := s.GraphQLOp(s.handler, "trashedRole", trashRoleQuery,
		map[string]interface{}{"input": map[string]interface{}{"role_id": s.roleID}})
	s.Equal("success", trashed["status"])

	// 8. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedRole", findTrashedRolesQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashedList["status"])

	// 9. Restore
	restored := s.GraphQLOp(s.handler, "restoreRole", restoreRoleQuery,
		map[string]interface{}{"input": map[string]interface{}{"role_id": s.roleID}})
	s.Equal("success", restored["status"])

	// 10. DeletePermanent
	s.GraphQLRaw(s.handler, trashRoleQuery,
		map[string]interface{}{"input": map[string]interface{}{"role_id": s.roleID}})
	deleted := s.GraphQLOp(s.handler, "deleteRolePermanent", deleteRolePermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"role_id": s.roleID}})
	s.Equal("success", deleted["status"])

	// 11. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllRole", restoreAllRolesQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 12. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllRolePermanent", deleteAllRolesPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestRoleGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleGraphQLTestSuite))
}
