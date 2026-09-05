package repository

import (
	"github.com/MamangRust/microservice-point-of-sale-pkg/adapter"
	pbproduct "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	"gorm.io/gorm"
)

type Repositories struct {
	CashierQuery         CashierQueryRepository
	MerchantQuery        MerchantQueryRepository
	ProductQuery         ProductQueryRepository
	ProductCommand       ProductCommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	OrderItemQuery       OrderItemQueryRepository
	OrderItemCommand     OrderItemCommandRepository
}

// GuardOptions groups dependency guards by downstream gRPC dependency. Each
// non-nil slice is applied to the matching repository at construction (F6 §8.1
// poin 5 — deadline + circuit breaker + bulkhead per dependency).
type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Product   []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

func NewRepositories(
	DB *gorm.DB,
	cashierClient pbcashier.CashierServiceClient,
	merchantClient pbmerchant.MerchantServiceClient,
	productClient pbproduct.ProductServiceClient,
	orderItemClient pborderitem.OrderItemServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}
	return &Repositories{
		CashierQuery:         NewCashierQueryRepository(cashierClient, g.Cashier...),
		MerchantQuery:        NewMerchantQueryRepository(merchantClient, g.Merchant...),
		ProductQuery:         NewProductQueryRepository(productClient, g.Product...),
		ProductCommand:       NewProductCommandRepository(DB),
		OrderQuery:           NewOrderQueryRepository(DB),
		OrderCommand:         NewOrderCommandRepository(DB),
		OrderItemQuery:       NewOrderItemQueryRepository(orderItemClient, g.OrderItem...),
		OrderItemCommand:     NewOrderItemCommandRepository(DB),
	}
}
