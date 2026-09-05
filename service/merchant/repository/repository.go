package repository

import (
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(
	DB *gorm.DB,
	userClient pbuser.UserServiceClient,
) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB),
		MerchantCommand:         NewMerchantCommandRepository(DB),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB),
		UserQuery:               NewUserQueryRepository(userClient),
	}
}
