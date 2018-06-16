package tango

var middlewares = []MiddlewareInterface{}

type MiddlewareInterface interface {
	ProcessRequest(request *HttpRequest) *HttpResponse
	ProcessResponse(request *HttpRequest, response *HttpResponse)
}

type BaseMiddleware struct{}

func (m *BaseMiddleware) ProcessRequest(request *HttpRequest) *HttpResponse {
	return nil
}

func (m *BaseMiddleware) ProcessResponse(request *HttpRequest, response *HttpResponse) {
	// pass
}

func Middleware(m MiddlewareInterface) {
	middlewares = append(middlewares, m)
}

func runMiddlewarePreprocess(request *HttpRequest) *HttpResponse {
	// Top to bottom.
	for i := 0; i < len(middlewares); i++ {
		m := middlewares[i]
		r := m.ProcessRequest(request)
		if r != nil {
			return r
		}
	}
	return nil
}

func runMiddlewarePostprocess(request *HttpRequest, response *HttpResponse) {
	// Bottom to top.
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		m.ProcessResponse(request, response)
	}
}
