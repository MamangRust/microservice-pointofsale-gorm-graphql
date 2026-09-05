package role_test

import (
	"context"
	"testing"

	role_cache "github.com/MamangRust/microservice-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/microservice-point-of-sale-role/handler"
	"github.com/MamangRust/microservice-point-of-sale-role/repository"
	"github.com/MamangRust/microservice-point-of-sale-role/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleGapiTestSuite struct {
	tests.BaseTestSuite
	client pb.RoleServiceClient
}

func (s *RoleGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	roleQueries := s.GormDB()

	// Role dependencies
	mencache := role_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(roleQueries)
	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	// Handler
	handler := role_handler.NewHandler(&role_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterRoleServiceServer(server, handler.Role)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pb.NewRoleServiceClient(conn)
}

func (s *RoleGapiTestSuite) TestRoleGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.client.CreateRole(ctx, &pb.CreateRoleRequest{
		Name: "Gapi Role",
	})
	s.Require().NotNil(createRes)
	roleID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.client.FindByIdRole(ctx, &pb.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)
	s.Equal("Gapi Role", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.client.FindAllRole(ctx, &pb.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pb.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.client.UpdateRole(ctx, &pb.UpdateRoleRequest{
		Id:   roleID,
		Name: "Gapi Role Updated",
	})
	s.Require().NoError(err)
	s.Equal("Gapi Role Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.client.TrashedRole(ctx, &pb.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pb.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.client.RestoreRole(ctx, &pb.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.client.TrashedRole(ctx, &pb.FindByIdRoleRequest{RoleId: roleID})
	_, err = s.client.DeleteRolePermanent(ctx, &pb.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.client.RestoreAllRole(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.client.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestRoleGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleGapiTestSuite))
}
