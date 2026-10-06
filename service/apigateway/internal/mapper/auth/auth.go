package authgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	authpb "github.com/MamangRust/microservice-point-of-sale-pb/auth"
	userpb "github.com/MamangRust/microservice-point-of-sale-pb/user"
)

type authGraphqlMapper struct {
}

func NewAuthGraphqlMapper() *authGraphqlMapper {
	return &authGraphqlMapper{}
}

func (s *authGraphqlMapper) ToGraphqlVerifyCode(res *authpb.ApiResponseVerifyCode) *model.APIResponseVerifyCode {
	return &model.APIResponseVerifyCode{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authGraphqlMapper) ToGraphqlForgotPassword(res *authpb.ApiResponseForgotPassword) *model.APIResponseForgotPassword {
	return &model.APIResponseForgotPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authGraphqlMapper) ToGraphqlResetPassword(res *authpb.ApiResponseResetPassword) *model.APIResponseResetPassword {
	return &model.APIResponseResetPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authGraphqlMapper) ToGraphqlResponseLogin(res *authpb.ApiResponseLogin) *model.APIResponseLogin {
	return &model.APIResponseLogin{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseToken(res.Data),
	}
}

func (s *authGraphqlMapper) ToGraphqlResponseRegister(res *authpb.ApiResponseRegister) *model.APIResponseRegister {
	return &model.APIResponseRegister{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseUser(res.Data),
	}
}

func (s *authGraphqlMapper) ToGraphqlResponseRefreshToken(res *authpb.ApiResponseRefreshToken) *model.APIResponseRefreshToken {
	return &model.APIResponseRefreshToken{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseToken(res.Data),
	}
}

func (s *authGraphqlMapper) ToGraphqlResponseGetMe(res *authpb.ApiResponseGetMe) *model.APIResponseGetMe {
	return &model.APIResponseGetMe{
		Status:  res.Status,
		Message: res.Message,
		Data:    s.mapResponseUser(res.Data),
	}
}

func (s *authGraphqlMapper) mapResponseUser(res *userpb.UserResponse) *model.UserResponse {
	return &model.UserResponse{
		ID:        res.Id,
		Firstname: res.Firstname,
		Lastname:  res.Lastname,
		Email:     res.Email,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}

func (s *authGraphqlMapper) mapResponseToken(res *authpb.TokenResponse) *model.TokenResponse {
	return &model.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}
