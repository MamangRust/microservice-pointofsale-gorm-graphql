package auth_test

import (
	"net/http"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerApiTestSuite struct {
	tests.BaseTestSuite
	handler     http.Handler
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupAuthService()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()

	// Seed roles required by auth/user service
	s.GormDB().WithContext(s.Ctx).Exec(`
		INSERT INTO roles (role_name, created_at, updated_at)
		VALUES ('Admin Access 1', current_timestamp, current_timestamp),
		       ('ROLE_ADMIN', current_timestamp, current_timestamp)
		ON CONFLICT (role_name) DO NOTHING
	`)

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)

	s.email = "auth.graphql.handler@test.com"
	s.password = "password123"
}

func (s *AuthHandlerApiTestSuite) Test_AuthLifecycle() {
	// 1. Register
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		registerUser(input: {
			firstname: "Auth", lastname: "GraphQL",
			email: "auth.graphql.handler@test.com",
			password: "password123", confirm_password: "password123"
		}) {
			status message data { id email }
		}
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["registerUser"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.userID = int(data["id"].(float64))
	s.NotZero(s.userID)

	// Verify user
	err = s.GormDB().WithContext(s.Ctx).Exec("UPDATE users SET is_verified = true WHERE email = ?", s.email).Error
	s.Require().NoError(err)

	// 2. Login
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		loginUser(input: { email: "auth.graphql.handler@test.com", password: "password123" }) {
			status message data { access_token }
		}
	}`, nil, "")
	s.Require().NoError(err)
	loginResult := resp.Data["loginUser"].(map[string]interface{})
	s.Equal("success", loginResult["status"])
	loginData := loginResult["data"].(map[string]interface{})
	s.accessToken = loginData["access_token"].(string)
	s.NotEmpty(s.accessToken)
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}
