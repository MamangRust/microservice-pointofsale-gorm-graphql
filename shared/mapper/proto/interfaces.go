package protomapper

import (
	"github.com/MamangRust/microservice-point-of-sale-shared/domain/response"
	pb "github.com/MamangRust/microservice-pointofsale-grpc/pb"
	pbcashier "github.com/MamangRust/microservice-pointofsale-grpc/pb/cashier"
	pbcategory "github.com/MamangRust/microservice-pointofsale-grpc/pb/category"
	pbcommon "github.com/MamangRust/microservice-pointofsale-grpc/pb/common"
	pbmerchant "github.com/MamangRust/microservice-pointofsale-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-pointofsale-grpc/pb/order"
	pborderitem "github.com/MamangRust/microservice-pointofsale-grpc/pb/order_item"
	pbproduct "github.com/MamangRust/microservice-pointofsale-grpc/pb/product"
	pbrole "github.com/MamangRust/microservice-pointofsale-grpc/pb/role"
	pbstats "github.com/MamangRust/microservice-pointofsale-grpc/pb/stats"
	pbtransaction "github.com/MamangRust/microservice-pointofsale-grpc/pb/transaction"
	pbuser "github.com/MamangRust/microservice-pointofsale-grpc/pb/user"
)

type AuthProtoMapper interface {
	ToProtoResponseVerifyCode(status string, message string) *pb.ApiResponseVerifyCode
	ToProtoResponseForgotPassword(status string, message string) *pb.ApiResponseForgotPassword
	ToProtoResponseResetPassword(status string, message string) *pb.ApiResponseResetPassword
	ToProtoResponseLogin(status string, message string, response *response.TokenResponse) *pb.ApiResponseLogin
	ToProtoResponseRegister(status string, message string, response *response.UserResponse) *pb.ApiResponseRegister
	ToProtoResponseRefreshToken(status string, message string, response *response.TokenResponse) *pb.ApiResponseRefreshToken
	ToProtoResponseGetMe(status string, message string, response *response.UserResponse) *pb.ApiResponseGetMe
}

type UserProtoMapper interface {
	ToProtoResponseUserDeleteAt(status string, message string, pbResponse *response.UserResponseDeleteAt) *pbuser.ApiResponseUserDeleteAt
	ToProtoResponsesUser(status string, message string, pbResponse []*response.UserResponse) *pbuser.ApiResponsesUser
	ToProtoResponseUser(status string, message string, pbResponse *response.UserResponse) *pbuser.ApiResponseUser
	ToProtoResponseUserDelete(status string, message string) *pbuser.ApiResponseUserDelete
	ToProtoResponseUserAll(status string, message string) *pbuser.ApiResponseUserAll
	ToProtoResponsePaginationUserDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponseDeleteAt) *pbuser.ApiResponsePaginationUserDeleteAt
	ToProtoResponsePaginationUser(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.UserResponse) *pbuser.ApiResponsePaginationUser
}

type RoleProtoMapper interface {
	ToProtoResponseRoleAll(status string, message string) *pbrole.ApiResponseRoleAll
	ToProtoResponseRoleDelete(status string, message string) *pbrole.ApiResponseRoleDelete
	ToProtoResponseRole(status string, message string, pbResponse *response.RoleResponse) *pbrole.ApiResponseRole
	ToProtoResponsesRole(status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsesRole
	ToProtoResponsePaginationRole(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponse) *pbrole.ApiResponsePaginationRole
	ToProtoResponsePaginationRoleDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, pbResponse []*response.RoleResponseDeleteAt) *pbrole.ApiResponsePaginationRoleDeleteAt
}

