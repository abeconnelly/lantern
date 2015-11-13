package main

import "time"
import gouuid "github.com/nu7hatch/gouuid"

import "github.com/abeconnelly/sloppyjson"

import "github.com/julienschmidt/httprouter"
import "net/http"
import "log"

import "github.com/abeconnelly/cgf"
import "github.com/abeconnelly/cglf"

const gAPIVersionString = "0.1.0"
const gLanternVersion = "0.1.0"


type APILocusStruct struct {
  ChromosomeName string
  Indexing int
  StartPosition int
  EndPosition int
}

type APIAssembliesStruct struct {
  Name string
  PDH string
  Locus []APILocusStruct
}

type LanternStatStruct struct {
  StartTime time.Time
  Requests int
  Failures int
  AvgRespMs []float64
  AvgRespDt float64
}

type LanternConnInfo struct {
  UUID *gouuid.UUID
}

type LanternTileInfo struct {
  Path int
  Step int
  Variant int
  Span int
  Md5sum string
}

type LanternContext struct {
  VerboseFlag bool
  PrettyAPIFlag bool
  VerboseAPIFlag bool

  Version string
  Config *sloppyjson.SloppyJSON

  CGFBytes [][]byte
  CGFPath [][]cgf.PathIntermediate
  CGFi []cgf.HeaderIntermediate
  CGFPathi []cgf.PathIntermediate

  SGLF cglf.SGLF

  // assembly-pdh, path, step
  //
  Assembly map[string]map[int][]int

  // assembly-pdh, path, chromosome name ('chr' prefix)
  //
  AssemblyChrom map[string]map[int]string

  // path, step, md5sum
  //
  TileInfoMap map[int]map[int]map[string]LanternTileInfo
  TileInfo map[int]map[int][]LanternTileInfo
  TileMap []cgf.TileMapEntry
}

func (ctx *LanternContext) Qux(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  log.Printf("Qux\n")
}
