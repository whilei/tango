package tango

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unrolled/assert"
)

func init() {
	// Disable logging when running the tests.
	LogInfo = log.New(ioutil.Discard, "", log.LstdFlags)
	LogError = log.New(ioutil.Discard, "", log.LstdFlags)
}

//---
type AlmostBaseHandler struct{ BaseHandler }

func (h *AlmostBaseHandler) New() HandlerInterface {
	return &AlmostBaseHandler{}
}

func TestHandlerMethods(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/", &AlmostBaseHandler{})

	methods := []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	for _, meth := range methods {
		r, _ := http.NewRequest(meth, "/", nil)
		rec := httptest.NewRecorder()
		Mux.ServeHTTP(rec, r)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	}
}

//---
type AppendSlashBaseHandler struct{ BaseHandler }

func (h *AppendSlashBaseHandler) New() HandlerInterface {
	return &AppendSlashBaseHandler{}
}
func (h *AppendSlashBaseHandler) Get(request *HttpRequest) *HttpResponse {
	return NewHttpResponse("hello")
}

func TestHandlerAppendSlash(t *testing.T) {
	oldAppendSetting := Settings.Bool("append_slash", false)
	oldAppendRedirectSetting := Settings.Bool("append_slash_should_redirect", true)

	Settings.Set("append_slash", false)
	Pattern("/hello", &AppendSlashBaseHandler{})
	r, _ := http.NewRequest("GET", "/hello/", nil)
	rec := httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	Mux = &PatternServeMux{}

	Settings.Set("append_slash", true)
	Pattern("/hello", &AppendSlashBaseHandler{})
	r, _ = http.NewRequest("GET", "/hello/", nil)
	rec = httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusOK, rec.Code)
	Mux = &PatternServeMux{}

	Settings.Set("append_slash", true)
	Settings.Set("append_slash_should_redirect", true)
	Pattern("/hello", &AppendSlashBaseHandler{})
	r, _ = http.NewRequest("GET", "/hello/", nil)
	rec = httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	Mux = &PatternServeMux{}

	// Set it back to what it was.
	Settings.Set("append_slash", oldAppendSetting)
	Settings.Set("append_slash_should_redirect", oldAppendRedirectSetting)
}

//---
type GetHandler struct{ BaseHandler }

func (h *GetHandler) New() HandlerInterface {
	return &GetHandler{}
}
func (h *GetHandler) Get(request *HttpRequest) *HttpResponse {
	return NewHttpResponse("hello")
}

func TestHandlerHead(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/", &GetHandler{})

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
}

//---
var OneOffPreFinTestObj string

type PrepFinHandler struct{ BaseHandler }

func (h *PrepFinHandler) New() HandlerInterface {
	OneOffPreFinTestObj = "NEW"
	return &PrepFinHandler{}
}
func (h *PrepFinHandler) Get(request *HttpRequest) *HttpResponse {
	return NewHttpResponse(OneOffPreFinTestObj)
}
func (h *PrepFinHandler) Prepare(r *HttpRequest) *HttpResponse {
	OneOffPreFinTestObj = OneOffPreFinTestObj + "-PRE"

	return nil
}

func (h *PrepFinHandler) Finish(r *HttpRequest, response *HttpResponse) {
	OneOffPreFinTestObj = OneOffPreFinTestObj + "-FIN"
}

func TestHandlerPrepareFinish(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/", &PrepFinHandler{})

	r, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "NEW-PRE", rec.Body.String())
	assert.Equal(t, "NEW-PRE-FIN", OneOffPreFinTestObj)
}

//---

//---
type GenericHandler struct{ BaseHandler }

func (h *GenericHandler) New() HandlerInterface {
	return &GenericHandler{}
}
func (h *GenericHandler) Get(request *HttpRequest) *HttpResponse {
	return h.HttpResponseGone()
}
func (h *GenericHandler) Post(request *HttpRequest) *HttpResponse {
	return h.HttpResponseNotFound()
}
func (h *GenericHandler) Put(request *HttpRequest) *HttpResponse {
	return h.HttpResponseForbidden()
}
func (h *GenericHandler) Delete(request *HttpRequest) *HttpResponse {
	return h.HttpResponseBadRequest()
}
func (h *GenericHandler) Patch(request *HttpRequest) *HttpResponse {
	return h.HttpResponseNotModified()
}
func (h *GenericHandler) Option(request *HttpRequest) *HttpResponse {
	return h.HttpResponseServerError()
}

func TestHandlerGenericResponses(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/", &GenericHandler{})

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
}

//---
type RedirectHandler struct{ BaseHandler }

func (h *RedirectHandler) New() HandlerInterface {
	return &RedirectHandler{}
}
func (h *RedirectHandler) Get(request *HttpRequest) *HttpResponse {
	return h.PermanentRedirect(request, "/next/")
}
func (h *RedirectHandler) Post(request *HttpRequest) *HttpResponse {
	return h.TemporaryRedirect(request, "next/")
}
func (h *RedirectHandler) Put(request *HttpRequest) *HttpResponse {
	return h.TemporaryRedirect(request, "next")
}
func (h *RedirectHandler) Delete(request *HttpRequest) *HttpResponse {
	return h.PermanentRedirect(request, "../next")
}
func (h *RedirectHandler) Options(request *HttpRequest) *HttpResponse {
	return h.PermanentRedirect(request, "next/?foo=bar")
}
func (h *RedirectHandler) Patch(request *HttpRequest) *HttpResponse {
	request.URL.Path = ""
	return h.TemporaryRedirect(request, "/next/")
}

func TestHandlerRedirect(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/start/", &RedirectHandler{})

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
}

//---
type ErrHandler struct{ BaseHandler }
type CustomErrHandler struct{ BaseHandler }

func (h *ErrHandler) New() HandlerInterface {
	return &ErrHandler{}
}
func (h *ErrHandler) Get(request *HttpRequest) *HttpResponse {
	panic("foo")
}

func (h *CustomErrHandler) New() HandlerInterface {
	return &CustomErrHandler{}
}
func (h *CustomErrHandler) Get(request *HttpRequest) *HttpResponse {
	panic("bar")
}
func (h *CustomErrHandler) ErrorHandler(errorStr string) *HttpResponse {
	resp := NewHttpResponse(`{"hello": "world"}`, 400, "application/json")
	return resp
}

func TestHandlerErrors(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/err/", &ErrHandler{})
	Pattern("/custom/", &CustomErrHandler{})

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
}

//---
type AllowedHandler struct{ BaseHandler }

func (h *AllowedHandler) New() HandlerInterface {
	return &AllowedHandler{}
}
func (h *AllowedHandler) Get(request *HttpRequest) *HttpResponse {
	return NewHttpResponse("This is allowed")
}
func (h *AllowedHandler) Post(request *HttpRequest) *HttpResponse {
	return NewHttpResponse("This is also allowed")
}

func TestHandlerAllowedMethods(t *testing.T) {
	defer func() { Mux = &PatternServeMux{} }()
	Pattern("/allowed/", &AllowedHandler{})

	r, _ := http.NewRequest("GET", "/allowed/", nil)
	rec := httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusOK, rec.Code)

	r, _ = http.NewRequest("POST", "/allowed/", nil)
	rec = httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusOK, rec.Code)

	r, _ = http.NewRequest("PUT", "/allowed/", nil)
	rec = httptest.NewRecorder()
	Mux.ServeHTTP(rec, r)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	// TODO: Still need to fix this sometime down the road.
	//assert.Equal(t, "GET,POST", rec.Header().Get("Allow"))
}
