package mapper

import (
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
	commonpb "github.com/MamangRust/microservice-point-of-sale-pb/common"
)

func MapPaginationMeta(s *commonpb.PaginationMeta) *model.PaginationMeta {
	return &model.PaginationMeta{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalRecords: int32(s.TotalRecords),
		TotalPages:   int32(s.TotalPages),
	}
}
