package main

import "net/http"
import "io"
import "log"

var gAPIVersionString string = "0.1.0"

func handle_status(w http.ResponseWriter, r *http.Request) {
  if gVerboseFlag { log.Printf("status (%s)", gAPIVersionString) }
  io.WriteString(w,`{"api-version":"` + gAPIVersionString + `"}`)
}


