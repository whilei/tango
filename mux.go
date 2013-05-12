package tango

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
)

type PatternServeMux struct {
    handlers []*patHandler
}

var Mux = &PatternServeMux{}

func init() {
    http.Handle("/", Mux)
}

// ServeHTTP matches r.URL.Path against its routing table using the rules
// described above.
func (p *PatternServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, ph := range p.handlers {
        if params, ok := ph.try(r.URL.Path); ok {
            if ph.isSlashRedirect {
                http.Redirect(w, r, buildUrlWithSlash(r), http.StatusMovedPermanently)
                return
            }

            ph.ServeHandlerHttp(w, r, params)
            return
        }
    }
}

func Pattern(pat string, h HandlerInterface) {
    Mux.handlers = append(Mux.handlers, &patHandler{pat, h, false})

    if Settings.Bool("append_slash", false) {
        n := len(pat)
        if n > 0 && pat[n-1] == '/' {
            Mux.handlers = append(Mux.handlers, &patHandler{pat[:n-1], h, true})
        }
    }
}

func buildUrlWithSlash(r *http.Request) string {
    result := r.URL.Path + "/"

    if len(r.URL.Query()) != 0 {
        result = result + "?" + r.URL.RawQuery
    }

    if len(r.URL.Fragment) != 0 {
        result = result + "#" + r.URL.Fragment
    }

    return result
}

type patHandler struct {
    pat string
    HandlerInterface
    isSlashRedirect bool
}

func (ph *patHandler) ServeHandlerHttp(w http.ResponseWriter, r *http.Request, params url.Values) {
    handler := ph.New()

    // Any panic errors will be caught and passed over to our ErrorHandler.
    defer func() {
        if rec := recover(); rec != nil {
            LogError.Printf("Panic Recovered: %s", rec)
            writePatternResponse(handler.ErrorHandler(fmt.Sprintf("%q", rec)), w)
        }
    }()

    start_request := time.Now()
    runMixinPrepare(handler)

    request := NewHttpRequest(r, params)
    response := NewHttpResponse()
    runMiddlewarePreprocess(request, response)

    // Only if the response has not finished should we let the handler touch it.
    if !response.isFinished {
        handler.Prepare(request, response)
        defer handler.Finish(request, response)
    }

    // And again, the prepare method has the ability to halt the response, so check again.
    if !response.isFinished {
        switch strings.ToUpper(r.Method) {
        case "HEAD":
            // If HEAD is not implemented, just trip the content from a regular GET request.
            response = handler.Head(request)
            if response.StatusCode == http.StatusMethodNotAllowed {
                getResp := handler.Get(request)
                if getResp.StatusCode == http.StatusOK {
                    response = getResp
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
        default:
            response = handler.ErrorHandler("Unsupported HTTP Method")
        }
    }

    // Always run postprocess for middlewares.
    runMiddlewarePostprocess(request, response)

    // Finish off the response by writing the output.
    writePatternResponse(response, w)

    runMixinFinish(handler)

    LogInfo.Printf("%d %s %s (%s) %s",
        response.StatusCode,
        strings.ToUpper(r.Method),
        r.RequestURI,
        r.RemoteAddr,
        time.Since(start_request))
}

func (ph *patHandler) try(path string) (url.Values, bool) {
    p := make(url.Values)
    var i, j int
    for i < len(path) {
        switch {
        case j >= len(ph.pat):
            if ph.pat != "/" && len(ph.pat) > 0 && ph.pat[len(ph.pat)-1] == '/' {
                return p, true
            }
            return nil, false
        case ph.pat[j] == ':':
            var name, val string
            var nextc byte
            name, nextc, j = match(ph.pat, isAlnum, j+1)
            val, _, i = match(path, matchPart(nextc), i)
            p.Add(":"+name, val)
        case path[i] == ph.pat[j]:
            i++
            j++
        default:
            return nil, false
        }
    }
    if j != len(ph.pat) {
        return nil, false
    }
    return p, true
}

func matchPart(b byte) func(byte) bool {
    return func(c byte) bool {
        return c != b && c != '/'
    }
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
    j = i
    for j < len(s) && f(s[j]) {
        j++
    }
    if j < len(s) {
        next = s[j]
    }
    return s[i:j], next, j
}

func isAlpha(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
    return isAlpha(ch) || isDigit(ch)
}

func writePatternResponse(response *HttpResponse, w http.ResponseWriter) {
    for k, v := range response.Header {
        w.Header().Set(k, strings.Join(v, ","))
    }

    w.Header().Set("Content-Type", response.ContentType)

    w.WriteHeader(response.StatusCode)
    w.Write([]byte(response.Content))
}
