package merchant_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type MerchantApiTestSuite struct {
	tests.BaseTestSuite
	handler    http.Handler
	merchantID int
	userID     int
}

func (s *MerchantApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()
	s.SetupTransactionService()
	s.SetupCashierService()
	s.userID = s.SeedUser(context.Background())

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *MerchantApiTestSuite) TestMerchantApiLifecycle() {
	uid := strconv.Itoa(s.userID)
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		createMerchant(input: {
			user_id: `+uid+`, name: "Test Merchant", description: "Test Description",
			address: "Test Address", contact_email: "merchant@example.com",
			contact_phone: "123456789", status: "active"
		}) { status message data { id name } }
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createMerchant"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.Equal("Test Merchant", data["name"])
	s.merchantID = int(data["id"].(float64))

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByIdMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllMerchant(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveMerchant(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		updateMerchant(input: {
			merchant_id: `+strconv.Itoa(s.merchantID)+`, user_id: `+uid+`, name: "Updated Merchant",
			description: "Updated", address: "Updated", contact_email: "upd@example.com", contact_phone: "987", status: "active"
		}) { status message data { id name } }
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { trashedMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["trashedMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedMerchant(input: { page: 1, page_size: 10 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreMerchant(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteMerchantPermanent(input: { id: `+strconv.Itoa(s.merchantID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteMerchantPermanent"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllMerchant { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllMerchant"].(map[string]interface{})["status"])

	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllMerchantPermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllMerchantPermanent"].(map[string]interface{})["status"])
}

func TestMerchantApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantApiTestSuite))
}
