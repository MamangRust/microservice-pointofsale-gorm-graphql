package user_test

import (
	"context"
	"testing"

	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	gapi "github.com/MamangRust/microservice-point-of-sale-user/handler"
	"github.com/MamangRust/microservice-point-of-sale-user/repository"
	"github.com/MamangRust/microservice-point-of-sale-user/service"

	user_cache "github.com/MamangRust/microservice-point-of-sale-user/cache"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGapiTestSuite struct {
	tests.BaseTestSuite
	client      pb.UserServiceClient
	queryClient pb.UserServiceClient
	userID      int
}

func (s *UserGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()

	// Seed default role required by user service CreateUser
	s.GormDB().WithContext(s.Ctx).Exec(
		`INSERT INTO roles (role_name, created_at, updated_at)
		 VALUES ('Admin Access 1', current_timestamp, current_timestamp)
		 ON CONFLICT (role_name) DO NOTHING`)

	userQueries := s.GormDB()
	repos := repository.NewRepositories(userQueries)

	log, _ := logger.NewLogger("test")
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Hash:          hasher,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userHandler.User)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pb.NewUserServiceClient(conn)
	s.queryClient = pb.NewUserServiceClient(conn)
}

func (s *UserGapiTestSuite) TestUserGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateUserRequest{
		Firstname:       "Gapi",
		Lastname:        "User",
		Email:           "gapi.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	res, err := s.client.Create(ctx, createReq)
	s.Require().NoError(err)
	s.Equal(createReq.Email, res.Data.Email)
	userID := res.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pb.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)
	s.Equal(userID, getRes.Data.Id)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, _ := s.client.Update(ctx, &pb.UpdateUserRequest{
		Id:              userID,
		Firstname:       "GapiUpdated",
		Lastname:        "UserUpdated",
		Email:           "gapi.updated@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	s.Equal("GapiUpdated", updateRes.Data.Firstname)

	// 6. Trash
	_, err = s.client.TrashedUser(ctx, &pb.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.client.RestoreUser(ctx, &pb.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.client.TrashedUser(ctx, &pb.FindByIdUserRequest{Id: userID})
	_, err = s.client.DeleteUserPermanent(ctx, &pb.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.client.RestoreAllUser(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestUserGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserGapiTestSuite))
}
