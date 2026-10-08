package tests

import (
	"net/http"

	"github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest"
	"github.com/MamangRust/microservice-point-of-sale-apigateway/testhelper"
)

// NewGraphQLHandler builds a GraphQL HTTP handler wired to every gRPC service
// connection set up by the suite (s.Conns). It is the GraphQL-gateway analogue
// of the old REST handler registration.
func (s *BaseTestSuite) NewGraphQLHandler() http.Handler {
	conns := testhelper.SetupServiceConnections(s.Conns)
	return testhelper.NewTestHandler(conns, s.RedisClient(), s.Log)
}

// GraphQLOp executes a single-root-field query/mutation and returns that root
// field's response object (e.g. the value of "createOrder"). It fails the test
// on transport errors or GraphQL errors.
func (s *BaseTestSuite) GraphQLOp(handler http.Handler, op, query string, variables map[string]interface{}) map[string]interface{} {
	resp, err := graphtest.ExecuteGraphQL(handler, query, variables, "")
	s.Require().NoError(err)
	s.Require().Empty(resp.Errors, "graphql errors: %v", resp.Errors)

	raw, ok := resp.Data[op]
	s.Require().True(ok, "missing operation %q in response %v", op, resp.Data)
	obj, ok := raw.(map[string]interface{})
	s.Require().True(ok, "operation %q is %T, not an object", op, raw)
	return obj
}

// GraphQLRaw executes a query/mutation and returns the whole response, without
// asserting on errors. It is used by tests that expect a GraphQL error.
func (s *BaseTestSuite) GraphQLRaw(handler http.Handler, query string, variables map[string]interface{}) *graphtest.GraphQLResponse {
	resp, err := graphtest.ExecuteGraphQL(handler, query, variables, "")
	s.Require().NoError(err)
	return resp
}

// GQLData returns the "data" object of an operation response (which may be nil
// for delete-only responses).
func GQLData(op map[string]interface{}) map[string]interface{} {
	data, _ := op["data"].(map[string]interface{})
	return data
}

// GQLID extracts a numeric id from a GraphQL data object.
func GQLID(data map[string]interface{}, key string) int {
	v, ok := data[key]
	if !ok {
		return 0
	}
	f, _ := v.(float64)
	return int(f)
}

// GQLList extracts a list from a GraphQL data object.
func GQLList(data map[string]interface{}, key string) []interface{} {
	list, _ := data[key].([]interface{})
	return list
}

// GQLFieldList extracts a list-valued field directly from an operation
// response, e.g. the "data" list of an ApiResponses* payload.
func GQLFieldList(op map[string]interface{}, key string) []interface{} {
	list, _ := op[key].([]interface{})
	return list
}
