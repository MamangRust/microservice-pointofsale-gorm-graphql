package repository

import (
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	UserQuery              UserQueryRepository
	MerchantQuery          MerchantQueryRepository
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
}

func NewRepositories(
	DB *gorm.DB,
	userClient pbuser.UserServiceClient,
	merchantClient pbmerchant.MerchantServiceClient,
) *Repositories {
	return &Repositories{
		UserQuery:              NewUserQueryRepository(userClient),
		MerchantQuery:          NewMerchantQueryRepository(merchantClient),
		CashierQuery:           NewCashierQueryRepository(DB),
		CashierCommand:         NewCashierCommandRepository(DB),
	}
}
