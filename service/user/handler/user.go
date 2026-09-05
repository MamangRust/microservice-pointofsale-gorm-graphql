package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/user_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"github.com/MamangRust/microservice-point-of-sale-user/repository"
	"github.com/MamangRust/microservice-point-of-sale-user/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleGrpc struct {
	pb.UnimplementedUserServiceServer
	userQuery   service.UserQueryService
	userCommand service.UserCommandService
	logger      logger.LoggerInterface
}

func NewUserHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.UserServiceServer {
	return &userHandleGrpc{
		userQuery:   service.UserQuery,
		userCommand: service.UserCommand,
		logger:      logger,
	}
}


func (s *userHandleGrpc) FindAll(ctx context.Context, req *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{Page: page, PageSize: pageSize, Search: search}
	users, totalRecords, err := s.userQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationUser{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       mapUserResults(users),
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) FindById(ctx context.Context, req *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	userID := int(req.GetId())
	if userID <= 0 {
		return nil, user_errors.ErrUserNotFound
	}
	user, err := s.userQuery.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUser{Status: "success", Message: "Successfully fetched user", Data: mapUserModel(user)}, nil
}

func (s *userHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{Page: page, PageSize: pageSize, Search: search}
	users, totalRecords, err := s.userQuery.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active users",
		Data:       mapUserResultsDeleteAt(users),
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{Page: page, PageSize: pageSize, Search: search}
	users, totalRecords, err := s.userQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed users",
		Data:       mapUserResultsDeleteAt(users),
		Pagination: paginationMeta,
	}, nil
}

func (s *userHandleGrpc) Create(ctx context.Context, reqPb *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
	req := &requests.CreateUserRequest{
		FirstName:       reqPb.GetFirstname(),
		LastName:        reqPb.GetLastname(),
		Email:           reqPb.GetEmail(),
		Password:        reqPb.GetPassword(),
		ConfirmPassword: reqPb.GetConfirmPassword(),
	}
	user, err := s.userCommand.CreateUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUser{Status: "success", Message: "Successfully created user", Data: mapUserModel(user)}, nil
}

func (s *userHandleGrpc) Update(ctx context.Context, reqPb *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
	userID := int(reqPb.GetId())
	req := &requests.UpdateUserRequest{
		UserID:          &userID,
		FirstName:       reqPb.GetFirstname(),
		LastName:        reqPb.GetLastname(),
		Email:           reqPb.GetEmail(),
		Password:        reqPb.GetPassword(),
		ConfirmPassword: reqPb.GetConfirmPassword(),
	}
	user, err := s.userCommand.UpdateUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUser{Status: "success", Message: "Successfully updated user", Data: mapUserModel(user)}, nil
}

func (s *userHandleGrpc) TrashedUser(ctx context.Context, req *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	userID := int(req.GetId())
	if userID <= 0 {
		return nil, user_errors.ErrUserNotFound
	}
	user, err := s.userCommand.TrashedUser(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    mapUserModelDeleteAt(user),
	}, nil
}

func (s *userHandleGrpc) RestoreUser(ctx context.Context, req *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	userID := int(req.GetId())
	if userID <= 0 {
		return nil, user_errors.ErrUserNotFound
	}
	user, err := s.userCommand.RestoreUser(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    mapUserModelDeleteAt(user),
	}, nil
}

func (s *userHandleGrpc) DeleteUserPermanent(ctx context.Context, req *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	userID := int(req.GetId())
	if userID <= 0 {
		return nil, user_errors.ErrUserNotFound
	}
	_, err := s.userCommand.DeleteUserPermanent(ctx, userID)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUserDelete{Status: "success", Message: "Successfully deleted user permanently"}, nil
}

func (s *userHandleGrpc) RestoreAllUser(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userCommand.RestoreAllUser(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUserAll{Status: "success", Message: "Successfully restored all users"}, nil
}

func (s *userHandleGrpc) DeleteAllUserPermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	_, err := s.userCommand.DeleteAllUserPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseUserAll{Status: "success", Message: "Successfully deleted all users permanently"}, nil
}

// Map helpers
func mapUserModel(user *models.User) *pb.UserResponse {
	if user == nil {
		return nil
	}
	return &pb.UserResponse{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: convert.FormatTimePtr(user.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(user.UpdatedAt),
	}
}

func mapUserModelDeleteAt(user *models.User) *pb.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	return &pb.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: convert.FormatTimePtr(user.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(user.UpdatedAt),
		DeletedAt: convert.TimeToWrappers(user.DeletedAt),
	}
}

func mapUserResults(users []*repository.UserResult) []*pb.UserResponse {
	var res []*pb.UserResponse
	for _, u := range users {
		res = append(res, &pb.UserResponse{
			Id:        u.UserID,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Email:     u.Email,
			CreatedAt: userStrVal(u.CreatedAt),
			UpdatedAt: userStrVal(u.UpdatedAt),
		})
	}
	return res
}

func mapUserResultsDeleteAt(users []*repository.UserResult) []*pb.UserResponseDeleteAt {
	var res []*pb.UserResponseDeleteAt
	for _, u := range users {
		deletedAt := convert.StrValToWrappers(u.DeletedAt)
		res = append(res, &pb.UserResponseDeleteAt{
			Id:        u.UserID,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Email:     u.Email,
			CreatedAt: userStrVal(u.CreatedAt),
			UpdatedAt: userStrVal(u.UpdatedAt),
			DeletedAt: deletedAt,
		})
	}
	return res
}

func userStrVal(s *string) string {
	return convert.StrVal(s)
}
