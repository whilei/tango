package middleware

import (
    "fmt"
    "github.com/cojac/tango"
    "time"
)

const runTimeContextKey string = "__middleware_run_time_profile_start__"

type RuntimeProfile struct {
    tango.BaseMiddleware
}

func (m RuntimeProfile) ProcessRequest(request *tango.HttpRequest, response *tango.HttpResponse) {
    request.Registry[runTimeContextKey] = time.Now()
}

func (m RuntimeProfile) ProcessResponse(request *tango.HttpRequest, response *tango.HttpResponse) {
    started := request.Registry[runTimeContextKey]
    response.Header.Set("X-Runtime", fmt.Sprintf("%s", time.Since(started.(time.Time))))
}
