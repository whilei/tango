package tango

import (
    "github.com/bmizerany/assert"
    "github.com/cojac/mux"
    "net/http"
    "net/http/httptest"
    "testing"
)

//---
type MiddleHandler struct{ BaseHandler }

func (h MiddleHandler) Get(request *HttpRequest) *HttpResponse {
    return NewHttpResponse("foo")
}

//---
type Firstware struct{ BaseMiddleware }

func (m Firstware) ProcessRequest(request *HttpRequest) {
    request.Header().Set("X-pre", "superman")
}

func (m Firstware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    response.AddHeader("X-post", request.Header().Get("X-pre"))
    response.Content = "bar"
}

//---
type Secondware struct{ BaseMiddleware }

func (m Secondware) ProcessRequest(request *HttpRequest) {
    request.Header().Set("X-pre", "batman")
}

//---
type Thirdware struct{ BaseMiddleware }

func (m Thirdware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    response.Content = "foobar"
}

//---
func TestBasicMiddleware(t *testing.T) {
    Pattern("/", MiddleHandler{})
    Middleware(Firstware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "superman", rec.Header().Get("X-post"))

    Mux = mux.NewRouter()
    Middlewares = []MiddlewareInterface{}
}

func TestOrderOfMiddleware(t *testing.T) {
    Pattern("/", MiddleHandler{})
    Middleware(Firstware{})
    Middleware(Secondware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "batman", rec.Header().Get("X-post"))

    Mux = mux.NewRouter()
    Middlewares = []MiddlewareInterface{}
}

func TestOrderOfMiddlewareReversed(t *testing.T) {
    Pattern("/", MiddleHandler{})
    Middleware(Secondware{})
    Middleware(Firstware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "bar", rec.Body.String())
    assert.Equal(t, "superman", rec.Header().Get("X-post"))

    Mux = mux.NewRouter()
    Middlewares = []MiddlewareInterface{}
}

func TestMultiMiddleware(t *testing.T) {
    Pattern("/", MiddleHandler{})
    Middleware(Thirdware{})
    Middleware(Firstware{})
    Middleware(Secondware{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "foobar", rec.Body.String())
    assert.Equal(t, "batman", rec.Header().Get("X-post"))

    Mux = mux.NewRouter()
    Middlewares = []MiddlewareInterface{}
}
