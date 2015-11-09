package main

import "fmt"
import "net/http"
import "io"
import "log"

import "github.com/julienschmidt/httprouter"
import "github.com/abeconnelly/cgf"

import "os"
import "runtime/pprof"

func (ctx *LanternContext) APIStatus(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  if ctx.VerboseFlag { log.Printf("status (%s)", gAPIVersionString) }

  if ctx.VerboseAPIFlag {
    io.WriteString(w,"{\n")
    io.WriteString(w,`  "api-version":"` + gAPIVersionString + `"` + ",\n")
    io.WriteString(w, fmt.Sprintf(`  "cgf-count":%d,` + "\n",  len(ctx.CGFBytes)))
    io.WriteString(w, fmt.Sprintf(`  "lantern-version":"%s"` + "\n",  gLanternVersion))
    io.WriteString(w,"}\n")

    ctx.CGFSimpleStats()

  } else {
    io.WriteString(w,`{"api-version":"` + gAPIVersionString + `"}`)
  }
}


func (ctx *LanternContext) CGFSimpleStats() {
  fmt.Printf("start\n")
  for ii:=0; ii<len(ctx.CGFPathi); ii++ {

    pathi := ctx.CGFPathi[ii]

    count:=0
    for i:=0; i<len(pathi.VecUint64); i++ {
      u := (pathi.VecUint64[i]>>32)

      for j:=uint32(0); j<32; j++ {
        if (u&(1<<j)) > 0 { count++; }
      }

    }

    c_count := len(pathi.VecUint64)*32 - count ; _ = c_count

  }

  fmt.Printf("end\n")

  pprof.StartCPUProfile(g_pprof)
  pprof.StopCPUProfile()
  os.Exit(0)

}

func (ctx *LanternContext) CGFSimpleStats_() {
  n := len(ctx.CGFBytes)

  paths := []int{ 0x247, 0x2c5 }

  fmt.Printf("start\n")
  for ii:=0; ii<len(paths); ii++ {

    path := paths[ii]

    for i:=0; i<n; i++ {
      hdri,_ := cgf.HeaderIntermediateFromBytes(ctx.CGFBytes[i])
      path_bytes := hdri.PathBytes[path]
      pathi,_ := cgf.PathIntermediateFromBytes(path_bytes)

      count:=0
      for i:=0; i<len(pathi.VecUint64); i++ {
        u := (pathi.VecUint64[i]>>32)

        for j:=uint32(0); j<32; j++ {
          if (u&(1<<j)) > 0 { count++; }
        }

      }

      c_count := len(pathi.VecUint64)*32 - count ; _ = c_count

      //fmt.Printf("[%d] path %x canon: %d/%d, (non-cacnon %d/%d)\n", i, path, c_count, len(pathi.VecUint64)*32, count, len(pathi.VecUint64)*32)

    }

  }

  fmt.Printf("end\n")

  pprof.StartCPUProfile(g_pprof)
  pprof.StopCPUProfile()
  os.Exit(0)

}


func handle_status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  if gVerboseFlag { log.Printf("status (%s)", gAPIVersionString) }

  //if ctx.VerboseAPIFlag {
  if true {
    io.WriteString(w,"{\n")
    io.WriteString(w,`  "api-version":"` + gAPIVersionString + `"` + "\n")
    io.WriteString(w,"}\n")
  } else {
    io.WriteString(w,`{"api-version":"` + gAPIVersionString + `"}`)
  }
}


