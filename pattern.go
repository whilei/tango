package tango

import (
    "fmt"
    "github.com/cojac/context"
    "github.com/cojac/mux"
    "net/http"
    "strings"
    "time"
)

const tangoRunTimeContextKey string = "TANGO_RUN_TIME"

var Mux = mux.NewRouter()

func Pattern(path string, handler HandlerInterface) {

    Mux.HandleFunc(path, http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            // Any panic errors will be caught and passed over to our ErrorHandler.
            defer func() {
                if rec := recover(); rec != nil {
                    LogError.Printf("Panic Recovered: %s", rec)
                    writePatResponse(handler.ErrorHandler(fmt.Sprintf("%s", rec)), w, r)
                }
            }()

            context.Set(r, tangoRunTimeContextKey, time.Now())
            request := NewHttpRequest(r)
            RunMiddlewarePreprocess(request)
            handler.Prepare(request)

            response := HttpResponseServerError()
            switch strings.ToUpper(r.Method) {
            case "HEAD":
                response = handler.Head(request)
                if response.StatusCode == http.StatusMethodNotAllowed {
                    resp2 := handler.Get(request)
                    if resp2.StatusCode == http.StatusOK {
                        response = resp2
                        response.Content = ""
                    }
                }
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

            writePatResponse(response, w, r)

            LogInfo.Printf("%d %s %s (%s) %s",
                response.StatusCode,
                strings.ToUpper(r.Method),
                r.RequestURI,
                r.RemoteAddr,
                time.Since(context.Get(r, tangoRunTimeContextKey).(time.Time)))
        }))
}

func writePatResponse(response *HttpResponse, w http.ResponseWriter, r *http.Request) {
    for k, v := range response.headers {
        w.Header().Add(k, v)
    }

    w.Header().Add("Content-Type", response.ContentType)
    w.WriteHeader(response.StatusCode)
    w.Write([]byte(response.Content))
}
