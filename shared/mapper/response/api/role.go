package response_api

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
)

type roleResponseMapper struct {
}

func NewRoleResponseMapper() *roleResponseMapper {
	return &roleResponseMapper{}
}

func (s *roleResponseMapper) ToApiResponseRoleAll(pbResponse *pb.ApiResponseRoleAll) *response.ApiResponseRoleAll {
	return &response.ApiResponseRoleAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRoleDelete(pbResponse *pb.ApiResponseRoleDelete) *response.ApiResponseRoleDelete {
	return &response.ApiResponseRoleDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleResponseMapper) ToApiResponseRole(pbResponse *pb.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsesRole(pbResponse *pb.ApiResponsesRole) *response.ApiResponsesRole {
	return &response.ApiResponsesRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponsesRole(pbResponse.Data),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRole(pbResponse *pb.ApiResponsePaginationRole) *response.ApiResponsePaginationRole {
	return &response.ApiResponsePaginationRole{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRole(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleResponseMapper) ToApiResponsePaginationRoleDeleteAt(pbResponse *pb.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt {
	return &response.ApiResponsePaginationRoleDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse.Data),
		Pagination: mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleResponseMapper) mapResponseRole(role *pb.RoleResponse) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRole(roles []*pb.RoleResponse) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}

	return responseRoles
}

func (s *roleResponseMapper) mapResponseRoleDeleteAt(role *pb.RoleResponseDeleteAt) *response.RoleResponseDeleteAt {
	return &response.RoleResponseDeleteAt{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: role.DeletedAt,
	}
}

func (s *roleResponseMapper) mapResponsesRoleDeleteAt(roles []*pb.RoleResponseDeleteAt) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}

	return responseRoles
}

func mapPaginationMeta(s *pbcommon.PaginationMeta) *response.PaginationMeta {
	return &response.PaginationMeta{
		CurrentPage:  int(s.CurrentPage),
		PageSize:     int(s.PageSize),
		TotalRecords: int(s.TotalRecords),
		TotalPages:   int(s.TotalPages),
	}
}
