package order_cache

import "github.com/MamangRust/microservice-point-of-sale-shared/cache"

type OrderMencache interface {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

type orderMencache struct {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

func NewOrderMencache(store *cache.CacheStore) OrderMencache {
	return &orderMencache{
		OrderQueryCache:           NewOrderQueryCache(store),
		OrderCommandCache:         NewOrderCommandCache(store),
		OrderStatsCache:           NewOrderStatsCache(store),
		OrderStatsByMerchantCache: NewOrderStatsByMerchantCache(store),
	}
}
