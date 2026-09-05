package mapper

import (
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/internal/model"
)

func MapPaginationMeta(s *pbcommon.PaginationMeta) *model.PaginationMeta {
	return &model.PaginationMeta{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalRecords: int32(s.TotalRecords),
		TotalPages:   int32(s.TotalPages),
	}
}
