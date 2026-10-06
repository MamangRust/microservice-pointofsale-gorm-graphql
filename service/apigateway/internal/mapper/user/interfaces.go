package usergraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	userpb "github.com/MamangRust/microservice-point-of-sale-pb/user"
)

type UserGraphqlMapper interface {
	ToGraphqlResponseUser(resp *userpb.ApiResponseUser) *model.APIResponseUserResponse
	ToGraphqlResponseUserDeleteAt(resp *userpb.ApiResponseUserDeleteAt) *model.APIResponseUserResponseDeleteAt
	ToGraphqlResponseUsers(resp *userpb.ApiResponsesUser) *model.APIResponsesUser
	ToGraphqlResponseUserDelete(resp *userpb.ApiResponseUserDelete) *model.APIResponseUserDelete
	ToGraphqlResponseUserAll(resp *userpb.ApiResponseUserAll) *model.APIResponseUserAll
	ToGraphqlResponsePaginationUser(resp *userpb.ApiResponsePaginationUser) *model.APIResponsePaginationUser
	ToGraphqlResponsePaginationUserDeleteAt(resp *userpb.ApiResponsePaginationUserDeleteAt) *model.APIResponsePaginationUserDeleteAt
}
