package repository

import (
    "gorm.io/gorm"
)

type Repositories struct {
    ProductQuery   ProductQueryRepository
    ProductCommand ProductCommandRepository
    CategoryQuery  CategoryQueryRepository
    MerchantQuery  MerchantQueryRepository
}

func NewRepositories(DB *gorm.DB) *Repositories {
    return &Repositories{
        ProductQuery:   NewProductQueryRepository(DB),
        ProductCommand: NewProductCommandRepository(DB),
        CategoryQuery:  NewCategoryQueryRepository(DB),
        MerchantQuery:  NewMerchantQueryRepository(DB),
    }
}
