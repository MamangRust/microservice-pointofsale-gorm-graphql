package merchant_document_cache

import "github.com/MamangRust/microservice-point-of-sale-shared/cache"

type MerchantDocumentMencache interface {
	MerchantDocumentCommandCache
	MerchantDocumentQueryCache
}

type merchantDocumentMencache struct {
	MerchantDocumentCommandCache
	MerchantDocumentQueryCache
}

func NewMerchantDocumentMencache(cacheStore *cache.CacheStore) MerchantDocumentMencache {
	return &merchantDocumentMencache{
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
	}
}