type CategoryProtoMapper interface {
	ToProtoResponseMonthlyTotalPrice(status string, message string, row []*response.CategoriesMonthlyTotalPriceResponse) *pbstats.ApiResponseCategoryMonthlyTotalPrice
	ToProtoResponseYearlyTotalPrice(status string, message string, row []*response.CategoriesYearlyTotalPriceResponse) *pbstats.ApiResponseCategoryYearlyTotalPrice
	ToProtoResponseCategoryMonthlyPrice(status string, message string, row []*response.CategoryMonthPriceResponse) *pbstats.ApiResponseCategoryMonthPrice
	ToProtoResponseCategoryYearlyPrice(status string, message string, row []*response.CategoryYearPriceResponse) *pbstats.ApiResponseCategoryYearPrice

	ToProtoResponsesCategory(status string, message string, pbResponse []*response.CategoryResponse) *pbcategory.ApiResponsesCategory
	ToProtoResponseCategoryDeleteAt(status string, message string, pbResponse *response.CategoryResponseDeleteAt) *pbcategory.ApiResponseCategoryDeleteAt

	ToProtoResponseCategoryAll(status string, message string) *pbcategory.ApiResponseCategoryAll
	ToProtoResponseCategory(status string, message string, pbResponse *response.CategoryResponse) *pbcategory.ApiResponseCategory
	ToProtoResponseCategoryDelete(status string, message string) *pbcategory.ApiResponseCategoryDelete
	ToProtoResponsePaginationCategoryDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponseDeleteAt) *pbcategory.ApiResponsePaginationCategoryDeleteAt
	ToProtoResponsePaginationCategory(pagination *pbcommon.PaginationMeta, status string, message string, categories []*response.CategoryResponse) *pbcategory.ApiResponsePaginationCategory
}

type CashierProtoMapper interface {
	ToProtoMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthTotalSales) *pbstats.ApiResponseCashierMonthlyTotalSales
	ToProtoYearlyTotalSales(status, message string, row []*response.CashierResponseYearTotalSales) *pbstats.ApiResponseCashierYearlyTotalSales

	ToProtoResponseMonthlyTotalSales(status, message string, row []*response.CashierResponseMonthSales) *pbstats.ApiResponseCashierMonthSales
	ToProtoResponseYearlyTotalSales(status, message string, row []*response.CashierResponseYearSales) *pbstats.ApiResponseCashierYearSales

	ToProtoResponseCashier(status string, message string, pbResponse *response.CashierResponse) *pbcashier.ApiResponseCashier
	ToProtoResponseCashierDeleteAt(status string, message string, pbResponse *response.CashierResponseDeleteAt) *pbcashier.ApiResponseCashierDeleteAt
	ToProtoResponsesCashier(status string, message string, pbResponse []*response.CashierResponse) *pbcashier.ApiResponsesCashier
	ToProtoResponseCashierDelete(status string, message string) *pbcashier.ApiResponseCashierDelete
	ToProtoResponseCashierAll(status string, message string) *pbcashier.ApiResponseCashierAll
	ToProtoResponsePaginationCashierDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponseDeleteAt) *pbcashier.ApiResponsePaginationCashierDeleteAt
	ToProtoResponsePaginationCashier(pagination *pbcommon.PaginationMeta, status string, message string, users []*response.CashierResponse) *pbcashier.ApiResponsePaginationCashier
}

type MerchantProtoMapper interface {
	ToProtoResponseMerchant(status string, message string, pbResponse *response.MerchantResponse) *pbmerchant.ApiResponseMerchant
	ToProtoResponseMerchantDeleteAt(status string, message string, pbResponse *response.MerchantResponseDeleteAt) *pbmerchant.ApiResponseMerchantDeleteAt

	ToProtoResponsesMerchant(status string, message string, pbResponse []*response.MerchantResponse) *pbmerchant.ApiResponsesMerchant
	ToProtoResponseMerchantDelete(status string, message string) *pbmerchant.ApiResponseMerchantDelete
	ToProtoResponseMerchantAll(status string, message string) *pbmerchant.ApiResponseMerchantAll
	ToProtoResponsePaginationMerchantDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pbmerchant.ApiResponsePaginationMerchantDeleteAt
	ToProtoResponsePaginationMerchant(pagination *pbcommon.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pbmerchant.ApiResponsePaginationMerchant
}

type MerchantDocumentProtoMapper interface {
	ToProtoResponseMerchantDocument(status string, message string, doc *response.MerchantDocumentResponse) *pbmerchant.ApiResponseMerchantDocument
	ToProtoResponsesMerchantDocument(status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant.ApiResponsesMerchantDocument

	ToProtoResponsePaginationMerchantDocument(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponse) *pbmerchant.ApiResponsePaginationMerchantDocument
	ToProtoResponsePaginationMerchantDocumentDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, docs []*response.MerchantDocumentResponseDeleteAt) *pbmerchant.ApiResponsePaginationMerchantDocumentAt

	ToProtoResponseMerchantDocumentDelete(status string, message string) *pbmerchant.ApiResponseMerchantDocumentDelete

	ToProtoResponseMerchantDocumentAll(status string, message string) *pbmerchant.ApiResponseMerchantDocumentAll
}

type OrderItemProtoMapper interface {
	ToProtoResponseOrderItem(status string, message string, pbResponse *response.OrderItemResponse) *pborderitem.ApiResponseOrderItem
	ToProtoResponsesOrderItem(status string, message string, pbResponse []*response.OrderItemResponse) *pborderitem.ApiResponsesOrderItem
	ToProtoResponseOrderItemDelete(status string, message string) *pborderitem.ApiResponseOrderItemDelete
	ToProtoResponseOrderItemAll(status string, message string) *pborderitem.ApiResponseOrderItemAll
	ToProtoResponsePaginationOrderItemDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponseDeleteAt) *pborderitem.ApiResponsePaginationOrderItemDeleteAt
	ToProtoResponsePaginationOrderItem(pagination *pbcommon.PaginationMeta, status string, message string, orderItems []*response.OrderItemResponse) *pborderitem.ApiResponsePaginationOrderItem
}

