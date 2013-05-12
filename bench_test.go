// Go 1.1: go test -bench=.
package tango

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

//---
type IndexHandler struct {
    BaseHandler
    BenchMixin
}

func (h *IndexHandler) New() HandlerInterface {
    return &IndexHandler{}
}
func (h *IndexHandler) Get(request *HttpRequest) *HttpResponse {
    passedId, _ := request.PathValue(":id")
    return NewHttpResponse(fmt.Sprintf("Hello, world: %s | %s", passedId, h.TesterStr))
}

//---
type Benchware struct{ BaseMiddleware }

func (m *Benchware) ProcessRequest(request *HttpRequest, response *HttpResponse) {
    request.Header.Set("X-pre", "superman")
}

func (m *Benchware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    response.Header.Set("X-post", request.Header.Get("X-pre"))
}

//---
type BenchMixin struct {
    BaseMixin
    TesterStr string
}

func (m *BenchMixin) PrepareBenchMixin() {
    m.TesterStr = "superman"
}

func (m *BenchMixin) FinishBenchMixin() {
    m.TesterStr = "batman"
}

//---
func BenchmarkTango(b *testing.B) {
    Pattern("/hello/:id", &IndexHandler{})
    Middleware(&Benchware{})
    Mixin(&BaseMixin{})

    for i := 0; i < b.N; i++ {
        url := fmt.Sprintf("/hello/%d", i)
        rec := httptest.NewRecorder()
        request, _ := http.NewRequest("GET", url, nil)
        Mux.ServeHTTP(rec, request)
    }
}
