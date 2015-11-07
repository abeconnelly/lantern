package main

import "io"
import "fmt"
import "net/http"

import "log"

import "github.com/julienschmidt/httprouter"

func (ctx *LanternContext) APIAssemblies(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  if ctx.VerboseFlag {
    log.Printf("APIAssemblies\n")
  }

  assembly_map := make(map[string]string)

  tagset_sj := ctx.Config.O["tagset"].O
  for tagset_pdh := range tagset_sj {
    for assembly_pdh := range tagset_sj[tagset_pdh].O["assembly"].O {
      assembly_name := tagset_sj[tagset_pdh].O["assembly"].O[assembly_pdh].O["name"].S
      assembly_map[assembly_pdh] = assembly_name
    }
  }


  count:=0
  io.WriteString(w,"[")
  for pdh := range assembly_map {
    if count>0 { io.WriteString(w,",") }
    io.WriteString(w,"{")
    io.WriteString(w, fmt.Sprintf(`"assembly-name":"%s","assembly-pdh":"%s"`, assembly_map[pdh], pdh))
    io.WriteString(w,"}")
    count++
  }
  io.WriteString(w,"]")

}

func (ctx *LanternContext) APIAssembliesId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  assembly_pdh := params.ByName("id")
  if len(assembly_pdh) == 0 {
    if ctx.VerboseFlag { log.Printf("APIAssembliesId\n") }

    send_error_bad_request(w, "invalid assembly id")
    return
  }


  if ctx.VerboseFlag {
    log.Printf("APIAssembliesId %v\n", assembly_pdh)
  }

  assembly_map := make(map[string]string)
  tagset_sj := ctx.Config.O["tagset"].O
  for tagset_pdh := range tagset_sj {
    for assembly_pdh := range tagset_sj[tagset_pdh].O["assembly"].O {
      assembly_name := tagset_sj[tagset_pdh].O["assembly"].O[assembly_pdh].O["name"].S
      assembly_map[assembly_pdh] = assembly_name
    }
  }

  count:=0
  io.WriteString(w,"[")
  for assembly_pdh := range ctx.Assembly {
    assembly_name := assembly_map[assembly_pdh]
    for path := range ctx.Assembly[assembly_pdh] {

      chrom := ctx.AssemblyChrom[assembly_pdh][path]

      beg_pos := 0
      if (path>0) && (ctx.AssemblyChrom[assembly_pdh][path-1]==chrom) {
        n := len(ctx.Assembly[assembly_pdh][path-1])
        beg_pos = ctx.Assembly[assembly_pdh][path-1][n-1]
      }

      if (len(chrom)>3) && (chrom[0:3] == "chr") { chrom = chrom[3:] }
      for step:=0; step<len(ctx.Assembly[assembly_pdh][path]); step++ {

        end_pos := ctx.Assembly[assembly_pdh][path][step]

        if count>0 { io.WriteString(w,",") }
        io.WriteString(w,"{")
        io.WriteString(w, fmt.Sprintf(`"assembly-name":"%s","assembly-pdh":"%s"`, assembly_name, assembly_pdh))
        io.WriteString(w, fmt.Sprintf(`,"chromosome-name":"%s","indexing":0`, chrom))
        io.WriteString(w, fmt.Sprintf(`,"start-position":%d,"end-position":%d`, beg_pos, end_pos))
        io.WriteString(w,"}")

        if ctx.PrettyAPIFlag {
          io.WriteString(w,"\n")
        }
        count++

        beg_pos = end_pos
      }

    }
  }
  io.WriteString(w,"]")


}
