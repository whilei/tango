package tango

var Middlewares = []MiddlewareInterface{}

type MiddlewareInterface interface {
    ProcessRequest(request *HttpRequest, response *HttpResponse)
    ProcessResponse(request *HttpRequest, response *HttpResponse)
}

type BaseMiddleware struct{}

func (m BaseMiddleware) ProcessRequest(request *HttpRequest, response *HttpResponse) {
    // pass
}

func (m BaseMiddleware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
    // pass
}

func Middleware(m MiddlewareInterface) {
    Middlewares = append(Middlewares, m)
}

func RunMiddlewarePreprocess(request *HttpRequest, response *HttpResponse) {
    // Top to bottom.
    for i := 0; i < len(Middlewares); i++ {
        m := Middlewares[i]
        m.ProcessRequest(request, response)

        // If the middleware obj has set the response as finished, then no need to continue.
        if response.isFinished {
            return
        }
    }
}

func RunMiddlewarePostprocess(request *HttpRequest, response *HttpResponse) {
    // Bottom to top.
    for i := len(Middlewares) - 1; i >= 0; i-- {
        m := Middlewares[i]
        m.ProcessResponse(request, response)
    }
}
