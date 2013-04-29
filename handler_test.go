package tango

import (
    "github.com/bmizerany/assert"
    "github.com/cojac/mux"
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
func TestHandlerMethods(t *testing.T) {
    Pattern("/", BaseHandler{})

    methods := []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    for _, meth := range methods {
        r, _ := http.NewRequest(meth, "/", nil)
        rec := httptest.NewRecorder()
        Mux.ServeHTTP(rec, r)
        assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
    }

    Mux = mux.NewRouter()
}

//---
type GetHandler struct{ BaseHandler }

func (h GetHandler) Get(request *HttpRequest) *HttpResponse {
    return NewHttpResponse("hello")
}

func TestHandlerHead(t *testing.T) {
    Pattern("/", GetHandler{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "hello", rec.Body.String())

    getResult := rec.Body

    r, _ = http.NewRequest("HEAD", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusOK, rec.Code)

    assert.NotEqual(t, getResult, rec.Body)
    assert.Equal(t, 0, len(rec.Body.String()))

    Mux = mux.NewRouter()
}

//---
type PrepFinHandler struct{ BaseHandler }

func (h PrepFinHandler) Get(request *HttpRequest) *HttpResponse {
    return NewHttpResponse("PrepFin")
}
func (h PrepFinHandler) Prepare(r *HttpRequest) {
    r.Header().Add("X-pre", "superman")
}

func (h PrepFinHandler) Finish(r *HttpRequest, response *HttpResponse) {
    response.AddHeader("X-pre", r.Header().Get("X-pre"))
    response.AddHeader("X-fin", "batman")
}

func TestHandlerPrepareFinish(t *testing.T) {
    Pattern("/", PrepFinHandler{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "superman", rec.Header().Get("X-pre"))
    assert.Equal(t, "batman", rec.Header().Get("X-fin"))

    Mux = mux.NewRouter()
}

//---
type GenericHandler struct{ BaseHandler }

func (h GenericHandler) Get(request *HttpRequest) *HttpResponse {
    return HttpResponseGone()
}
func (h GenericHandler) Post(request *HttpRequest) *HttpResponse {
    return HttpResponseNotFound()
}
func (h GenericHandler) Put(request *HttpRequest) *HttpResponse {
    return HttpResponseForbidden()
}
func (h GenericHandler) Delete(request *HttpRequest) *HttpResponse {
    return HttpResponseBadRequest()
}
func (h GenericHandler) Patch(request *HttpRequest) *HttpResponse {
    return HttpResponseNotModified()
}
func (h GenericHandler) Option(request *HttpRequest) *HttpResponse {
    return HttpResponseServerError()
}

func TestHandlerGenericResponses(t *testing.T) {
    Pattern("/", GenericHandler{})

    r, _ := http.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusGone, rec.Code)

    r, _ = http.NewRequest("POST", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusNotFound, rec.Code)

    r, _ = http.NewRequest("PUT", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusForbidden, rec.Code)

    r, _ = http.NewRequest("DELETE", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusBadRequest, rec.Code)

    r, _ = http.NewRequest("PATCH", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusNotModified, rec.Code)

    r, _ = http.NewRequest("OPTION", "/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusInternalServerError, rec.Code)

    Mux = mux.NewRouter()
}

//---
type RedirectHandler struct{ BaseHandler }

func (h RedirectHandler) Get(request *HttpRequest) *HttpResponse {
    return h.PermanentRedirect(request, "/next/")
}
func (h RedirectHandler) Post(request *HttpRequest) *HttpResponse {
    return h.TemporaryRedirect(request, "next/")
}
func (h RedirectHandler) Put(request *HttpRequest) *HttpResponse {
    return h.TemporaryRedirect(request, "next")
}
func (h RedirectHandler) Delete(request *HttpRequest) *HttpResponse {
    return h.PermanentRedirect(request, "../next")
}
func (h RedirectHandler) Options(request *HttpRequest) *HttpResponse {
    return h.PermanentRedirect(request, "next/?foo=bar")
}
func (h RedirectHandler) Patch(request *HttpRequest) *HttpResponse {
    request.RawRequest.URL.Path = ""
    return h.TemporaryRedirect(request, "/next/")
}

func TestHandlerRedirect(t *testing.T) {
    Pattern("/start/", RedirectHandler{})

    r, _ := http.NewRequest("GET", "/start/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusMovedPermanently, rec.Code)
    assert.Equal(t, "/next/", rec.Header().Get("Location"))

    r, _ = http.NewRequest("POST", "/start/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
    assert.Equal(t, "/start/next/", rec.Header().Get("Location"))

    r, _ = http.NewRequest("PUT", "/start/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
    assert.Equal(t, "/start/next", rec.Header().Get("Location"))

    r, _ = http.NewRequest("DELETE", "/start/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusMovedPermanently, rec.Code)
    assert.Equal(t, "/next", rec.Header().Get("Location"))

    r, _ = http.NewRequest("OPTIONS", "/start/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusMovedPermanently, rec.Code)
    assert.Equal(t, "/start/next/?foo=bar", rec.Header().Get("Location"))

    r, _ = http.NewRequest("PATCH", "/start/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
    assert.Equal(t, "/next/", rec.Header().Get("Location"))

    Mux = mux.NewRouter()
}

//---
type ErrHandler struct{ BaseHandler }
type CustomErrHandler struct{ BaseHandler }

func (h ErrHandler) Get(request *HttpRequest) *HttpResponse {
    panic("foo")
}

func (h CustomErrHandler) Get(request *HttpRequest) *HttpResponse {
    panic("bar")
}
func (h CustomErrHandler) ErrorHandler(errorStr string) *HttpResponse {
    resp := NewHttpResponse(`{"hello": "world"}`, 400, "application/json")
    return resp
}

func TestHandlerErrors(t *testing.T) {
    Pattern("/err/", ErrHandler{})
    Pattern("/custom/", CustomErrHandler{})

    r, _ := http.NewRequest("GET", "/err/", nil)
    rec := httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, http.StatusInternalServerError, rec.Code)
    assert.Equal(t, http.StatusText(500), rec.Body.String())

    r, _ = http.NewRequest("GET", "/custom/", nil)
    rec = httptest.NewRecorder()
    Mux.ServeHTTP(rec, r)
    assert.Equal(t, 400, rec.Code)
    assert.Equal(t, `{"hello": "world"}`, rec.Body.String())
    assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

    Mux = mux.NewRouter()
}
