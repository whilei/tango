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
    in.Header.Add("X-Forwarded-Proto", "httPS")
    req = NewHttpRequest(in, make(url.Values))
    assert.Equal(t, true, req.IsSecure())
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

func TestRequestGetValueTypes(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com?one=foo&two=222&three=11&four=2.0", nil)
    req := NewHttpRequest(in, make(url.Values))

    vs := req.GetString("one")
    assert.Equal(t, "foo", vs)

    vs = req.GetString("two")
    assert.Equal(t, "222", vs)

    vi := req.GetInt("three")
    assert.Equal(t, int64(11), vi)

    vf := req.GetFloat("four")
    assert.Equal(t, 2.0, vf)

    vs = req.GetString("five", "hello")
    assert.Equal(t, "hello", vs)
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

func TestRequestFormValueTypes(t *testing.T) {
    in, _ := http.NewRequest("GET", "example.com", nil)
    in.PostForm = make(url.Values)
    in.PostForm.Add("one", "foo")
    in.PostForm.Add("two", "222")
    in.PostForm.Add("three", "11")
    in.PostForm.Add("four", "2.0")
    req := NewHttpRequest(in, make(url.Values))

    vs := req.FormString("one")
    assert.Equal(t, "foo", vs)

    vs = req.FormString("two")
    assert.Equal(t, "222", vs)

    vi := req.FormInt("three")
    assert.Equal(t, int64(11), vi)

    vf := req.FormFloat("four")
    assert.Equal(t, 2.0, vf)

    vs = req.FormString("five", "hello")
    assert.Equal(t, "hello", vs)

    vs = req.FormString("six")
    assert.Equal(t, "", vs)

    vi = req.FormInt("seven")
    assert.Equal(t, int64(0), vi)

    vf = req.FormFloat("eight")
    assert.Equal(t, 0.0, vf)

    vi = req.FormInt("nine", 22)
    assert.Equal(t, int64(22), vi)

    vf = req.FormFloat("ten", 10.10)
    assert.Equal(t, 10.10, vf)
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
