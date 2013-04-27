package tango

import (
    "log"
    "net/http"
    "os"
)

var LogInfo = log.New(os.Stdout, "[Tango I] ", log.Ldate|log.Ltime)
var LogError = log.New(os.Stderr, "[Tango E] ", log.Ldate|log.Ltime|log.Lshortfile)

func ListenAndServe() {
    addr := Settings.String("serve_address")
    LogInfo.Printf("Starting server at %s.", addr)

    http.Handle("/", Mux)
    http.ListenAndServe(addr, nil)
}
