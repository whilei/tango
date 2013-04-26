package tango

import (
    "net/http"
    "net/url"
)

type HttpRequest struct {
    Raw      *http.Request
    Host     string
    Path     string
    Body     string
    Scheme   string
    RawQuery string
    Fragment string

    PathParams map[string]string
    GetParams  url.Values
    PostParams map[string]string
    MetaParams map[string]string
}

func NewHttpRequest() *HttpRequest {
    r := new(HttpRequest)
    return r
}

func (r HttpRequest) AbsoluteURI() string {
    return "n/a"
}

func (r HttpRequest) isSecure() bool {
    return false
}

func (r HttpRequest) isAjax() bool {
    return false
}
