package tango

import (
    "fmt"
    "net/http"
)

type HttpResponse struct {
    Content     string
    StatusCode  int
    ContentType string
    Context     map[string]interface{}
    Header      http.Header
    isFinished  bool
}

func NewHttpResponse(args ...interface{}) *HttpResponse {
    r := new(HttpResponse)

    content := ""
    status := 200
    contentType := "text/html; charset=" + Settings.String("charset", "utf-8")

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
    r.isFinished = false

    r.Header = make(http.Header)
    r.Context = make(map[string]interface{})

    return r
}

func (h HttpResponse) Finish() {
    h.isFinished = true
}
