package tango

import (
    "bytes"
    "github.com/cojac/assert"
    "net/http"
    "net/url"
    "testing"
)

func TestRequestNew(t *testing.T) {
    b := bytes.NewBufferString("sample body here")
    in, _ := http.NewRequest("GET", "http://example.com", b)
    req := NewHttpRequest(in, make(url.Values))

    assert.Equal(t, "sample body here", req.BodyString)
    assert.Equal(t, 0, len(req.PathArgs))
    assert.Equal(t, 0, len(req.Registry))
    assert.Equal(t, "example.com", req.URL.Host)
    assert.Equal(t, "GET", req.Method)
}

func TestRequestIsSecure(t *testing.T) {
    in, _ := http.NewRequest("GET", "http://example.com", nil)
    req := NewHttpRequest(in, make(url.Values))
    assert.Equal(t, false, req.IsSecure())

    in, _ = http.NewRequest("GET", "https://secure-example.com", nil)
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, true, req.IsSecure())

    in, _ = http.NewRequest("GET", "example.com", nil)
    in.Proto = "HTTPS/1.1"
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, true, req.IsSecure())

    in, _ = http.NewRequest("GET", "example.com", nil)
    in.Header.Add("X-Forwarded-Proto", "httPS")
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, true, req.IsSecure())
}

func TestRequestIsAjax(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com", nil)
    req := NewHttpRequest(in, make(url.Values))
    assert.Equal(t, false, req.IsAjax())

    in, _ = http.NewRequest("GET", "example.com", nil)
    in.Header.Add("X-Requested-With", "xmlsucks")
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, false, req.IsAjax())

    in, _ = http.NewRequest("GET", "example.com", nil)
    in.Header.Add("X-Requested-With", "xmlhttprequest")
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, true, req.IsAjax())
}

func TestRequestGetValue(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com?one=foo&two=bar&three=1&three=2&three=3", nil)
    req := NewHttpRequest(in, make(url.Values))

    v, ok := req.GetValue("one")
    assert.Equal(t, true, ok)
    assert.Equal(t, "foo", v)

    v, ok = req.GetValue("two")
    assert.Equal(t, true, ok)
    assert.Equal(t, "bar", v)

    v, ok = req.GetValue("foobar")
    assert.Equal(t, false, ok)
    assert.Equal(t, "", v)

    v, ok = req.GetValue("three")
    assert.Equal(t, true, ok)
    assert.Equal(t, "1", v)

    varray, ok := req.GetArray("three")
    assert.Equal(t, true, ok)
    assert.Equal(t, []string{"1", "2", "3"}, varray)

    varray, ok = req.GetArray("barfoo")
    assert.Equal(t, false, ok)
    assert.Equal(t, []string(nil), varray)
}

func TestRequestFragementValue(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com?foo=bar#batman", nil)
    req := NewHttpRequest(in, make(url.Values))
    assert.Equal(t, "batman", req.FragmentValue())

    in, _ = http.NewRequest("GET", "example.com?foo=bar", nil)
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, "", req.FragmentValue())
}

func TestRequestFormValue(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com", nil)
    in.PostForm = make(url.Values)
    in.PostForm.Add("one", "foo")
    in.PostForm.Add("two", "bar")
    in.PostForm.Add("three", "1")
    in.PostForm.Add("three", "2")
    in.PostForm.Add("three", "3")
    req := NewHttpRequest(in, make(url.Values))

    v, ok := req.FormValue("one")
    assert.Equal(t, true, ok)
    assert.Equal(t, "foo", v)

    v, ok = req.FormValue("two")
    assert.Equal(t, true, ok)
    assert.Equal(t, "bar", v)

    v, ok = req.FormValue("foobar")
    assert.Equal(t, false, ok)
    assert.Equal(t, "", v)

    v, ok = req.FormValue("three")
    assert.Equal(t, true, ok)
    assert.Equal(t, "1", v)

    varray, ok := req.FormArray("three")
    assert.Equal(t, true, ok)
    assert.Equal(t, []string{"1", "2", "3"}, varray)

    varray, ok = req.FormArray("barfoo")
    assert.Equal(t, false, ok)
    assert.Equal(t, []string(nil), varray)
}

func TestRequestPathValue(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com", nil)
    path := make(url.Values)
    path.Add(":one", "foo")
    path.Add(":two", "bar")
    req := NewHttpRequest(in, path)

    v, ok := req.PathValue(":one")
    assert.Equal(t, true, ok)
    assert.Equal(t, "foo", v)

    v, ok = req.PathValue(":two")
    assert.Equal(t, true, ok)
    assert.Equal(t, "bar", v)

    v, ok = req.PathValue(":foobar")
    assert.Equal(t, false, ok)
    assert.Equal(t, "", v)
}
