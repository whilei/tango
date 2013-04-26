package tango

import (
    "log"
    "net/http"
    "os"
)

var LogInfo = log.New(os.Stdout, "[Tango I] ", log.Ldate|log.Ltime)
var LogError = log.New(os.Stderr, "[Tango E] ", log.Ldate|log.Ltime|log.Lshortfile)

func ListenAndServe() {
    LogInfo.Print("Starting server at :8000.")

    http.Handle("/", Mux)
    http.ListenAndServe(":8000", nil)
}
