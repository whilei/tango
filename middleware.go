package tango

var Middlewares = []MiddlewareInterface{}

type MiddlewareInterface interface {
    ProcessRequest(request *HttpRequest)
    ProcessResponse(request *HttpRequest, response *HttpResponse)
}

type BaseMiddleware struct{}

func (m BaseMiddleware) ProcessRequest(req *HttpRequest) {
    // pass
}

func (m BaseMiddleware) ProcessResponse(req *HttpRequest, resp *HttpResponse) {
    // pass
}

func Middleware(m MiddlewareInterface) {
    Middlewares = append(Middlewares, m)
}

func RunMiddlewarePreprocess(req *HttpRequest) {
    for i := 0; i < len(Middlewares); i++ {
        m := Middlewares[i]
        m.ProcessRequest(req)
    }
}

func RunMiddlewarePostprocess(req *HttpRequest, resp *HttpResponse) {
    for i := len(Middlewares) - 1; i >= 0; i-- {
        m := Middlewares[i]
        m.ProcessResponse(req, resp)
    }
}
