package main

//import "fmt"
import "io/ioutil"

import "github.com/abeconnelly/cgf"


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

func (ctx *LanternContext) LoadCGFBytes(cgf_dir string) error {

  files,e := ioutil.ReadDir(cgf_dir)
  if e!=nil { return e }

  ctx.CGFBytes = make([][]byte, 0, 1024)

  for _,f := range files {
    cgf_bytes,e := ioutil.ReadFile( cgf_dir + "/" + f.Name())
    if e!=nil { return e }
    ctx.CGFBytes = append(ctx.CGFBytes, cgf_bytes)
  }

  return nil
}
