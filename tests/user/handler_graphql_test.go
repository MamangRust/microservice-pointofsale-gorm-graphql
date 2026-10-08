package user_test

import (
	"net/http"
	"testing"

	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type UserGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	userID  int
}

func (s *UserGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()

	// The user service assigns the "Admin Access 1" role on create, so it must
	// exist before the GraphQL lifecycle runs.
	err := s.GormDB().Exec(`
		INSERT INTO roles (role_name, created_at, updated_at)
		VALUES ('Admin Access 1', current_timestamp, current_timestamp),
		       ('ROLE_ADMIN', current_timestamp, current_timestamp)
		ON CONFLICT (role_name) DO NOTHING
	`).Error
	s.Require().NoError(err)

	s.handler = s.NewGraphQLHandler()
}

const createUserQuery = `mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) { status message data { id firstname lastname email } }
}`

const updateUserQuery = `mutation UpdateUser($input: UpdateUserInput!) {
  updateUser(input: $input) { status message data { id firstname lastname email } }
}`

const findUserByIdQuery = `query FindUserById($input: FindByIdUserInput!) {
  findByIdUser(input: $input) { status message data { id firstname email } }
}`

const findAllUsersQuery = `query FindAllUsers($input: FindAllUserInput) {
  findAllUsers(input: $input) { status message data { id email } pagination { total_records } }
}`

const findActiveUsersQuery = `query FindActiveUsers($input: FindAllUserInput) {
  findByActiveUsers(input: $input) { status message data { id email } pagination { total_records } }
}`

const findTrashedUsersQuery = `query FindTrashedUsers($input: FindAllUserInput) {
  findByTrashedUsers(input: $input) { status message data { id email } pagination { total_records } }
}`

const trashUserQuery = `mutation TrashUser($input: FindByIdUserInput!) {
  trashedUser(input: $input) { status message data { id email deleted_at } }
}`

const restoreUserQuery = `mutation RestoreUser($input: FindByIdUserInput!) {
  restoreUser(input: $input) { status message data { id email deleted_at } }
}`

const deleteUserPermanentQuery = `mutation DeleteUserPermanent($input: FindByIdUserInput!) {
  deleteUserPermanent(input: $input) { status message }
}`

const restoreAllUsersQuery = `mutation RestoreAllUsers { restoreAllUser { status message } }`

const deleteAllUsersPermanentQuery = `mutation DeleteAllUsersPermanent { deleteAllUserPermanent { status message } }`

func (s *UserGraphQLTestSuite) TestUserGraphQLLifecycle() {
	// 1. Create
	created := s.GraphQLOp(s.handler, "createUser", createUserQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"firstname":        "Test",
			"lastname":         "User",
			"email":            "user.graphql.test@example.com",
			"password":         "password123",
			"confirm_password": "password123",
		},
	})
	s.Equal("success", created["status"])
	s.Equal("user.graphql.test@example.com", tests.GQLData(created)["email"])
	s.userID = tests.GQLID(tests.GQLData(created), "id")
	s.Require().NotZero(s.userID)

	// 2. FindById
	found := s.GraphQLOp(s.handler, "findByIdUser", findUserByIdQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.userID}})
	s.Equal("success", found["status"])
	s.Equal(s.userID, tests.GQLID(tests.GQLData(found), "id"))

	// 3. FindAll
	all := s.GraphQLOp(s.handler, "findAllUsers", findAllUsersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", all["status"])

	// 4. FindByActive
	active := s.GraphQLOp(s.handler, "findByActiveUsers", findActiveUsersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", active["status"])

	// 5. Update
	updated := s.GraphQLOp(s.handler, "updateUser", updateUserQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"id":               s.userID,
			"firstname":        "Updated",
			"lastname":         "User",
			"email":            "user.graphql.updated@example.com",
			"password":         "newpassword123",
			"confirm_password": "newpassword123",
		},
	})
	s.Equal("success", updated["status"])
	s.Equal("Updated", tests.GQLData(updated)["firstname"])

	// 6. Trash
	trashed := s.GraphQLOp(s.handler, "trashedUser", trashUserQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.userID}})
	s.Equal("success", trashed["status"])

	// 7. FindByTrashed
	trashedList := s.GraphQLOp(s.handler, "findByTrashedUsers", findTrashedUsersQuery,
		map[string]interface{}{"input": map[string]interface{}{"page": 1, "page_size": 10}})
	s.Equal("success", trashedList["status"])

	// 8. Restore
	restored := s.GraphQLOp(s.handler, "restoreUser", restoreUserQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.userID}})
	s.Equal("success", restored["status"])

	// 9. DeletePermanent
	s.GraphQLRaw(s.handler, trashUserQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.userID}})
	deleted := s.GraphQLOp(s.handler, "deleteUserPermanent", deleteUserPermanentQuery,
		map[string]interface{}{"input": map[string]interface{}{"id": s.userID}})
	s.Equal("success", deleted["status"])

	// 10. RestoreAll
	restoredAll := s.GraphQLOp(s.handler, "restoreAllUser", restoreAllUsersQuery, nil)
	s.Equal("success", restoredAll["status"])

	// 11. DeleteAll
	deletedAll := s.GraphQLOp(s.handler, "deleteAllUserPermanent", deleteAllUsersPermanentQuery, nil)
	s.Equal("success", deletedAll["status"])
}

func TestUserGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserGraphQLTestSuite))
}
