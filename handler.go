package tango

import (
    "html"
    "net/http"
    "net/url"
    "path"
    "strings"
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
    return h.HttpResponseServerError()
}

func (h BaseHandler) Prepare(r *HttpRequest) {
    // pass
}

func (h BaseHandler) Finish(r *HttpRequest, response *HttpResponse) {
    // pass
}

func (h BaseHandler) Head(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Get(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Post(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Put(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Patch(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Delete(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) Options(request *HttpRequest) *HttpResponse {
    return h.HttpResponseNotAllowed()
}

func (h BaseHandler) PermanentRedirect(request *HttpRequest, urlStr string) *HttpResponse {
    return h.redirect(request, urlStr, http.StatusMovedPermanently)
}

func (h BaseHandler) TemporaryRedirect(request *HttpRequest, urlStr string) *HttpResponse {
    return h.redirect(request, urlStr, http.StatusTemporaryRedirect)
}

func (h BaseHandler) redirect(r *HttpRequest, urlStr string, code int) *HttpResponse {
    if u, err := url.Parse(urlStr); err == nil {
        oldpath := r.URL.Path
        if oldpath == "" {
            oldpath = "/"
        }

        if u.Scheme == "" {
            // no leading http://server
            if urlStr == "" || urlStr[0] != '/' {
                // make relative path absolute
                olddir, _ := path.Split(oldpath)
                urlStr = olddir + urlStr
            }

            var query string
            if i := strings.Index(urlStr, "?"); i != -1 {
                urlStr, query = urlStr[:i], urlStr[i:]
            }

            // clean up but preserve trailing slash
            trailing := urlStr[len(urlStr)-1] == '/'
            urlStr = path.Clean(urlStr)
            if trailing && urlStr[len(urlStr)-1] != '/' {
                urlStr += "/"
            }
            urlStr += query
        }
    }

    response := NewHttpResponse()
    response.Header.Set("Location", urlStr)
    response.StatusCode = code

    // RFC2616 recommends that a short note "SHOULD" be included in the
    // response because older user agents may not understand 301/307.
    // Shouldn't send the response for POST or HEAD; that leaves GET.
    if r.Method == "GET" {
        response.Content = "<a href=\"" + html.EscapeString(urlStr) + "\">" + http.StatusText(code) + "</a>.\n"
    }

    return response
}

func (h BaseHandler) HttpResponseNotModified() *HttpResponse {
    return shortHttpReturn(http.StatusNotModified)
}

func (h BaseHandler) HttpResponseBadRequest() *HttpResponse {
    return shortHttpReturn(http.StatusBadRequest)
}

func (h BaseHandler) HttpResponseForbidden() *HttpResponse {
    return shortHttpReturn(http.StatusForbidden)
}

func (h BaseHandler) HttpResponseNotFound() *HttpResponse {
    return shortHttpReturn(http.StatusNotFound)
}

func (h BaseHandler) HttpResponseNotAllowed() *HttpResponse {
    return shortHttpReturn(http.StatusMethodNotAllowed)
}

func (h BaseHandler) HttpResponseGone() *HttpResponse {
    return shortHttpReturn(http.StatusGone)
}

func (h BaseHandler) HttpResponseServerError() *HttpResponse {
    return shortHttpReturn(http.StatusInternalServerError)
}

func shortHttpReturn(code int) *HttpResponse {
    return NewHttpResponse(http.StatusText(code), code, "text/plain")
}
