package seeder

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
	"gorm.io/gorm"
)

type Deps struct {
	DB          *gorm.DB
	Ctx         context.Context
	Logger      logger.LoggerInterface
	Hash        hash.HashPassword
}

type Seeder struct {
	User        *userSeeder
	Role        *roleSeeder
	UserRole    *userRoleSeeder
	Cashier     *cashierSeeder
	Category    *categorySeeder
	Product     *productSeeder
	Merchant    *merchantSeeder
	Order       *orderSeeder
	Transaction *transactionSeeder
}

func NewSeeder(deps Deps) *Seeder {
	db := deps.DB
	return &Seeder{
		User:        NewUserSeeder(db, deps.Hash, deps.Ctx, deps.Logger),
		Role:        NewRoleSeeder(db, deps.Ctx, deps.Logger),
		UserRole:    NewUserRoleSeeder(db, deps.Ctx, deps.Logger),
		Merchant:    NewMerchantSeeder(db, deps.Ctx, deps.Logger),
		Cashier:     NewCashierSeeder(db, deps.Ctx, deps.Logger),
		Category:    NewCategorySeeder(db, deps.Ctx, deps.Logger),
		Product:     NewProductSeeder(db, deps.Ctx, deps.Logger),
		Order:       NewOrderSeeder(db, deps.Ctx, deps.Logger),
		Transaction: NewTransactionSeeder(db, deps.Ctx, deps.Logger),
	}
}

func (s *Seeder) Run() error {
	if err := s.seedWithDelay("users", s.User.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("roles", s.Role.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("user_roles", s.UserRole.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchant", s.Merchant.Seed); err != nil {
		return nil
	}

	if err := s.seedWithDelay("cashier", s.Cashier.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("category", s.Category.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("product", s.Product.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("order", s.Order.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("transaction", s.Transaction.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) seedWithDelay(entityName string, seedFunc func() error) error {
	if err := seedFunc(); err != nil {
		return fmt.Errorf("failed to seed %s: %w", entityName, err)
	}
	time.Sleep(seedDelay())
	return nil
}

func seedDelay() time.Duration {
	if raw := os.Getenv("SEED_DELAY_SECONDS"); raw != "" {
		if secs, err := strconv.Atoi(raw); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
	}
	return 30 * time.Second
}

func ptrString(s string) *string {
	return &s
}

func ptrInt32(i int32) *int32 {
	return &i
}
