package auth_test

import (
	"net/http"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-apigateway/testhelper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type AuthGraphQLTestSuite struct {
	tests.BaseTestSuite
	handler      http.Handler
	email        string
	password     string
	userID       int
	accessToken  string
	refreshToken string
}

func (s *AuthGraphQLTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupAuthService()

	s.email = "auth.graphql.test@example.com"
	s.password = "password123"

	s.handler = s.NewGraphQLHandler()
}

const registerUserQuery = `mutation RegisterUser($input: RegisterInput!) {
  registerUser(input: $input) { status message data { id firstname lastname email } }
}`

const loginUserQuery = `mutation LoginUser($input: LoginInput!) {
  loginUser(input: $input) { status message data { access_token refresh_token } }
}`

const refreshTokenQuery = `mutation RefreshToken($input: RefreshTokenInput!) {
  refreshToken(input: $input) { status message data { access_token refresh_token } }
}`

const verifyCodeQuery = `mutation VerifyCode($input: VerifyCodeInput!) {
  verifyCode(input: $input) { status message }
}`

const forgotPasswordQuery = `mutation ForgotPassword($input: ForgotPasswordInput!) {
  forgotPassword(input: $input) { status message }
}`

const resetPasswordQuery = `mutation ResetPassword($input: ResetPasswordInput!) {
  resetPassword(input: $input) { status message }
}`

const getMeQuery = `query GetMe($input: GetMeInput!) {
  getMe(input: $input) { status message data { id firstname email } }
}`

func (s *AuthGraphQLTestSuite) TestAuthGraphQLLifecycle() {
	// 1. Register
	registered := s.GraphQLOp(s.handler, "registerUser", registerUserQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"firstname":        "Auth",
			"lastname":         "GraphQL",
			"email":            s.email,
			"password":         s.password,
			"confirm_password": s.password,
		},
	})
	s.Equal("success", registered["status"])
	s.Equal(s.email, tests.GQLData(registered)["email"])
	s.userID = tests.GQLID(tests.GQLData(registered), "id")
	s.Require().NotZero(s.userID)

	// Login requires a verified account; registration leaves it unverified.
	err := s.GormDB().Exec("UPDATE users SET is_verified = true WHERE email = $1", s.email).Error
	s.Require().NoError(err)

	// 2. Login
	loggedIn := s.GraphQLOp(s.handler, "loginUser", loginUserQuery, map[string]interface{}{
		"input": map[string]interface{}{"email": s.email, "password": s.password},
	})
	s.Equal("success", loggedIn["status"])
	tokenData := tests.GQLData(loggedIn)
	s.accessToken, _ = tokenData["access_token"].(string)
	s.refreshToken, _ = tokenData["refresh_token"].(string)
	s.Require().NotEmpty(s.accessToken)
	s.Require().NotEmpty(s.refreshToken)

	// 3. RefreshToken
	refreshed := s.GraphQLOp(s.handler, "refreshToken", refreshTokenQuery, map[string]interface{}{
		"input": map[string]interface{}{"refresh_token": s.refreshToken},
	})
	s.Equal("success", refreshed["status"])
	s.NotEmpty(tests.GQLData(refreshed)["access_token"])

	// 4. GetMe — the resolver reads the caller from the request context, so the
	// handler must be wrapped to emulate AuthMiddleware.
	meHandler := testhelper.WithUserAndToken(s.handler, s.userID, s.accessToken)
	me := s.GraphQLOp(meHandler, "getMe", getMeQuery, map[string]interface{}{
		"input": map[string]interface{}{"access_token": s.accessToken},
	})
	s.Equal("success", me["status"])
	s.Equal(s.email, tests.GQLData(me)["email"])
}

func (s *AuthGraphQLTestSuite) TestAuthPasswordResetGraphQLLifecycle() {
	email := "auth.reset.graphql@example.com"
	password := "password123"

	// 1. Register
	registered := s.GraphQLOp(s.handler, "registerUser", registerUserQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"firstname":        "Reset",
			"lastname":         "GraphQL",
			"email":            email,
			"password":         password,
			"confirm_password": password,
		},
	})
	s.Equal("success", registered["status"])
	userID := tests.GQLID(tests.GQLData(registered), "id")
	s.Require().NotZero(userID)

	// 2. VerifyCode — consumes the code generated at registration
	var code string
	s.Require().NoError(s.GormDB().Raw("SELECT verification_code FROM users WHERE user_id = $1", userID).Scan(&code).Error)
	s.Require().NotEmpty(code)

	verified := s.GraphQLOp(s.handler, "verifyCode", verifyCodeQuery,
		map[string]interface{}{"input": map[string]interface{}{"code": code}})
	s.Equal("success", verified["status"])

	// 3. ForgotPassword — issues a reset token
	forgot := s.GraphQLOp(s.handler, "forgotPassword", forgotPasswordQuery,
		map[string]interface{}{"input": map[string]interface{}{"email": email}})
	s.Equal("success", forgot["status"])

	var token string
	s.Require().NoError(s.GormDB().Raw("SELECT token FROM reset_tokens WHERE user_id = $1", userID).Scan(&token).Error)
	s.Require().NotEmpty(token)

	// 4. ResetPassword
	newPassword := "freshpassword456"
	reset := s.GraphQLOp(s.handler, "resetPassword", resetPasswordQuery, map[string]interface{}{
		"input": map[string]interface{}{
			"reset_token":      token,
			"password":         newPassword,
			"confirm_password": newPassword,
		},
	})
	s.Equal("success", reset["status"])

	// 5. Login with the new password (account already verified via verifyCode)
	loggedIn := s.GraphQLOp(s.handler, "loginUser", loginUserQuery, map[string]interface{}{
		"input": map[string]interface{}{"email": email, "password": newPassword},
	})
	s.Equal("success", loggedIn["status"])
	s.NotEmpty(tests.GQLData(loggedIn)["access_token"])
}

func TestAuthGraphQLSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthGraphQLTestSuite))
}