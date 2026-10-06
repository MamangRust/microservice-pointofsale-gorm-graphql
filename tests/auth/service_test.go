package auth_test

import (
	"context"
	"testing"

	mencache "github.com/MamangRust/microservice-point-of-sale-auth/cache"
	"github.com/MamangRust/microservice-point-of-sale-auth/repository"
	"github.com/MamangRust/microservice-point-of-sale-auth/service"
	"github.com/MamangRust/microservice-point-of-sale-pkg/auth"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	pbrole "github.com/MamangRust/microservice-point-of-sale-pb/role"
	pbuser "github.com/MamangRust/microservice-point-of-sale-pb/user"
	pbuserrole "github.com/MamangRust/microservice-point-of-sale-pb/user_role"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/microservice-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	tests.BaseTestSuite
	authService *service.Service
	email       string
	password    string
}

func (s *AuthServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Register resolves ROLE_ADMIN and persists through the real user service.
	s.SetupUserService()

	userQueryClient := pbuser.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := pbuser.NewUserCommandServiceClient(s.Conns["user"])
	roleClient := pbrole.NewRoleQueryServiceClient(s.Conns["role"])
	userRoleClient := pbuserrole.NewUserRoleServiceClient(s.Conns["role"])
	repos := repository.NewRepositories(s.GormDB(), userQueryClient, userCommandClient, roleClient, userRoleClient)

	hasher := hash.NewHashingPassword()
	mencacheService := mencache.NewMencache(s.GetCacheStore())

	tokenManager, _ := auth.NewManager("mysecretkey")

	s.authService = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencacheService,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})

	s.GormDB().WithContext(s.Ctx).Exec(
		`INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN') ON CONFLICT (role_name) DO NOTHING`)

	s.email = "auth.service.test@example.com"
	s.password = "password123"
}

func (s *AuthServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *AuthServiceTestSuite) TestAuthLifecycle() {
	ctx := context.Background()

	// 1. Register
	regReq := &requests.RegisterRequest{
		FirstName:       "Auth",
		LastName:        "Service",
		Email:           s.email,
		Password:        s.password,
		ConfirmPassword: s.password,
	}

	created, err := s.authService.Register.Register(ctx, regReq)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	s.Equal(s.email, created.Email)

	// 1b. Verify email (login hanya menerima user is_verified = true)
	var verifyCode string
	err = s.GormDB().WithContext(ctx).Raw(
		"SELECT verification_code FROM users WHERE email = ?", s.email).Scan(&verifyCode).Error
	s.Require().NoError(err)

	verified, err := s.authService.PasswordReset.VerifyCode(ctx, verifyCode)
	s.Require().NoError(err)
	s.True(verified)

	// 2. Login
	loginReq := &requests.AuthRequest{
		Email:    s.email,
		Password: s.password,
	}

	tokenRes, err := s.authService.Login.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotNil(tokenRes)
	s.NotEmpty(tokenRes.AccessToken)
	s.NotEmpty(tokenRes.RefreshToken)

	// 3. ForgotPassword
	success, err := s.authService.PasswordReset.ForgotPassword(ctx, s.email)
	s.Require().NoError(err)
	s.True(success)
}

func TestAuthServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthServiceTestSuite))
}