type OrderProtoMapper interface {
	ToProtoResponseMonthlyTotalRevenue(status string, message string, row []*response.OrderMonthlyTotalRevenueResponse) *pbstats.ApiResponseOrderMonthlyTotalRevenue
	ToProtoResponseYearlyTotalRevenue(status string, message string, row []*response.OrderYearlyTotalRevenueResponse) *pbstats.ApiResponseOrderYearlyTotalRevenue

	ToProtoResponseMonthlyRevenue(status string, message string, row []*response.OrderMonthlyResponse) *pbstats.ApiResponseOrderMonthly
	ToProtoResponseYearlyRevenue(status string, message string, row []*response.OrderYearlyResponse) *pbstats.ApiResponseOrderYearly

	ToProtoResponseOrder(status string, message string, pbResponse *response.OrderResponse) *pborder.ApiResponseOrder
	ToProtoResponseOrderDeleteAt(status string, message string, pbResponse *response.OrderResponseDeleteAt) *pborder.ApiResponseOrderDeleteAt
	ToProtoResponsesOrder(status string, message string, pbResponse []*response.OrderResponse) *pborder.ApiResponsesOrder
	ToProtoResponseOrderDelete(status string, message string) *pborder.ApiResponseOrderDelete
	ToProtoResponseOrderAll(status string, message string) *pborder.ApiResponseOrderAll
	ToProtoResponsePaginationOrderDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponseDeleteAt) *pborder.ApiResponsePaginationOrderDeleteAt
	ToProtoResponsePaginationOrder(pagination *pbcommon.PaginationMeta, status string, message string, orders []*response.OrderResponse) *pborder.ApiResponsePaginationOrder
}

type ProductProtoMapper interface {
	ToProtoResponseProduct(status string, message string, pbResponse *response.ProductResponse) *pbproduct.ApiResponseProduct
	ToProtoResponseProductDeleteAt(status string, message string, pbResponse *response.ProductResponseDeleteAt) *pbproduct.ApiResponseProductDeleteAt

	ToProtoResponsesProduct(status string, message string, pbResponse []*response.ProductResponse) *pbproduct.ApiResponsesProduct
	ToProtoResponseProductDelete(status string, message string) *pbproduct.ApiResponseProductDelete
	ToProtoResponseProductAll(status string, message string) *pbproduct.ApiResponseProductAll
	ToProtoResponsePaginationProductDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponseDeleteAt) *pbproduct.ApiResponsePaginationProductDeleteAt
	ToProtoResponsePaginationProduct(pagination *pbcommon.PaginationMeta, status string, message string, products []*response.ProductResponse) *pbproduct.ApiResponsePaginationProduct
}

type TransactionProtoMapper interface {
	ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbstats.ApiResponseTransactionMonthAmountSuccess
	ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbstats.ApiResponseTransactionYearAmountSuccess
	ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbstats.ApiResponseTransactionMonthAmountFailed
	ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbstats.ApiResponseTransactionYearAmountFailed
	ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbstats.ApiResponseTransactionMonthPaymentMethod
	ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbstats.ApiResponseTransactionYearPaymentmethod

	ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pbtransaction.ApiResponseTransaction
	ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pbtransaction.ApiResponseTransactionDeleteAt
	ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pbtransaction.ApiResponsesTransaction
	ToProtoResponseTransactionDelete(status string, message string) *pbtransaction.ApiResponseTransactionDelete
	ToProtoResponseTransactionAll(status string, message string) *pbtransaction.ApiResponseTransactionAll
	ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pbtransaction.ApiResponsePaginationTransactionDeleteAt
	ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pbtransaction.ApiResponsePaginationTransaction
}
