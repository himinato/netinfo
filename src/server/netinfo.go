package main

import (
  "fmt"
  "regexp"
  "encoding/json"
  "net/http"
  "net"
  "os"
)

type IP string

func (ip IP) isIPv4() bool {
  return regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`).Match([]byte(ip))
}

func getIp(req *http.Request) string {
  if req.Header.Get("X-FORWARDED-FOR") != "" {
    return req.Header.Get("X-FORWARDED-FOR")
  } else if req.RemoteAddr != "" {
    ip, _, _ := net.SplitHostPort(req.RemoteAddr)
    return ip
  }
  return ""
}

func backip (w http.ResponseWriter, r *http.Request) {
  ip := IP(getIp(r))

  if ip.isIPv4() {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "%v", ip)
  } else {
    http.Error(w, "can't detect IP", http.StatusInternalServerError)
  }
}

func backjson (w http.ResponseWriter, r *http.Request) {
  ip := IP(getIp(r))

  r.Header.Set("ip", fmt.Sprint(ip))
  js, err := json.Marshal(r.Header)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if ip.isIPv4() {
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
  } else {
    http.Error(w, "can't detect IP", http.StatusInternalServerError)
  }
}


func main() {
  http.HandleFunc("/ip", backip)
  http.HandleFunc("/json", backjson)

  port := os.Getenv("PORT")
  if port == "" {
    port = "80"
  }
  err := http.ListenAndServe(":" + port, nil)
  if err != nil {
    panic(err)
  }
}
