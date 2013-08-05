package tango

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

type HttpRequest struct {
    http.Request
    BodyString string
    PathArgs   url.Values
    Registry   map[string]interface{}
}

func NewHttpRequest(orig *http.Request, params url.Values) *HttpRequest {
    r := &HttpRequest{*orig, "", params, make(map[string]interface{})}

    if r.Body != nil {
        strBody, err := ioutil.ReadAll(r.Body)
        if err != nil {
            panic("Error reading request body")
        }
        r.BodyString = string(strBody)
    }

    return r
}

//Info methods.
func (r *HttpRequest) IsSecure() bool {
    if strings.EqualFold(r.URL.Scheme, "https") {
        return true
    }

    if strings.HasPrefix(r.Proto, "HTTPS") {
        return true
    }

    if strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
        return true
    }

    return false
}

func (r *HttpRequest) IsAjax() bool {
    xhr, ok := r.Header["X-Requested-With"]
    if !ok {
        return false
    }

    for _, v := range xhr {
        if strings.EqualFold(v, "xmlhttprequest") {
            return true
        }
    }

    return false
}

// Retrive specific values.
func (r *HttpRequest) PathValue(key string) (string, bool) {
    val, ok := r.PathArgs[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r *HttpRequest) GetValue(key string) (string, bool) {
    val, ok := r.URL.Query()[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r *HttpRequest) GetString(key string, args ...string) string {
    v, o := r.GetValue(key)
    return convToString(v, o, args...)
}

func (r *HttpRequest) GetInt(key string, args ...int64) int64 {
    v, o := r.GetValue(key)
    return convToInt(v, o, args...)
}

func (r *HttpRequest) GetFloat(key string, args ...float64) float64 {
    v, o := r.GetValue(key)
    return convToFloat(v, o, args...)
}

func (r *HttpRequest) GetArray(key string) ([]string, bool) {
    val, ok := r.URL.Query()[key]
    return val, ok
}

func (r *HttpRequest) FormValue(key string) (string, bool) {
    val, ok := r.PostForm[key]
    if ok {
        return val[0], ok
    }
    return "", ok
}

func (r *HttpRequest) FormString(key string, args ...string) string {
    v, o := r.FormValue(key)
    return convToString(v, o, args...)
}

func (r *HttpRequest) FormInt(key string, args ...int64) int64 {
    v, o := r.FormValue(key)
    return convToInt(v, o, args...)
}

func (r *HttpRequest) FormFloat(key string, args ...float64) float64 {
    v, o := r.FormValue(key)
    return convToFloat(v, o, args...)
}

func (r *HttpRequest) FormArray(key string) ([]string, bool) {
    val, ok := r.PostForm[key]
    return val, ok
}

func (r *HttpRequest) FragmentValue() string {
    return r.URL.Fragment
}

func convToString(val string, ok bool, args ...string) string {
    if ok {
        return val
    } else if len(args) >= 1 {
        return args[0]
    }

    return ""
}

func convToInt(val string, ok bool, args ...int64) int64 {
    if ok {
        newInt, err := strconv.ParseInt(val, 10, 64)
        if err == nil {
            return newInt
        }
    } else if len(args) >= 1 {
        return args[0]
    }

    return 0
}

func convToFloat(val string, ok bool, args ...float64) float64 {
    if ok {
        newFloat, err := strconv.ParseFloat(val, 64)
        if err == nil {
            return newFloat
        }
    } else if len(args) >= 1 {
        return args[0]
    }

    return 0.0
}
