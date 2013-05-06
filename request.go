package tango

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

type HttpRequest struct {
    http.Request
    BodyString string
    PathValues url.Values
    Args       map[string]interface{} //Subject to a better name!
}

func NewHttpRequest(orig *http.Request, params url.Values) *HttpRequest {
    r := HttpRequest{*orig, "", params, make(map[string]interface{})}

    if r.Body != nil {
        strBody, err := ioutil.ReadAll(r.Body)
        if err != nil {
            panic("Error reading body")
        }
        r.BodyString = string(strBody)
    }

    return &r
}

//Info methods.
func (r HttpRequest) IsSecure() bool {
    if r.URL.Scheme == "https" {
        return true
    }

    if strings.HasPrefix(r.Proto, "HTTPS") {
        return true
    }

    if r.Header.Get("X-Forwarded-Proto") == "https" {
        return true
    }

    return false
}

func (r HttpRequest) IsAjax() bool {
    xhr, ok := r.Header["X-Requested-With"]
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
    val, ok := r.PathValues[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r HttpRequest) GetValue(key string) (string, bool) {
    val, ok := r.URL.Query()[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r HttpRequest) GetArray(key string) ([]string, bool) {
    val, ok := r.URL.Query()[key]
    return val, ok
}

func (r HttpRequest) FormValue(key string) (string, bool) {
    val, ok := r.Form[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r HttpRequest) FormArray(key string) ([]string, bool) {
    val, ok := r.Form[key]
    return val, ok
}

func (r HttpRequest) FragmentValue() string {
    return r.URL.Fragment
}
