package tango

import (
    "fmt"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

const tangoRunTimeContextKey string = "TANGO_RUN_TIME"

var Mux = mux.NewRouter()

func Pattern(path string, handler HandlerInterface) {

    Mux.HandleFunc(path, http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    LogError.Printf("Panic Recovered: %s", rec)
                    writeResponse(handler.ErrorHandler(fmt.Sprintf("%s", rec)), w, r)
                }
            }()

            context.Set(r, tangoRunTimeContextKey, time.Now())

            request := NewHttpRequest()
            request.Raw = r
            request.Host = r.URL.Host
            request.Path = r.URL.Path
            request.Scheme = r.URL.Scheme
            request.RawQuery = r.URL.RawQuery
            request.Fragment = r.URL.Fragment
            request.PathParams = mux.Vars(r)
            request.GetParams = r.URL.Query()

            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                panic("Error reading body")
            }
            request.Body = string(body)

            RunMiddlewarePreprocess(request)
            handler.Prepare(request)

            response := HttpResponseServerError()
            switch strings.ToUpper(r.Method) {
            case "HEAD":
                response = handler.Head(request)
            case "GET":
                response = handler.Get(request)
            case "POST":
                response = handler.Post(request)
            case "PUT":
                response = handler.Put(request)
            case "PATCH":
                response = handler.Patch(request)
            case "DELETE":
                response = handler.Delete(request)
            case "OPTIONS":
                response = handler.Options(request)
            }

            handler.Finish(request, response)
            RunMiddlewarePostprocess(request, response)

            writeResponse(response, w, r)

            LogInfo.Printf("%d %s %s (%s) %s",
                response.StatusCode,
                strings.ToUpper(r.Method),
                r.RequestURI,
                r.RemoteAddr,
                time.Since(context.Get(r, tangoRunTimeContextKey).(time.Time)))
        }))
}

func writeResponse(response *HttpResponse, w http.ResponseWriter, r *http.Request) {
    for k, v := range response.headers {
        w.Header().Add(k, v)
    }

    w.Header().Add("Content-Type", response.ContentType)
    w.WriteHeader(response.StatusCode)
    w.Write([]byte(response.Content))
}
