package main

import "net/http"
import "io"
import "log"

import "github.com/julienschmidt/httprouter"

func handle_status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  if gVerboseFlag { log.Printf("status (%s)", gAPIVersionString) }
  io.WriteString(w,`{"api-version":"` + gAPIVersionString + `"}`)
}


