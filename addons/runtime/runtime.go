package runtime

import (
	"fmt"
	"time"

	"github.com/unrolled/tango"
)

const runTimeContextKey string = "__middleware_run_time_profiler_key__"

type Profiler struct {
	tango.BaseMiddleware
}

func (m *Profiler) ProcessRequest(request *tango.HttpRequest) *tango.HttpResponse {
	request.Registry[runTimeContextKey] = time.Now()

	return nil
}

func (m *Profiler) ProcessResponse(request *tango.HttpRequest, response *tango.HttpResponse) {
	started := request.Registry[runTimeContextKey]
	response.Header.Set("X-Runtime", fmt.Sprintf("%s", time.Since(started.(time.Time))))
}
