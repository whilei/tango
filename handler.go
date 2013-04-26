package tango

import (
    "net/http"
)

type HandlerInterface interface {
    Head(r *HttpRequest) *HttpResponse
    Get(r *HttpRequest) *HttpResponse
    Post(r *HttpRequest) *HttpResponse
    Put(r *HttpRequest) *HttpResponse
    Patch(r *HttpRequest) *HttpResponse
    Delete(r *HttpRequest) *HttpResponse
    Options(r *HttpRequest) *HttpResponse

    Prepare(r *HttpRequest)
    Finish(r *HttpRequest, response *HttpResponse)
    ErrorHandler(errorStr string) *HttpResponse
}

type BaseHandler struct{}

func (h BaseHandler) ErrorHandler(errorStr string) *HttpResponse {
    return HttpResponseServerError()
}

func (h BaseHandler) Prepare(r *HttpRequest) {
    // pass
}

func (h BaseHandler) Finish(r *HttpRequest, response *HttpResponse) {
    // pass
}

func (h BaseHandler) Head(request *HttpRequest) *HttpResponse {
    // By default, HEAD requests will mimic a GET request sans content.
    resp := h.Get(request)
    resp.Content = ""
    return resp
}

func (h BaseHandler) Get(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

func (h BaseHandler) Post(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

func (h BaseHandler) Put(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

func (h BaseHandler) Patch(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

func (h BaseHandler) Delete(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

func (h BaseHandler) Options(request *HttpRequest) *HttpResponse {
    return HttpResponseNotAllowed()
}

// The following are convience methods.
// TODO: Fix the redirect handler.
// func HttpResponsePermanentRedirect(url string) *HttpResponse {
//     // do redirect
//     //Redirect(w ResponseWriter, r *Request, urlStr string, code int) {
//     //http://golang.org/src/pkg/net/http/server.go?s=21797:21865#L740
//     return NewHttpResponse("Moved Permanently", http.StatusMovedPermanently, "text/plain")
// }

func HttpResponseRedirect(url string) *HttpResponse {
    // do redirect
    return NewHttpResponse("Found", http.StatusFound, "text/plain")
}

func HttpResponseNotModified() *HttpResponse {
    return NewHttpResponse("Not Modified", http.StatusNotModified, "text/plain")
}

func HttpResponseBadRequest() *HttpResponse {
    return NewHttpResponse("Bad Request", http.StatusBadRequest, "text/plain")
}

func HttpResponseForbidden() *HttpResponse {
    return NewHttpResponse("Forbidden", http.StatusForbidden, "text/plain")
}

func HttpResponseNotFound() *HttpResponse {
    return NewHttpResponse("Not Found", http.StatusNotFound, "text/plain")
}

func HttpResponseNotAllowed() *HttpResponse {
    response := NewHttpResponse("Method Not Allowed", http.StatusMethodNotAllowed, "text/plain")
    // TODO: How are we going to determine which methods are implemented on a given handler?
    //response.AddHeader("Allow", "")
    return response
}

func HttpResponseGone() *HttpResponse {
    return NewHttpResponse("Gone", http.StatusGone, "text/plain")
}

func HttpResponseServerError() *HttpResponse {
    return NewHttpResponse("Internal Server Error", http.StatusInternalServerError, "text/plain")
}
