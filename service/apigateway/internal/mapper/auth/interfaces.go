package authgraphqlmapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	authpb "github.com/MamangRust/microservice-point-of-sale-pb/auth"
)

type AuthGraphqlMapper interface {
	ToGraphqlVerifyCode(res *authpb.ApiResponseVerifyCode) *model.APIResponseVerifyCode
	ToGraphqlForgotPassword(res *authpb.ApiResponseForgotPassword) *model.APIResponseForgotPassword
	ToGraphqlResetPassword(res *authpb.ApiResponseResetPassword) *model.APIResponseResetPassword
	ToGraphqlResponseLogin(res *authpb.ApiResponseLogin) *model.APIResponseLogin
	ToGraphqlResponseRegister(res *authpb.ApiResponseRegister) *model.APIResponseRegister
	ToGraphqlResponseRefreshToken(res *authpb.ApiResponseRefreshToken) *model.APIResponseRefreshToken
	ToGraphqlResponseGetMe(res *authpb.ApiResponseGetMe) *model.APIResponseGetMe
}
