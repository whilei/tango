package middleware

import (
    "fmt"
    "github.com/cojac/context"
    "github.com/cojac/tango"
    "time"
)

const runTimeContextKey string = "__middleware_run_time_profile_start"

type RuntimeProfile struct {
    tango.BaseMiddleware
}

func (m RuntimeProfile) ProcessRequest(request *tango.HttpRequest) {
    context.Set(request.Raw, runTimeContextKey, time.Now())
}

func (m RuntimeProfile) ProcessResponse(request *tango.HttpRequest, response *tango.HttpResponse) {
    started := context.Get(request.Raw, runTimeContextKey)
    response.AddHeader("X-Runtime", fmt.Sprintf("%s", time.Since(started.(time.Time))))
}
