package tango

import (
    "fmt"
)

type HttpResponse struct {
    Content     string
    StatusCode  int
    ContentType string
    headers     map[string]string
}

func (r *HttpResponse) AddHeader(k, v string) {
    dst := make(map[string]string)
    for k, v := range r.headers {
        dst[k] = v
    }
    dst[k] = v

    r.headers = dst
}

func (r *HttpResponse) DeleteHeader(k string) {
    delete(r.headers, k)
}

func (r *HttpResponse) GetHeader(k string) (string, bool) {
    val, ok := r.headers[k]
    return val, ok
}

func (r *HttpResponse) HasHeader(k string) bool {
    _, exists := r.headers[k]
    return exists
}

func NewHttpResponse(args ...interface{}) *HttpResponse {
    r := new(HttpResponse)

    content := ""
    status := 200
    contentType := "text/html; charset=utf-8"

    switch len(args) {
    case 3:
        contentType = args[2].(string)
        fallthrough
    case 2:
        status = args[1].(int)
        fallthrough
    case 1:
        content = args[0].(string)
    case 0:
        break
    default:
        panic(fmt.Sprintf("NewHttpResponse received [%d] args, can only handle 3.", len(args)))
    }

    r.Content = content
    r.StatusCode = status
    r.ContentType = contentType

    return r
}
