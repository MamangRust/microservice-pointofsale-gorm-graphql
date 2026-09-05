package rolegraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
)

type RoleGraphqlMapper interface {
	ToGraphqlResponseRole(res *pb.ApiResponseRole) *model.APIResponseRole
	ToGraphqlResponseRoleDeleteAt(res *pb.ApiResponseRole) *model.APIResponseRoleDeleteAt
	ToGraphqlResponsesRole(res *pb.ApiResponsesRole) *model.APIResponsesRole
	ToGraphqlResponseDelete(res *pb.ApiResponseRoleDelete) *model.APIResponseRoleDelete
	ToGraphqlResponseAll(res *pb.ApiResponseRoleAll) *model.APIResponseRoleAll
	ToGraphqlResponsePaginationRole(res *pb.ApiResponsePaginationRole) *model.APIResponsePaginationRole
	ToGraphqlResponsePaginationRoleDeleteAt(res *pb.ApiResponsePaginationRoleDeleteAt) *model.APIResponsePaginationRoleDeleteAt
}
