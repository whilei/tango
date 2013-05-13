package tango

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
)

var LogInfo = log.New(os.Stdout, "[Tango I] ", log.Ldate|log.Ltime)
var LogError = log.New(os.Stderr, "[Tango E] ", log.Ldate|log.Ltime|log.Lshortfile)

var Version = "0.0.1"

func VersionMap() [3]int {
    var out [3]int
    t := strings.Split(Version, ".")
    for k, v := range t {
        i, _ := strconv.ParseInt(v, 10, 0)
        out[k] = int(i)
    }
    return out
}

func ListenAndServe() {
    // Lets leave this function bare bones... then App Engine can do everything
    // except call this function. (So call this in your main func)
    addr := Settings.String("serve_address", ":8000")
    LogInfo.Printf("Starting server at %s.", addr)

    http.ListenAndServe(addr, nil)
}
