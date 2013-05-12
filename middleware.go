package tango

var middlewares = []MiddlewareInterface{}

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
    middlewares = append(middlewares, m)
}

func runMiddlewarePreprocess(request *HttpRequest, response *HttpResponse) {
    // Top to bottom.
    for i := 0; i < len(middlewares); i++ {
        m := middlewares[i]
        m.ProcessRequest(request, response)

        // If the middleware obj has set the response as finished, then no need to continue.
        if response.isFinished {
            return
        }
    }
}

func runMiddlewarePostprocess(request *HttpRequest, response *HttpResponse) {
    // Bottom to top.
    for i := len(middlewares) - 1; i >= 0; i-- {
        m := middlewares[i]
        m.ProcessResponse(request, response)
    }
}
