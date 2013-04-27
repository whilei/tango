package tango

import (
    "fmt"
    "strings"
)

type dictObj struct {
    data map[string]interface{}
}

var Debug = false
var Settings = dictObj{data: make(map[string]interface{})}

func (s *dictObj) Set(key string, val interface{}) {
    s.data[key] = val

    // Special case... only applies to debug!
    if strings.ToLower(key) == "debug" {
        Debug = s.Bool(key)
    }
}

func (s *dictObj) Bool(key string, args ...bool) bool {
    def := false

    switch len(args) {
    case 0:
        break
    case 1:
        def = bool(args[0])
    default:
        panic(fmt.Sprintf("Bool received too many args: [%d]", len(args)))
    }

    x, ok := s.data[key]
    if !ok {
        return def
    }
    return x.(bool)
}

func (s *dictObj) Int(key string, args ...int) int {
    var def int = -1

    switch len(args) {
    case 0:
        break
    case 1:
        def = args[0]
    default:
        panic(fmt.Sprintf("Int received too many args: [%d]", len(args)))
    }

    x, ok := s.data[key]
    if !ok {
        return def
    }
    return int(x.(float64))
}

func (s *dictObj) Float(key string, args ...float64) float64 {
    var def float64 = -1

    switch len(args) {
    case 0:
        break
    case 1:
        def = args[0]
    default:
        panic(fmt.Sprintf("Float received too many args: [%d]", len(args)))
    }

    x, ok := s.data[key]
    if !ok {
        return def
    }
    return x.(float64)
}

func (s *dictObj) String(key string, args ...string) string {
    var def string
    switch len(args) {
    case 0:
        break
    case 1:
        def = args[0]
    default:
        panic(fmt.Sprintf("String received too many args: [%d]", len(args)))
    }

    result, present := s.data[key]
    if !present {
        return def
    }
    return result.(string)
}
