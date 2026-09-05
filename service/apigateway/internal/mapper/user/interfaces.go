package usergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
)

type UserGraphqlMapper interface {
	ToGraphqlResponseUser(resp *pb.ApiResponseUser) *model.APIResponseUserResponse
	ToGraphqlResponseUserDeleteAt(resp *pb.ApiResponseUserDeleteAt) *model.APIResponseUserResponseDeleteAt
	ToGraphqlResponseUsers(resp *pb.ApiResponsesUser) *model.APIResponsesUser
	ToGraphqlResponseUserDelete(resp *pb.ApiResponseUserDelete) *model.APIResponseUserDelete
	ToGraphqlResponseUserAll(resp *pb.ApiResponseUserAll) *model.APIResponseUserAll
	ToGraphqlResponsePaginationUser(resp *pb.ApiResponsePaginationUser) *model.APIResponsePaginationUser
	ToGraphqlResponsePaginationUserDeleteAt(resp *pb.ApiResponsePaginationUserDeleteAt) *model.APIResponsePaginationUserDeleteAt
}
