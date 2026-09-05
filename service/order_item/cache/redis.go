package mencache

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/cache"
)

type Mencache interface {
	OrderItemQueryCache
}

type mencache struct {
	OrderItemQueryCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		OrderItemQueryCache: NewOrderItemQueryCache(cacheStore),
	}
}
