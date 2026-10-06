package rolegraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	rolepb "github.com/MamangRust/microservice-point-of-sale-pb/role"
)

type RoleGraphqlMapper interface {
	ToGraphqlResponseRole(res *rolepb.ApiResponseRole) *model.APIResponseRole
	ToGraphqlResponseRoleDeleteAt(res *rolepb.ApiResponseRole) *model.APIResponseRoleDeleteAt
	ToGraphqlResponsesRole(res *rolepb.ApiResponsesRole) *model.APIResponsesRole
	ToGraphqlResponseDelete(res *rolepb.ApiResponseRoleDelete) *model.APIResponseRoleDelete
	ToGraphqlResponseAll(res *rolepb.ApiResponseRoleAll) *model.APIResponseRoleAll
	ToGraphqlResponsePaginationRole(res *rolepb.ApiResponsePaginationRole) *model.APIResponsePaginationRole
	ToGraphqlResponsePaginationRoleDeleteAt(res *rolepb.ApiResponsePaginationRoleDeleteAt) *model.APIResponsePaginationRoleDeleteAt
}
