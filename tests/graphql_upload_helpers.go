package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/MamangRust/microservice-point-of-sale-apigateway/graphtest"
)

// GraphQLUploadOp executes a single-root-field mutation whose input carries a
// file upload, using the GraphQL multipart request spec. filePath is the
// variable path the uploaded part replaces (e.g. "variables.input.image"); the
// caller leaves that path set to nil in variables. It fails the test on
// transport or GraphQL errors and returns that root field's response object.
func (s *BaseTestSuite) GraphQLUploadOp(handler http.Handler, op, query string, variables map[string]interface{}, filePath, filename string) map[string]interface{} {
	resp := s.GraphQLUploadRaw(handler, query, variables, filePath, filename)
	s.Require().Empty(resp.Errors, "graphql errors: %v", resp.Errors)

	raw, ok := resp.Data[op]
	s.Require().True(ok, "missing operation %q in response %v", op, resp.Data)
	obj, ok := raw.(map[string]interface{})
	s.Require().True(ok, "operation %q is %T, not an object", op, raw)
	return obj
}

// GraphQLUploadRaw is GraphQLUploadOp without the error assertions.
func (s *BaseTestSuite) GraphQLUploadRaw(handler http.Handler, query string, variables map[string]interface{}, filePath, filename string) *graphtest.GraphQLResponse {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	ops, err := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	s.Require().NoError(err)
	s.Require().NoError(w.WriteField("operations", string(ops)))

	mapping, err := json.Marshal(map[string][]string{"0": {filePath}})
	s.Require().NoError(err)
	s.Require().NoError(w.WriteField("map", string(mapping)))

	fw, err := w.CreateFormFile("0", filename)
	s.Require().NoError(err)
	_, err = fw.Write([]byte("dummy image content"))
	s.Require().NoError(err)
	s.Require().NoError(w.Close())

	req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Result().Body)
	s.Require().NoError(err)

	var resp graphtest.GraphQLResponse
	s.Require().NoError(json.Unmarshal(body, &resp))
	return &resp
}
