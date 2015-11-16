package main

import "log"
import "io/ioutil"
import "github.com/abeconnelly/cgf"
//import "github.com/abeconnelly/cglf"

import "strings"

func (ctx *LanternContext) LoadCGFIntermediate() {
  for i:=0; i<len(ctx.CGFBytes); i++ {
    hdri,_ := cgf.HeaderIntermediateFromBytes(ctx.CGFBytes[i])
    ctx.CGFi = append(ctx.CGFi,hdri)

    m0,_ := cgf.PathIntermediateFromBytes(hdri.PathBytes[0x247])
    ctx.CGFPathi = append(ctx.CGFPathi, m0)

    m1,_ := cgf.PathIntermediateFromBytes(hdri.PathBytes[0x2c5])
    ctx.CGFPathi = append(ctx.CGFPathi, m1)
  }
}

func (ctx *LanternContext) CGFNameMap(fn string) string {
  if fn[0]=='h' && fn[1]=='u' {
    z := strings.Split(fn, ".")
    return "hupgp-" + z[0]
  }

  if strings.Contains(fn, "cg_data") {
    z := strings.Split(fn, ".")
    return "okg-" + z[0]
  }

  return fn
}

func (ctx *LanternContext) LoadCGFBytesConfig() error {
  ctx.CGFBytes = make([][]byte, 0, 1024)
  if ctx.CGFIndexMap == nil {
    ctx.CGFIndexMap = make(map[string]int)
  }

  for idx:=0; idx<len(ctx.Config.O["cgf"].L); idx++ {
    cgf_name := ctx.Config.O["cgf"].L[idx].O["name"].S
    cgf_locator := ctx.Config.O["cgf"].L[idx].O["locator"].S

    ctx.CGFIndexMap[cgf_name] = len(ctx.CGFBytes)
    cgf_bytes,e := ioutil.ReadFile(cgf_locator)
    if e!=nil { return e }
    ctx.CGFBytes = append(ctx.CGFBytes, cgf_bytes)
  }

  return nil
}

func (ctx *LanternContext) LoadCGFBytesDir(cgf_dir string) error {

  if ctx.CGFIndexMap == nil {
    ctx.CGFIndexMap = make(map[string]int)
  }

  files,e := ioutil.ReadDir(cgf_dir)
  if e!=nil { return e }

  ctx.CGFBytes = make([][]byte, 0, 1024)

  for _,f := range files {

    ctx.CGFIndexMap[ ctx.CGFNameMap(f.Name()) ] = len(ctx.CGFBytes)

    cgf_bytes,e := ioutil.ReadFile(cgf_dir + "/" + f.Name())
    if e!=nil { return e }
    ctx.CGFBytes = append(ctx.CGFBytes, cgf_bytes)
  }

  return nil
}

//====================================
//====================================
//====================================
//====================================
//====================================

func (ctx *LanternContext) LoadSGLF(path int, sglf_fn string) error {
  var e error

  /*
  if ctx.SGLF==nil {
    ctx.SGLF,e = cglf.LoadGenomeLibraryCSV(sglf_fn)
  } else {
    e = ctx.SGLF.AddGenomeLibraryCSV(sglf_fn)
  }
  */

  e = ctx.SGLF.AddGenomeLibraryCSV(sglf_fn)

  if ctx.VerboseFlag {
    log.Printf("Loaded SGLF 0x%x %v\n", path, e)
  }

  return e
}

func (ctx *LanternContext) PreCalc() {

  n_path := len(ctx.CGFi[0].StepPerPath)
  step_count := make([][]int, n_path)
  for i:=0; i<n_path; i++ {
    n_step := ctx.CGFi[0].StepPerPath[i]
    if n_step > 0 { step_count[i] = make([]int, n_step) }
  }

  count:=0
  tot_count:=0

  // Allocate PathIntermediate structures for easy processing.
  // These structures are probably too verbose to really be
  // in memory but they're convenient.  We'll have to
  // figure out how to access the in memory data structures
  // easily for fast queries later.
  //
  ctx.CGFPath = make([][]cgf.PathIntermediate, len(ctx.CGFi))
  for sample:=0; sample<len(ctx.CGFi); sample++ {

    n_path := len(ctx.CGFi[sample].StepPerPath)
    ctx.CGFPath[sample] = make([]cgf.PathIntermediate, n_path)
    for path:=0; path<n_path; path++ {

      if ctx.CGFi[sample].StepPerPath[path]==0 { continue }
      //n_step := ctx.CGFi[sample].StepPerPath[path]

      ctx.CGFPath[sample][path],_ = cgf.PathIntermediateFromBytes(ctx.CGFi[sample].PathBytes[path])
    }

  }



  // Do a simple 'non-canonical' (rough) count.
  // Rough because it counts some steps at the end of
  // the vector.
  //
  for sample:=0; sample<len(ctx.CGFPath); sample++ {

    n_path := len(ctx.CGFPath[sample])
    for path:=0; path<n_path; path++ {

      n_step := 32*len(ctx.CGFPath[sample][path].VecUint64)

      for step:=0; step<n_step; step++ {
        b := step/32
        c := step%32
        u := ctx.CGFPath[sample][path].VecUint64[b]>>32
        if (u & (1<<uint(c))) > 0 { count++ }
        tot_count++
      }

    }
  }

  ctx.TileMap = ctx.CGFi[0].TileMap

  // ...
  //
  /*
  n_sample := len(ctx.CGFPath)
  n_path := len(ctx.CGFPath[0])

  tcount := make([]int, n_sample)

  for path:=0; path<n_path; path++ {

    n_step := 32*len(ctx.CGFPath[sample][path].VecUint64)

    for step:=0; step<n_step; step++ {

      for i:=0; i<n_sample; i++ {
        tcount[i] = 0;
      }

      for sample:=0; sample<n_sample; sample++ {
        b := step/32
        c := step%32
        u := ctx.CGFPath[sample][path].VecUint64[b]>>32
        if (u & (1<<uint(c))) > 0 {
          count++
        }
        tot_count++
      }

    }
  }
  */



  log.Printf(">>> non-canon? %d\n", count)
  log.Printf(">>> tot_count %d\n", tot_count)

}
