package tango

import (
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
)

type HttpResponse struct {
    Content     string
    StatusCode  int
    ContentType string
    Context     map[string]interface{}
    Header      http.Header
}

func NewHttpResponse(args ...interface{}) *HttpResponse {
    r := &HttpResponse{}

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
        arg := args[0]
        val := reflect.ValueOf(arg)
        k := val.Kind()

        if k == reflect.Map || k == reflect.Slice || k == reflect.Array {
            // If out object is a map || slice || array, assume it's for Json output.
            contentType = "application/json; charset=" + Settings.String("charset", "utf-8")

            out, err := json.Marshal(arg)
            if err != nil {
                panic(fmt.Sprintf("NewJsonResponse could not convert map/array to json: %s", err))
            }
            content = string(out)
        } else {
            content = args[0].(string)
        }
    case 0:
        break
    default:
        panic(fmt.Sprintf("NewHttpResponse received [%d] args, can only handle 3.", len(args)))
    }

    r.Content = content
    r.StatusCode = status
    r.ContentType = contentType

    r.Header = make(http.Header)
    r.Context = make(map[string]interface{})

    return r
}
