package product_test

import (
	"context"
	"testing"

	prod_cache "github.com/MamangRust/microservice-point-of-sale-product/cache"
	prod_handler "github.com/MamangRust/microservice-point-of-sale-product/handler"
	prod_repo "github.com/MamangRust/microservice-point-of-sale-product/repository"
	prod_service "github.com/MamangRust/microservice-point-of-sale-product/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
	"github.com/MamangRust/microservice-point-of-sale-shared/observability"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductGapiTestSuite struct {
	tests.BaseTestSuite
	client pb.ProductServiceClient
}

func (s *ProductGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	productQueries := s.GormDB()

	// Product dependencies
	mencache := prod_cache.NewMencache(cacheStore)
	repos := prod_repo.NewRepositories(productQueries)
	svc := prod_service.NewService(&prod_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
		Ctx:           context.Background(),
	})

	// Handler
	handler := prod_handler.NewHandler(&prod_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, handler.Product)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pb.NewProductServiceClient(conn)
}

func (s *ProductGapiTestSuite) TestProductGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create
	createRes, err := s.client.Create(ctx, &pb.CreateProductRequest{
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item",
		Description:  "GAPI Description",
		Price:        1000,
		CountInStock: 10,
		Brand:        "GAPI Brand",
		Weight:       100,
		ImageProduct: "gapi.jpg",
	})
	s.Require().NotNil(createRes)
	prodID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.client.FindById(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)
	s.Equal("GAPI Item", getRes.Data.Name)

	// 4. FindAll
	allRes, err := s.client.FindAll(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.client.Update(ctx, &pb.UpdateProductRequest{
		ProductId:    prodID,
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item Updated",
		Description:  "Updated Description",
		Price:        2000,
		CountInStock: 20,
		Brand:        "Updated Brand",
		Weight:       200,
		ImageProduct: "gapi-updated.jpg",
	})
	s.Require().NoError(err)
	s.Equal("GAPI Item Updated", updateRes.Data.Name)

	// 7. Trash
	_, err = s.client.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.client.RestoreProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.client.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	_, err = s.client.DeleteProductPermanent(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.client.RestoreAllProduct(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.client.DeleteAllProductPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestProductGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGapiTestSuite))
}
