package role_test

import (
	"net/http"
	"strconv"
	"testing"

	graphtest "github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest_wrapper"
	tests "github.com/MamangRust/microservice-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type RoleApiTestSuite struct {
	tests.BaseTestSuite
	handler http.Handler
	roleID  int
}

func (s *RoleApiTestSuite) SetupSuite() {
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

	resolver := graphtest.NewResolver(graphtest.ConnMap(s.Conns), s.Log, s.RedisClient())
	s.handler = graphtest.NewHandler(resolver)
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	resp, err := graphtest.ExecuteGraphQL(s.handler, `mutation {
		createRole(input: { name: "API Role" }) {
			status message data { id name }
		}
	}`, nil, "")
	s.Require().NoError(err)
	createResult := resp.Data["createRole"].(map[string]interface{})
	s.Equal("success", createResult["status"])
	data := createResult["data"].(map[string]interface{})
	s.Equal("API Role", data["name"])
	s.roleID = int(data["id"].(float64))

	// 2. FindAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findAllRole(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findAllRole"].(map[string]interface{})["status"])

	// 3. FindById
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByIdRole(input: { role_id: `+strconv.Itoa(s.roleID)+` }) { status message data { id name } } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByIdRole"].(map[string]interface{})["status"])

	// 4. FindByActive
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByActiveRole(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByActiveRole"].(map[string]interface{})["status"])

	// 5. FindByTrashed
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByTrashedRole(input: {}) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByTrashedRole"].(map[string]interface{})["status"])

	// 6. FindByUserId
	resp, err = graphtest.ExecuteGraphQL(s.handler, `query { findByUserIdRole(input: { user_id: 1 }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["findByUserIdRole"].(map[string]interface{})["status"])

	// 7. Update
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation {
		updateRole(input: { id: `+strconv.Itoa(s.roleID)+`, name: "Updated API Role" }) {
			status message data { id name }
		}
	}`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["updateRole"].(map[string]interface{})["status"])

	// 8. Restore
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreRole(input: { role_id: `+strconv.Itoa(s.roleID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreRole"].(map[string]interface{})["status"])

	// 9. DeletePermanent
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteRolePermanent(input: { role_id: `+strconv.Itoa(s.roleID)+` }) { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteRolePermanent"].(map[string]interface{})["status"])

	// 10. RestoreAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { restoreAllRole { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["restoreAllRole"].(map[string]interface{})["status"])

	// 11. DeleteAll
	resp, err = graphtest.ExecuteGraphQL(s.handler, `mutation { deleteAllRolePermanent { status message } }`, nil, "")
	s.Require().NoError(err)
	s.Equal("success", resp.Data["deleteAllRolePermanent"].(map[string]interface{})["status"])
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
