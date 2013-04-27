package tango

import (
    "fmt"
    "strings"
)

type settingsObj struct {
    data map[string]interface{}
}

var Debug = false
var Settings = settingsObj{data: make(map[string]interface{})}

func (s *settingsObj) Set(key string, val interface{}) {
    s.data[key] = val

    // Special case... only applies to debug!
    if strings.ToLower(key) == "debug" {
        Debug = s.Bool(key)
    }
}

func (s *settingsObj) Bool(key string, args ...bool) bool {
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

func (s *settingsObj) Int(key string, args ...int) int {
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

func (s *settingsObj) Float(key string, args ...float64) float64 {
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

func (s *settingsObj) String(args ...string) string {
    var key, def string
    switch len(args) {
    case 2:
        def = args[1]
        fallthrough
    case 1:
        key = args[0]
    case 0:
        break
    default:
        panic(fmt.Sprintf("String received too many args: [%d]", len(args)))
    }

    result, present := s.data[key]
    if !present {
        return def
    }
    return result.(string)
}
