package orderitem_cache

import "github.com/MamangRust/microservice-point-of-sale-shared/cache"

type OrderItemCache interface {
	OrderItemQueryCache
}

type orderItemCache struct {
	OrderItemQueryCache
}

func NewOrderItemCache(store *cache.CacheStore) OrderItemCache {
	return &orderItemCache{
		OrderItemQueryCache: NewOrderItemQueryCache(store),
	}
}
