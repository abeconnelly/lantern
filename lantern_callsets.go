package main

import "io"
import "fmt"
import "log"
import "strings"
import "strconv"
import "net/http"
import "crypto/md5"
import "github.com/abeconnelly/cgf"
import "github.com/julienschmidt/httprouter"

func (ctx *LanternContext) APICallsets(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
  if ctx.VerboseFlag {
    log.Printf("APICallsets\n")
  }

  count:=0
  io.WriteString(w, `[`)
  for x,_ := range ctx.CGFIndexMap {
    if count>0 { io.WriteString(w, `,`) }
    io.WriteString(w, fmt.Sprintf(`"%s"`, x))
    count++
  }
  io.WriteString(w, `]`)
}

func (ctx *LanternContext) APICallsetsId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
  callset_name := param.ByName("callset_id")

  if ctx.VerboseFlag {
    log.Printf("APICallsetsId callset_id %v\n", callset_name)
  }

  if _,ok := ctx.CGFIndexMap[callset_name] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  callset_idx := ctx.CGFIndexMap[callset_name] ; _ = callset_idx
  callset_loc := ctx.Config.O["cgf"].L[callset_idx].O["locator"].S

  io.WriteString(w, `{`)
  io.WriteString(w, fmt.Sprintf(`"callset-name":"%s",`,callset_name))
  io.WriteString(w, fmt.Sprintf(`"callset-locator":"%s"`,callset_loc))
  io.WriteString(w, `}`)

}

func (ctx *LanternContext) APICallsetsIdTileVariants(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  callset_name := param.ByName("callset_id")
  tilepos := r.FormValue("tile-positions")

  if ctx.VerboseFlag {
    log.Printf("APICallsetsIdTileVariants callset_id %v, tile-positions %v\n", callset_name, tilepos)
  }

  if _,ok := ctx.CGFIndexMap[callset_name] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  tile_parts := strings.Split(tilepos, ".")
  if len(tile_parts)!=3 {
    send_error_bad_request(w, "invalid tilepos")
    return
  }

  _req_tagset,e := strconv.ParseInt(tile_parts[0], 16, 64)
  if e!=nil {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  _req_path,e := strconv.ParseInt(tile_parts[1], 16, 64)
  if e!=nil {
    send_error_bad_request(w, "invalid path id")
    return
  }

  _req_step,e := strconv.ParseInt(tile_parts[2], 16, 64)
  if e!=nil {
    send_error_bad_request(w, "invalid step id")
    return
  }

  req_tagset := int(_req_tagset) ; _ = req_tagset
  req_path := int(_req_path) ; _ = req_path
  req_step := int(_req_step) ; _ = req_step

  callset_idx := ctx.CGFIndexMap[callset_name]
  knot := cgf.GetKnot(ctx.TileMap, ctx.CGFPath[callset_idx][req_path], req_step)

  io.WriteString(w, `{`)
  io.WriteString(w, `"callset-name":"` + callset_name + `",`)
  io.WriteString(w, `"tile-variants":[`)

  count:=0
  for a:=0; a<len(knot); a++ {
    cur_step := knot[a][0].Step
    if a>0 {
      io.WriteString(w, `,[`)
    } else {
      io.WriteString(w, `[`)
    }
    for i:=0; i<len(knot[a]); i++ {
      varid := knot[a][i].VarId
      _seq := ctx.SGLF.Lib[req_path][cur_step][varid]
      seq := cgf.FillNocSeq(_seq, knot[a][i].NocallStartLen)

      m5str := cgf.Md5sum2str(md5.Sum([]byte(seq)))

      loq_flag := ""
      if len(knot[a][i].NocallStartLen)>0 {
        loq_flag = "*"
      }

      tilestr := fmt.Sprintf("%02x.%04x.%04x.%s+%x%s",
        req_tagset, req_path, cur_step, m5str, knot[a][i].Span, loq_flag)

      if i>0 { io.WriteString(w, `,`) }
      io.WriteString(w, `"` + tilestr + `"`)

      cur_step += knot[a][i].Span
      count++
    }
    io.WriteString(w, `]`)
  }
  io.WriteString(w, `]`)
  io.WriteString(w, `}`)

}
