package tango

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/smtp"
    "os"
    "strconv"
    "strings"
)

// This logger will out put with the prefix "[Tango D] " when Debug mode is true.
var LogDebug = log.New(ioutil.Discard, "", log.LstdFlags)

// Normal usage loggers.
var LogInfo = log.New(os.Stdout, "[Tango I] ", log.Ldate|log.Ltime)
var LogError = log.New(os.Stderr, "[Tango E] ", log.Ldate|log.Ltime|log.Lshortfile)

var Version = "0.0.2"

func VersionMap() [3]int {
    var out [3]int
    t := strings.Split(Version, ".")
    for k, v := range t {
        i, _ := strconv.ParseInt(v, 10, 0)
        out[k] = int(i)
    }
    return out
}

func SendMail(subject, message, to_address string) error {
    return SendMassMail(subject, message, []string{to_address})
}

func SendMailFrom(subject, message, to_address, from_address string) error {
    return SendMassMailFrom(subject, message, from_address, []string{to_address})
}

func SendMassMail(subject, message string, to_addresses []string) error {
    from_address := Settings.String("from_address", "tango_server@example.com")
    return SendMassMailFrom(subject, message, from_address, to_addresses)
}

func SendMassMailFrom(subject, message, from_address string, to_addresses []string) error {
    mail_server := Settings.String("mail_server", "localhost:25")
    body := fmt.Sprintf("Subject: %s\n%s", strings.TrimSpace(subject), message)
    return smtp.SendMail(mail_server, nil, from_address, to_addresses, []byte(body))
}

func ListenAndServe() {
    // Lets leave this function bare bones... then App Engine can do everything
    // except call this function. (So call this in your main func)
    addr := Settings.String("serve_address", ":8000")
    LogInfo.Printf("Starting server at %s.", addr)

    http.ListenAndServe(addr, nil)
}
