package tango

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
    "testing"
)

type HttpTestResponse struct {
    *HttpResponse
}

func (h *HttpTestResponse) Json() interface{} {
    v := new(interface{})
    err := json.Unmarshal([]byte(h.Content), v)
    if err != nil {
        return nil
    }

    return v
}

type testClient struct {
    t   *testing.T

    followRedirects bool
}

func NewTestClient(t *testing.T) *testClient {
    tc := &testClient{}
    tc.followRedirects = true

    return tc
}

// Follow redirect setting...

func (t *testClient) Head(path string, body ...string) *HttpTestResponse {
    return t.runMethod("HEAD", path, body)
}
func (t *testClient) Get(path string, body ...string) *HttpTestResponse {
    return t.runMethod("GET", path, body)
}
func (t *testClient) Post(path string, body ...string) *HttpTestResponse {
    return t.runMethod("POST", path, body)
}
func (t *testClient) Put(path string, body ...string) *HttpTestResponse {
    return t.runMethod("PUT", path, body)
}
func (t *testClient) Patch(path string, body ...string) *HttpTestResponse {
    return t.runMethod("PATCH", path, body)
}
func (t *testClient) Delete(path string, body ...string) *HttpTestResponse {
    return t.runMethod("DELETE", path, body)
}
func (t *testClient) Options(path string, body ...string) *HttpTestResponse {
    return t.runMethod("OPTIONS", path, body)
}

func (t *testClient) runMethod(method, path string, body []string) *HttpTestResponse {
    var bodyReader io.Reader = nil

    if len(body) != 0 {
        bodyReader = strings.NewReader(strings.TrimSpace(strings.Join(body, " ")))
    }

    req, _ := http.NewRequest(method, path, bodyReader)
    resp := Mux.ServeTestResponse(req)

    return &HttpTestResponse{resp}
}
