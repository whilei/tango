package tango

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "reflect"
    "strings"
    "testing"
)

type HttpTestResponse struct {
    *HttpResponse
}

func (h *HttpTestResponse) JsonMap() map[string]interface{} {
    v := make(map[string]interface{})
    err := json.Unmarshal([]byte(h.Content), &v)
    if err != nil {
        return nil
    }

    return v
}

func (h *HttpTestResponse) JsonArray() []interface{} {
    var v interface{}
    err := json.Unmarshal([]byte(h.Content), &v)
    if err != nil {
        return nil
    }

    return v.([]interface{})
}

type testClient struct {
    argTesting *testing.T

    followRedirects bool
}

func NewTestClient(t *testing.T) *testClient {
    tc := &testClient{}
    tc.followRedirects = true

    return tc
}

// Follow redirect setting...

func (t *testClient) Head(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("HEAD", path, input)
}
func (t *testClient) Get(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("GET", path, input)
}
func (t *testClient) Post(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("POST", path, input)
}
func (t *testClient) Put(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("PUT", path, input)
}
func (t *testClient) Patch(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("PATCH", path, input)
}
func (t *testClient) Delete(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("DELETE", path, input)
}
func (t *testClient) Options(path string, input ...interface{}) *HttpTestResponse {
    return t.runMethod("OPTIONS", path, input)
}

// Helper for above convenience methods.
func (t *testClient) runMethod(method, path string, input []interface{}) *HttpTestResponse {
    var body []string
    data := make(map[string]interface{})

    for _, arg := range input {
        k := reflect.ValueOf(arg).Kind()
        if k == reflect.Map {
            for kArg, vArg := range arg.(map[string]string) {
                data[kArg] = vArg
            }
        } else {
            body = append(body, arg.(string))
        }
    }

    var bodyReader io.Reader = nil

    if len(body) != 0 {
        bodyReader = strings.NewReader(strings.TrimSpace(strings.Join(body, " ")))
    }

    req, _ := http.NewRequest(method, path, bodyReader)
    if len(data) != 0 {
        uV := make(url.Values)
        for k, v := range data {
            uV.Add(k, fmt.Sprintf("%v", v))
        }
        req.PostForm = uV
    } else if len(data) == 0 && len(body) >= 1 {
        req.PostForm, _ = url.ParseQuery(body[0])
    }
    return t.Request(req)
}

func (t *testClient) Request(req *http.Request) *HttpTestResponse {
    resp := Mux.ServeTestResponse(req)
    return &HttpTestResponse{resp}
}
