package tango

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// ---
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
                handler := ph.New()
                response := handler.PermanentRedirect(NewHttpRequest(r, params), buildUrlWithOutSlash(r))
                writePatternResponse(response, w)
                LogInfo.Printf("%d %s %s (%s) -",
                    response.StatusCode,
                    strings.ToUpper(r.Method),
                    r.RequestURI,
                    r.RemoteAddr,
                )
                return
            }

            ph.ServeHandlerHttp(w, r, params)
            return
        }
    }

    notFoundPatHandler.ServeHandlerHttp(w, r, nil)
}

func (p *PatternServeMux) ServeTestResponse(r *http.Request) *HttpResponse {
    for _, ph := range p.handlers {
        if params, ok := ph.try(r.URL.Path); ok {
            if ph.isSlashRedirect {
                handler := ph.New()
                return handler.PermanentRedirect(NewHttpRequest(r, params), buildUrlWithOutSlash(r))
            }

            return ph.processRequest(r, params)
        }
    }

    return notFoundPatHandler.processRequest(r, nil)
}

func Pattern(pat string, h HandlerInterface) {
    Mux.handlers = append(Mux.handlers, &patHandler{pat, h, false})
    if Settings.Bool("append_slash", false) {
        // Allows urls like "/api/health" to be called like "/api/health/".
        n := len(pat)
        if n > 0 && pat[n-1] != '/' {
            Mux.handlers = append(Mux.handlers, &patHandler{pat + "/", h, Settings.Bool("append_slash_should_redirect", true)})
        }
    }
}

func buildUrlWithOutSlash(r *http.Request) string {
    result := r.URL.Path

    // Stripe slash.
    n := len(r.URL.Path)
    if n > 0 && r.URL.Path[n-1] == '/' {
        result = r.URL.Path[:n-1]
    }

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

func (ph *patHandler) processRequest(r *http.Request, params url.Values) *HttpResponse {
    handler := ph.New()

    var response *HttpResponse

    func() {
        // Any panic errors will be caught and passed over to our ErrorHandler.
        defer func() {
            if Debug == false {
                if rec := recover(); rec != nil {
                    LogError.Printf("Panic Recovered: %s", rec)
                    response = handler.ErrorHandler(fmt.Sprintf("%q", rec))
                }
            }
        }()

        response = NewHttpResponse()
        finished := false
        request := NewHttpRequest(r, params)

        runMixinPrepare(handler)
        defer runMixinFinish(handler)

        midResp := runMiddlewarePreprocess(request)
        if midResp != nil {
            finished = true
            response = midResp
        }

        // Only if the response has not finished should we let the handler touch it.
        if !finished {
            prepResp := handler.Prepare(request)
            if prepResp != nil {
                finished = true
                response = prepResp
            }
        }

        // And again, the prepare method has the ability to halt the response, so check again.
        if !finished {
            switch strings.ToUpper(r.Method) {
            case "HEAD":
                // If HEAD is not implemented, just strip the content from a regular GET request.
                response = handler.Head(request)
                if response == nil {
                    response = handler.Get(request)
                    response.Content = ""
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

        if response == nil {
            panic("Response cannot be nil.")
        }

        // Always call finish before middleware.
        handler.Finish(request, response)

        // Always run postprocess for middlewares.
        runMiddlewarePostprocess(request, response)
    }()

    return response
}

func (ph *patHandler) ServeHandlerHttp(w http.ResponseWriter, r *http.Request, params url.Values) {
    start_request := time.Now()

    response := ph.processRequest(r, params)

    // Finish off the response by writing the output.
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

// 404 - Not Found
type NotFoundHandler struct{ BaseHandler }

func (h *NotFoundHandler) New() HandlerInterface {
    return &NotFoundHandler{}
}
func (h *NotFoundHandler) Get(request *HttpRequest) *HttpResponse {
    return NewHttpResponse("404 - Page not found.", http.StatusNotFound)
}
func (h *NotFoundHandler) Post(request *HttpRequest) *HttpResponse {
    return h.Get(request)
}
func (h *NotFoundHandler) Put(request *HttpRequest) *HttpResponse {
    return h.Get(request)
}
func (h *NotFoundHandler) Patch(request *HttpRequest) *HttpResponse {
    return h.Get(request)
}
func (h *NotFoundHandler) Delete(request *HttpRequest) *HttpResponse {
    return h.Get(request)
}
func (h *NotFoundHandler) Options(request *HttpRequest) *HttpResponse {
    return h.Get(request)
}

var notFoundPatHandler *patHandler = &patHandler{"404", &NotFoundHandler{}, false}

func SetNotFoundHandler(h HandlerInterface) {
    notFoundPatHandler = &patHandler{"404", h, false}
}
