package tango

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
)

// Special Debug variable... used anywhere and everywhere, so let's make it easy.
var Debug = false

type DictObj struct {
    data map[string]interface{}
}

// Convention will be json style settings. Try to keep it a single level deep.
// Key values should be lower_case_with_underscores. Ex: debug, or secret_key.
var Settings = DictObj{data: make(map[string]interface{})}

func (s *DictObj) LoadFromFile(filepath string) {
    f, err := os.Open(filepath)
    if err != nil {
        panic(fmt.Sprintf("Opening settings file failed: %s", err))
    }
    defer f.Close()

    b := new(bytes.Buffer)
    _, err = b.ReadFrom(f)
    if err != nil {
        panic(fmt.Sprintf("Reading settings file failed: %s", err))
    }

    err = json.Unmarshal(b.Bytes(), &s.data)
    if err != nil {
        panic(fmt.Sprintf("Unpacking settings json failed: %s", err))
    }
}

func (s *DictObj) Set(key string, val interface{}) {
    s.data[key] = val

    // Special case... only applies to our global Debug!
    if key == "debug" {
        Debug = s.Bool(key)
    }
}

func (s *DictObj) Bool(key string, args ...bool) bool {
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

func (s *DictObj) Int(key string, args ...int) int {
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

func (s *DictObj) Float(key string, args ...float64) float64 {
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

func (s *DictObj) String(key string, args ...string) string {
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

func (s *DictObj) Array(key string, args ...[]interface{}) []interface{} {
    def := []interface{}(nil)

    switch len(args) {
    case 0:
        break
    case 1:
        def = args[0]
    default:
        panic(fmt.Sprintf("Array received too many args: [%d]", len(args)))
    }

    result, present := s.data[key]
    if !present {
        return def
    }
    return result.([]interface{})
}
