package tango

import (
    "github.com/cojac/assert"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
)

func init() {
    // Disable logging when running the tests.
    LogInfo = log.New(ioutil.Discard, "", log.LstdFlags)
    LogError = log.New(ioutil.Discard, "", log.LstdFlags)
}

//---
type MiddleHandler struct{ BaseHandler }

func (h *MiddleHandler) New() HandlerInterface {
    return &MiddleHandler{}
}
func (h *MiddleHandler) Get(request *HttpRequest) *HttpResponse {
    return NewHttpResponse("foo")
}

//---
type Firstware struct{ BaseMiddleware }

func (m *Firstware) ProcessRequest(request *HttpRequest, response *HttpResponse) {
    request.Header.Set("X-pre", "superman")
}

func (m *Firstware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    response.Header.Set("X-post", request.Header.Get("X-pre"))
    response.Content = "bar"
}

//---
type Secondware struct{ BaseMiddleware }

func (m *Secondware) ProcessRequest(request *HttpRequest, response *HttpResponse) {
    request.Header.Set("X-pre", "batman")
}

//---
type Thirdware struct{ BaseMiddleware }

func (m *Thirdware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    response.Content = "foobar"
}

//---
type Finishware struct{ BaseMiddleware }

func (m *Finishware) ProcessRequest(request *HttpRequest, response *HttpResponse) {
    request.Header.Set("X-pre", "FIRST")
    response.Finish()
}

//---
func TestBasicMiddleware(t *testing.T) {
    defer func() { Mux = &PatternServeMux{} }()
    defer func() { middlewares = []MiddlewareInterface{} }()
    Pattern("/", &MiddleHandler{})
    Middleware(&Firstware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "superman", rec.Header().Get("X-post"))

}

func TestFinishMiddleware(t *testing.T) {
    defer func() { Mux = &PatternServeMux{} }()
    defer func() { middlewares = []MiddlewareInterface{} }()
    Pattern("/", &MiddleHandler{})

    // Finish should not allow Firstware to be called.
    Middleware(&Finishware{})
    Middleware(&Firstware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "FIRST", rec.Header().Get("X-post"))

}

func TestOrderOfMiddleware(t *testing.T) {
    defer func() { Mux = &PatternServeMux{} }()
    defer func() { middlewares = []MiddlewareInterface{} }()
    Pattern("/", &MiddleHandler{})
    Middleware(&Firstware{})
    Middleware(&Secondware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "batman", rec.Header().Get("X-post"))
}

func TestOrderOfMiddlewareReversed(t *testing.T) {
    defer func() { Mux = &PatternServeMux{} }()
    defer func() { middlewares = []MiddlewareInterface{} }()
    Pattern("/", &MiddleHandler{})
    Middleware(&Secondware{})
    Middleware(&Firstware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "superman", rec.Header().Get("X-post"))
}

func TestMultiMiddleware(t *testing.T) {
    defer func() { Mux = &PatternServeMux{} }()
    defer func() { middlewares = []MiddlewareInterface{} }()
    Pattern("/", &MiddleHandler{})
    Middleware(&Thirdware{})
    Middleware(&Firstware{})
    Middleware(&Secondware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "foobar", rec.Body.String())
    assert.Equal(t, "batman", rec.Header().Get("X-post"))
}
