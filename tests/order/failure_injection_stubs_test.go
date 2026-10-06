package order_test

import (
	"context"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database/models"
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/requests"
)

// stubCashierRepo implements repository.CashierQueryRepository. In the
// failure-injection suite it is never reached (the merchant dependency fails
// first), so it returns a minimal valid row to keep the flow honest if the
// ordering ever changes.
type stubCashierRepo struct {
	cashierID int
}

func (s *stubCashierRepo) FindById(ctx context.Context, cashierID int) (*models.Cashier, error) {
	return &models.Cashier{CashierID: int32(cashierID), MerchantID: int32(s.cashierID), Name: "stub"}, nil
}

// stubProductRepo implements repository.ProductQueryRepository. Never reached
// in the failure-injection suite (merchant fails first).
type stubProductRepo struct{}

func (s *stubProductRepo) FindById(ctx context.Context, productID int) (*models.Product, error) {
	return &models.Product{ProductID: int32(productID), Price: 10000, CountInStock: 100}, nil
}

// stubProductCommandRepo implements repository.ProductCommandRepository. Never
// reached in the failure-injection suite (merchant fails first).
type stubProductCommandRepo struct{}

func (s *stubProductCommandRepo) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*models.Product, error) {
	return &models.Product{ProductID: int32(productID)}, nil
}

func (s *stubProductCommandRepo) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*models.Product, error) {
	return &models.Product{ProductID: int32(productID)}, nil
}

// stubOrderItemQueryRepo implements repository.OrderItemQueryRepository. Never
// reached in the failure-injection suite (merchant fails first).
type stubOrderItemQueryRepo struct{}

func (s *stubOrderItemQueryRepo) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error) {
	return nil, nil
}

func (s *stubOrderItemQueryRepo) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	total := int32(0)
	return &total, nil
}

// stubOrderItemCommandRepo implements repository.OrderItemCommandRepository.
// Never reached in the failure-injection suite (merchant fails first).
type stubOrderItemCommandRepo struct{}

func (s *stubOrderItemCommandRepo) DeleteOrderItem(ctx context.Context, orderID int) error {
	return nil
}

func (s *stubOrderItemCommandRepo) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	return &models.OrderItem{}, nil
}

func (s *stubOrderItemCommandRepo) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	return &models.OrderItem{}, nil
}
