package usergraphqlmapper

import (
	graphqlmapper "github.com/MamangRust/microservice-point-of-sale-apigateway/internal/mapper"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	userpb "github.com/MamangRust/microservice-point-of-sale-pb/user"
)

type userGraphqlMapper struct {
}

func NewUserGraphqlMapper() *userGraphqlMapper {
	return &userGraphqlMapper{}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserDelete(res *userpb.ApiResponseUserDelete) *model.APIResponseUserDelete {
	return &model.APIResponseUserDelete{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserAll(res *userpb.ApiResponseUserAll) *model.APIResponseUserAll {
	return &model.APIResponseUserAll{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUser(res *userpb.ApiResponseUser) *model.APIResponseUserResponse {
	return &model.APIResponseUserResponse{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponse(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUserDeleteAt(res *userpb.ApiResponseUserDeleteAt) *model.APIResponseUserResponseDeleteAt {
	return &model.APIResponseUserResponseDeleteAt{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponseDeleteAt(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponseUsers(res *userpb.ApiResponsesUser) *model.APIResponsesUser {
	return &model.APIResponsesUser{
		Status:  res.Status,
		Message: res.Message,
		Data:    u.mapUserResponses(res.Data),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponsePaginationUser(res *userpb.ApiResponsePaginationUser) *model.APIResponsePaginationUser {
	return &model.APIResponsePaginationUser{
		Status:     res.Status,
		Message:    res.Message,
		Data:       u.mapUserResponses(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (u *userGraphqlMapper) ToGraphqlResponsePaginationUserDeleteAt(res *userpb.ApiResponsePaginationUserDeleteAt) *model.APIResponsePaginationUserDeleteAt {
	return &model.APIResponsePaginationUserDeleteAt{
		Status:     res.Status,
		Message:    res.Message,
		Data:       u.mapUserResponsesDeleteAt(res.Data),
		Pagination: graphqlmapper.MapPaginationMeta(res.Pagination),
	}
}

func (u *userGraphqlMapper) mapUserResponse(user *userpb.UserResponse) *model.UserResponse {
	if user == nil {
		return nil
	}
	return &model.UserResponse{
		ID:        int32(user.Id),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userGraphqlMapper) mapUserResponses(users []*userpb.UserResponse) []*model.UserResponse {
	var responses []*model.UserResponse
	for _, user := range users {
		responses = append(responses, u.mapUserResponse(user))
	}
	return responses
}

func (u *userGraphqlMapper) mapUserResponseDeleteAt(user *userpb.UserResponseDeleteAt) *model.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	var deletedAt *string
	if user.DeletedAt != nil {
		deletedAt = &user.DeletedAt.Value
	}

	return &model.UserResponseDeleteAt{
		ID:        int32(user.Id),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (u *userGraphqlMapper) mapUserResponsesDeleteAt(users []*userpb.UserResponseDeleteAt) []*model.UserResponseDeleteAt {
	var responses []*model.UserResponseDeleteAt
	for _, user := range users {
		responses = append(responses, u.mapUserResponseDeleteAt(user))
	}
	return responses
}
