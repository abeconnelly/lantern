package main

import "io"
import "fmt"
import "net/http"

import "github.com/abeconnelly/sloppyjson"

import "log"
import "bytes"

import "github.com/julienschmidt/httprouter"

import gouuid "github.com/nu7hatch/gouuid"

var g_req_debug bool = true

func send_error_bad_request(w http.ResponseWriter, err_str string) {
  if gVerboseFlag {
    log.Printf(err_str)
  }
  w.Header().Set("Content-Type", "application/json")
  io.WriteString(w, `{"Type":"error","Message":"bad request"}`)
  return
}

func ping_req_handler(w http.ResponseWriter, sj *sloppyjson.SloppyJSON, conn_info *LanternConnInfo) {
  io.WriteString(w, `{"Type":"pong"}`)

  if gVerboseFlag {
    log.Printf("%s: sent pong", conn_info.UUID.String())
  }
}

func stat_req_handler(w http.ResponseWriter, sj *sloppyjson.SloppyJSON, conn_info *LanternConnInfo) {
  io.WriteString(w, `{"Type":"stat-resp"`)
  io.WriteString(w, fmt.Sprintf(`,"Requests":%d`, gLanternStat.Requests))
  io.WriteString(w, fmt.Sprintf(`,"Failures":%d`, gLanternStat.Failures))
  io.WriteString(w, fmt.Sprintf(`,"AvgRespMs":%v`, gLanternStat.AvgRespMs))
  io.WriteString(w, fmt.Sprintf(`,"AvgRespDt":%f`, gLanternStat.AvgRespDt))
  io.WriteString(w, "}")

  if gVerboseFlag {
    log.Printf("%s: stat-resp", conn_info.UUID.String())
  }

}

func handle_json_req(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  var body_reader io.Reader = r.Body

  b := bytes.Buffer{}
  n,e := b.ReadFrom(body_reader) ; _ = n
  if e!=nil {
    send_error_bad_request(w, fmt.Sprintf("%v", e))
    return
  }

  if g_req_debug {
    log.Printf(b.String())
  }

  sj,e := sloppyjson.Loads(b.String())
  if e!=nil {
    send_error_bad_request(w, fmt.Sprintf("%v", e))
    return
  }

  conn_info := LanternConnInfo{}
  conn_info.UUID,e = gouuid.NewV4()
  if e!=nil {
    send_error_bad_request(w, fmt.Sprintf("%v", e))
    return
  }


  if _,ok := sj.O["Type"] ; !ok {
    send_error_bad_request(w, fmt.Sprintf("%v", e))
    return
  }

  if gVerboseFlag {
    log.Printf("%s: received %s\n", conn_info.UUID.String(), sj.O["Type"].S)
  }

  switch sj.O["Type"].S {
  case "ping":
    ping_req_handler(w, sj, &conn_info)
    return
  case "stat":
    stat_req_handler(w, sj, &conn_info)
    return
  }


}
