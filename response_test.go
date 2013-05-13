package tango

import (
    "github.com/bmizerany/assert"
    "testing"
)

func TestResponseNewNoParams(t *testing.T) {
    resp := NewHttpResponse()

    assert.Equal(t, "", resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, false, resp.isFinished)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewOneParam(t *testing.T) {
    resp := NewHttpResponse("My Content")

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, false, resp.isFinished)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewTwoParams(t *testing.T) {
    resp := NewHttpResponse("My Content", 201)

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 201, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, false, resp.isFinished)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewThreeParams(t *testing.T) {
    resp := NewHttpResponse("My Content", 202, "text/plain; charset=latin1")

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 202, resp.StatusCode)
    assert.Equal(t, "text/plain; charset=latin1", resp.ContentType)
    assert.Equal(t, false, resp.isFinished)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewToManyParams(t *testing.T) {
    assert.Panic(t, "NewHttpResponse received [4] args, can only handle 3.", func() {
        NewHttpResponse("My Content", 202, "text/plain", "oops")
    })
}

func TestResponseIsFinished(t *testing.T) {
    resp := NewHttpResponse()

    assert.Equal(t, false, resp.isFinished)
    resp.Finish()
    assert.Equal(t, true, resp.isFinished)
}
