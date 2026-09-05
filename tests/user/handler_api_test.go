package user_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	userID  int
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	// Seed roles required by user service CreateUser
	s.GormDB().WithContext(s.Ctx).Exec(`
		INSERT INTO roles (role_name, created_at, updated_at) 
		VALUES ('Admin Access 1', current_timestamp, current_timestamp),
		       ('ROLE_ADMIN', current_timestamp, current_timestamp)
		ON CONFLICT (role_name) DO NOTHING
	`)

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		createUser(input: {
			firstname: "Handler", lastname: "User",
			email: "handler.user.graphql@test.com",
			password: "password123", confirm_password: "password123"
		}) {
			status message data { id email }
		}
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createUser"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.userID = int(data["id"].(float64))
	s.NotZero(s.userID)

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllUsers(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllUsers"].(map[string]interface{})["status"])

	// 3. FindById
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByIdUser(input: { id: `+strconv.Itoa(s.userID)+` }) { status message data { id email } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdUser"].(map[string]interface{})["status"])

	// 4. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveUsers(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveUsers"].(map[string]interface{})["status"])

	// 5. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		updateUser(input: { id: `+strconv.Itoa(s.userID)+`, firstname: "Updated", lastname: "User", email: "handler.user.graphql@test.com", password: "password123", confirm_password: "password123" }) {
			status message data { id firstname }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateUser"].(map[string]interface{})["status"])

	// 6. Trash
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedUser(input: { id: `+strconv.Itoa(s.userID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedUser"].(map[string]interface{})["status"])

	// 7. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedUsers(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedUsers"].(map[string]interface{})["status"])

	// 8. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreUser(input: { id: `+strconv.Itoa(s.userID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreUser"].(map[string]interface{})["status"])

	// 9. DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteUserPermanent(input: { id: `+strconv.Itoa(s.userID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteUserPermanent"].(map[string]interface{})["status"])

	// 10. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllUser { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllUser"].(map[string]interface{})["status"])

	// 11. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllUserPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllUserPermanent"].(map[string]interface{})["status"])
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
