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

// Tail returns the trailing string in path after the final slash for a pat ending with a slash.
//
// Examples:
//
//  Tail("/hello/:title/", "/hello/mr/mizerany") == "mizerany"
//  Tail("/:a/", "/x/y/z")                       == "y/z"
//
func Tail(pat, path string) string {
    var i, j int
    for i < len(path) {
        switch {
        case j >= len(pat):
            if pat[len(pat)-1] == '/' {
                return path[j-1:]
            }
            return ""
        case pat[j] == ':':
            var nextc byte
            _, nextc, j = match(pat, isAlnum, j+1)
            _, _, i = match(path, matchPart(nextc), i)
        case path[i] == pat[j]:
            i++
            j++
        default:
            return ""
        }
    }
    return ""
}

type patHandler struct {
    pat string
    HandlerInterface
    isSlashRedirect bool
}

func (ph *patHandler) ServeHandlerHttp(w http.ResponseWriter, r *http.Request, params url.Values) {
    // Any panic errors will be caught and passed over to our ErrorHandler.
    defer func() {
        if rec := recover(); rec != nil {
            LogError.Printf("Panic Recovered: %s", rec)
            writePatternResponse(ph.ErrorHandler(fmt.Sprintf("%q", rec)), w)
        }
    }()

    start_request := time.Now()

    request := NewHttpRequest(r, params)
    RunMiddlewarePreprocess(request)
    ph.Prepare(request)

    var response *HttpResponse

    switch strings.ToUpper(r.Method) {
    case "HEAD":
        response = ph.Head(request)
        if response.StatusCode == http.StatusMethodNotAllowed {
            resp2 := ph.Get(request)
            if resp2.StatusCode == http.StatusOK {
                response = resp2
                response.Content = ""
            }
        }
    case "GET":
        response = ph.Get(request)
    case "POST":
        response = ph.Post(request)
    case "PUT":
        response = ph.Put(request)
    case "PATCH":
        response = ph.Patch(request)
    case "DELETE":
        response = ph.Delete(request)
    case "OPTIONS":
        response = ph.Options(request)
    }

    ph.Finish(request, response)
    RunMiddlewarePostprocess(request, response)

    writePatternResponse(response, w)

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
