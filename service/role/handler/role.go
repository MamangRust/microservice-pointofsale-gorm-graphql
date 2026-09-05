package handler

import (
	"context"
	"math"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"github.com/MamangRust/microservice-point-of-sale-role/repository"
	"github.com/MamangRust/microservice-point-of-sale-role/service"
	"github.com/MamangRust/microservice-point-of-sale-shared/convert"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors"
	"github.com/MamangRust/microservice-point-of-sale-shared/errors/role_errors"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleGrpc struct {
	pb.UnimplementedRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewRoleHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.RoleServiceServer {
	return &roleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		logger:      logger,
	}
}


func (s *roleHandleGrpc) FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }

	reqService := requests.FindAllRoles{Page: page, PageSize: pageSize, Search: search}
	roles, totalRecords, err := s.roleQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}

	return &pb.ApiResponsePaginationRole{
		Status: "success", Message: "Successfully fetched role records",
		Data: mapRoleResults(roles), Pagination: paginationMeta,
	}, nil
}

func (s *roleHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())
	if roleID <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	role, err := s.roleQuery.FindById(ctx, roleID)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRole{Status: "success", Message: "Successfully fetched role", Data: mapRoleModel(role)}, nil
}

func (s *roleHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	userID := int(req.GetUserId())
	if userID <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponsesRole{Status: "success", Message: "Successfully fetched role by user ID", Data: mapRoleModels(roles)}, nil
}

func (s *roleHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }

	reqService := requests.FindAllRoles{Page: page, PageSize: pageSize, Search: search}
	roles, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, &reqService)
	if err != nil { return nil, errors.ToGrpcError(err) }

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}

	return &pb.ApiResponsePaginationRoleDeleteAt{
		Status: "success", Message: "Successfully fetched active roles",
		Data: mapRoleResultsDeleteAt(roles), Pagination: paginationMeta,
	}, nil
}

func (s *roleHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }

	reqService := requests.FindAllRoles{Page: page, PageSize: pageSize, Search: search}
	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, &reqService)
	if err != nil { return nil, errors.ToGrpcError(err) }

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))
	paginationMeta := &pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(*totalRecords)}

	return &pb.ApiResponsePaginationRoleDeleteAt{
		Status: "success", Message: "Successfully fetched trashed roles",
		Data: mapRoleResultsDeleteAt(roles), Pagination: paginationMeta,
	}, nil
}

func (s *roleHandleGrpc) CreateRole(ctx context.Context, reqPb *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{Name: reqPb.Name}
	if err := req.Validate(); err != nil { return nil, role_errors.ErrGrpcValidateCreateRole }
	role, err := s.roleCommand.CreateRole(ctx, req)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRole{Status: "success", Message: "Successfully created role", Data: mapRoleModel(role)}, nil
}

func (s *roleHandleGrpc) UpdateRole(ctx context.Context, reqPb *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(reqPb.GetId())
	if roleID <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	req := &requests.UpdateRoleRequest{ID: &roleID, Name: reqPb.GetName()}
	if err := req.Validate(); err != nil { return nil, role_errors.ErrGrpcValidateUpdateRole }
	role, err := s.roleCommand.UpdateRole(ctx, req)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRole{Status: "success", Message: "Successfully updated role", Data: mapRoleModel(role)}, nil
}

func (s *roleHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())
	if roleID <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	role, err := s.roleCommand.TrashedRole(ctx, roleID)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRole{Status: "success", Message: "Successfully trashed role", Data: mapRoleModel(role)}, nil
}

func (s *roleHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())
	if roleID <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	role, err := s.roleCommand.RestoreRole(ctx, roleID)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRole{Status: "success", Message: "Successfully restored role", Data: mapRoleModel(role)}, nil
}

func (s *roleHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())
	if id <= 0 { return nil, role_errors.ErrGrpcRoleInvalidId }
	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRoleDelete{Status: "success", Message: "Successfully deleted role permanently"}, nil
}

func (s *roleHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRoleAll{Status: "success", Message: "Successfully restored all roles"}, nil
}

func (s *roleHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)
	if err != nil { return nil, errors.ToGrpcError(err) }
	return &pb.ApiResponseRoleAll{Status: "success", Message: "Successfully deleted all roles permanently"}, nil
}

// Map helpers
func mapRoleModel(role *models.Role) *pb.RoleResponse {
	if role == nil { return nil }
	return &pb.RoleResponse{
		Id: role.RoleID, Name: role.RoleName,
		CreatedAt: convert.FormatTimePtr(role.CreatedAt), UpdatedAt: convert.FormatTimePtr(role.UpdatedAt),
	}
}

func mapRoleModels(roles []*models.Role) []*pb.RoleResponse {
	var res []*pb.RoleResponse
	for _, r := range roles { res = append(res, mapRoleModel(r)) }
	return res
}

func mapRoleResults(roles []*repository.RoleResult) []*pb.RoleResponse {
	var res []*pb.RoleResponse
	for _, r := range roles {
		res = append(res, &pb.RoleResponse{
			Id: r.RoleID, Name: r.RoleName,
			CreatedAt: strVal(r.CreatedAt), UpdatedAt: strVal(r.UpdatedAt),
		})
	}
	return res
}

func mapRoleResultsDeleteAt(roles []*repository.RoleResult) []*pb.RoleResponseDeleteAt {
	var res []*pb.RoleResponseDeleteAt
	for _, r := range roles {
		res = append(res, &pb.RoleResponseDeleteAt{
			Id: r.RoleID, Name: r.RoleName,
			CreatedAt: strVal(r.CreatedAt), UpdatedAt: strVal(r.UpdatedAt), DeletedAt: strVal(r.DeletedAt),
		})
	}
	return res
}

func strVal(s *string) string {
	return convert.StrVal(s)
}
