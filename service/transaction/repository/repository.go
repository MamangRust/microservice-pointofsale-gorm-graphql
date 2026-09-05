package repository

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	"gorm.io/gorm"
)

type Repositories struct {
	CashierQuery                 CashierQueryRepository
	MerchantQuery                MerchantQueryRepository
	OrderQuery                   OrderQueryRepository
	OrderItemQuery               OrderItemQueryRepository
	TransactionCommandRepository TransactionCommandRepository
	TransactionQueryRepository   TransactionQueryRepository
}

type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

func NewRepositories(
	db *gorm.DB,
	cashierClient pbcashier.CashierServiceClient,
	merchantClient pbmerchant.MerchantServiceClient,
	orderClient pborder.OrderServiceClient,
	orderItemClient pborderitem.OrderItemServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}
	return &Repositories{
		CashierQuery:                 NewCashierQueryRepository(cashierClient, g.Cashier...),
		MerchantQuery:                NewMerchantQueryRepository(merchantClient, g.Merchant...),
		OrderQuery:                   NewOrderQueryRepository(orderClient, g.Order...),
		OrderItemQuery:               NewOrderItemQueryRepository(orderItemClient, g.OrderItem...),
		TransactionCommandRepository: NewTransactionCommandRepository(db),
		TransactionQueryRepository:   NewTransactionQueryRepository(db),
	}
}
