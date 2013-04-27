package tango

import (
    "crypto/tls"
    "github.com/cojac/mux"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "net/url"
    "strings"
)

type HttpRequest struct {
    RawRequest *http.Request
    Body       string
}

func NewHttpRequest(orig *http.Request) *HttpRequest {
    r := new(HttpRequest)
    r.RawRequest = orig

    body, err := ioutil.ReadAll(orig.Body)
    if err != nil {
        panic("Error reading body")
    }
    r.Body = string(body)

    return r
}

// Info methods.
func (r HttpRequest) IsSecure() bool {
    if r.RawRequest.URL.Scheme == "https" {
        return true
    }

    if strings.HasPrefix(r.RawRequest.Proto, "HTTPS") {
        return true
    }

    if r.RawRequest.Header.Get("X-Forwarded-Proto") == "https" {
        return true
    }

    return false
}

func (r HttpRequest) IsAjax() bool {
    xhr, ok := r.RawRequest.Header["X-Requested-With"]
    if !ok {
        return false
    }

    for _, v := range xhr {
        if strings.ToLower(v) == "xmlhttprequest" {
            return true
        }
    }

    return false
}

// Retrive specific values.
func (r HttpRequest) PathValue(key string) (string, bool) {
    val, ok := mux.Vars(r.RawRequest)[key]
    return val, ok
}

func (r HttpRequest) GetValue(key string) (string, bool) {
    val, ok := r.RawRequest.URL.Query()[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r HttpRequest) GetArray(key string) ([]string, bool) {
    val, ok := r.RawRequest.URL.Query()[key]
    return val, ok
}

func (r HttpRequest) FormValue(key string) (string, bool) {
    val, ok := r.RawRequest.Form["key"]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r HttpRequest) FormArray(key string) ([]string, bool) {
    val, ok := r.RawRequest.Form["key"]
    return val, ok
}

func (r HttpRequest) FragmentValue() string {
    return r.RawRequest.URL.Fragment
}

// Why does it feel that there should be a better way to do this...
func (r HttpRequest) Proto() string {
    return r.RawRequest.Proto
}

func (r HttpRequest) ProtoMajor() int {
    return r.RawRequest.ProtoMajor
}

func (r HttpRequest) ProtoMinor() int {
    return r.RawRequest.ProtoMinor
}

func (r HttpRequest) Header() http.Header {
    return r.RawRequest.Header
}

func (r HttpRequest) ContentLength() int64 {
    return r.RawRequest.ContentLength
}

func (r HttpRequest) TransferEncoding() []string {
    return r.RawRequest.TransferEncoding
}

func (r HttpRequest) Close() bool {
    return r.RawRequest.Close
}

func (r HttpRequest) Host() string {
    return r.RawRequest.Host
}

func (r HttpRequest) Form() url.Values {
    return r.RawRequest.Form
}

func (r HttpRequest) MultipartForm() *multipart.Form {
    return r.RawRequest.MultipartForm
}

func (r HttpRequest) Trailer() http.Header {
    return r.RawRequest.Trailer
}

func (r HttpRequest) RemoteAddr() string {
    return r.RawRequest.RemoteAddr
}

func (r HttpRequest) RequestURI() string {
    return r.RawRequest.RequestURI
}

func (r HttpRequest) TLS() *tls.ConnectionState {
    return r.RawRequest.TLS
}
