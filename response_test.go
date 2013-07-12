package tango

import (
    "github.com/cojac/assert"
    "testing"
)

func TestResponseNewNoParams(t *testing.T) {
    resp := NewHttpResponse()

    assert.Equal(t, "", resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewOneParam(t *testing.T) {
    resp := NewHttpResponse("My Content")

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewTwoParams(t *testing.T) {
    resp := NewHttpResponse("My Content", 201)

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 201, resp.StatusCode)
    assert.Equal(t, "text/html; charset=utf-8", resp.ContentType)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewThreeParams(t *testing.T) {
    resp := NewHttpResponse("My Content", 202, "text/plain; charset=latin1")

    assert.Equal(t, "My Content", resp.Content)
    assert.Equal(t, 202, resp.StatusCode)
    assert.Equal(t, "text/plain; charset=latin1", resp.ContentType)
    assert.Equal(t, 0, len(resp.Header))
    assert.Equal(t, 0, len(resp.Context))
}

func TestResponseNewToManyParams(t *testing.T) {
    assert.Panic(t, "NewHttpResponse received [4] args, can only handle 3.", func() {
        NewHttpResponse("My Content", 202, "text/plain", "oops")
    })
}

func TestResponseNewArray(t *testing.T) {
    resp := NewHttpResponse([1]string{"hellowww"})

    assert.Equal(t, `["hellowww"]`, resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "application/json; charset=utf-8", resp.ContentType)
}

func TestResponseNewMap(t *testing.T) {
    v := map[string]interface{}{"k": "value"}
    resp := NewHttpResponse(v)

    assert.Equal(t, `{"k":"value"}`, resp.Content)
    assert.Equal(t, 200, resp.StatusCode)
    assert.Equal(t, "application/json; charset=utf-8", resp.ContentType)
}
