package main

import "io"
import "log"
import "net/http"
import "fmt"
import "strings"
import "strconv"

import "github.com/julienschmidt/httprouter"

type TagSetInfoStruct struct {
  PDH string
  Version int
}

type TmpStruct struct {
  Name string
  Id int
  Path int
  StepPerPath int
}

var gTagSets []TagSetInfoStruct
var gTagSetIdx map[string]int

var gCGF []TmpStruct



//====

/*
func api_tile_library_init() {
  gTagSets = make([]TagSetInfoStruct, 0, 3)
  gTagSetIdx = make(map[string]int)

  // TEMPORARY
  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"


  gTagSets = append(gTagSets, TagSetInfoStruct{ PDH: pdh, Version:0 } )
  gTagSetIdx[pdh] = 0


  gCGF = make([]TmpStruct, 0, 1024)
  for path:=0; path<=0x2c5; path++ {
    gCGF = append(gCGF, TmpStruct{Name:fmt.Sprintf("%04x", path), Path:path, StepPerPath: 101+path})
  }
}
*/

//====

func (ctx *LanternContext) APITileLibraryTagSets(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSets\n")
  }

  count := 0
  io.WriteString(w,"[")
  for tagset_pdh := range ctx.Config.O["tagset"].O {

    if count>0 { io.WriteString(w,",") }
    io.WriteString(w, fmt.Sprintf(`"%s"`, tagset_pdh))
    count++

  }
  io.WriteString(w,"]")

}

func (ctx *LanternContext) APITileLibraryTagSetsId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsId %v\n", tagset_id)
  }

  tagset_ver := int(ctx.Config.O["tagset"].O[tagset_id].O["id"].P+0.5)

  io.WriteString(w,"{")
  io.WriteString(w,fmt.Sprintf(`"tag-set-identifier":"%s","tag-set-integer":%02x`,tagset_id, tagset_ver))
  io.WriteString(w,"}")

}

func (ctx *LanternContext) APITileLibraryTagSetsIdPaths(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsIdPaths %v\n", tagset_id)
  }

  step_per_path_sjl := ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L


  io.WriteString(w,"[")
  for i:=0; i<len(step_per_path_sjl); i++ {
    if i>0 { io.WriteString(w,`,`) }
    io.WriteString(w,fmt.Sprintf(`%d`, i))
  }
  io.WriteString(w,"]")

}


func (ctx *LanternContext) APITileLibraryTagSetsIdPathsId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  path_id := param.ByName("path_id")
  _path,e := strconv.ParseInt(path_id, 16, 64)
  if e!=nil {
    send_error_bad_request(w, "invalid path id")
    return
  }
  ipath:=int(_path)

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsIdPathsId %v %v\n", tagset_id, path_id)
  }

  if ipath<0 || ipath>=len(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L) {
    send_error_bad_request(w, "invalid path id")
    return
  }

  n_path := int(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L[ipath].P+0.5)

  io.WriteString(w,"{")
  io.WriteString(w,fmt.Sprintf(`"path":%d,"num-positions":%d`, ipath,n_path))
  io.WriteString(w,"}")

}



func (ctx *LanternContext) APITileLibraryTagSetsIdTilePositions(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsIdTilePositions %v\n", tagset_id)
  }

  tagset_name := int(ctx.Config.O["tagset"].O[tagset_id].O["id"].P+0.5)

  count:=0
  io.WriteString(w,`[`)
  for i:=0; i<len(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L); i++ {

    n_path := int(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L[i].P+0.5)

    for j:=0; j<n_path; j++ {
      if count>0 { io.WriteString(w, `,`) }
      io.WriteString(w, fmt.Sprintf(`"%02x.%04x.%04x"`, tagset_name, i, j))
      count++
    }

  }
  io.WriteString(w,`]`)

}

func (ctx *LanternContext) APITileLibraryTagSetsIdTilePositionsId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  tilepos := param.ByName("tilepos_id")

  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  if len(tilepos) == 0 {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  tile_parts := strings.Split(tilepos, ".")
  if len(tile_parts)!=3 {
    send_error_bad_request(w, "invalid tilepos")
    return
  }

  _req_tagset_name,e := strconv.ParseInt(tile_parts[0], 16, 64)
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

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsIdTilePositionsId %v\n", tagset_id)
  }


  req_tagset_name := int(_req_tagset_name) ; _ = req_tagset_name
  req_path := int(_req_path)
  req_step := int(_req_step)

  if req_path<0 || req_path>=len(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L) {
    send_error_bad_request(w, "path out of range")
    return
  }

  if req_step<0 || req_step>=int(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L[req_path].P+0.5) {
    send_error_bad_request(w, "step out of range")
    return
  }

  tagset_name := int(ctx.Config.O["tagset"].O[tagset_id].O["id"].P+0.5)
  tilepos_str := fmt.Sprintf("%02x.%04x.%04x", tagset_name, req_path, req_step)

  //TODO: get this information
  // this is a placeholder for now
  n_tile_var:=-1
  n_hiq_tile_var:=-1
  n_genomes:=-1

  io.WriteString(w,`{`)
  io.WriteString(w, fmt.Sprintf(`"tile-position":"%s"`, tilepos_str))
  io.WriteString(w, fmt.Sprintf(`,"total-tile-variants":%d`, n_tile_var))
  io.WriteString(w, fmt.Sprintf(`,"well-sequenced-tile-variants":%d`, n_hiq_tile_var))
  io.WriteString(w, fmt.Sprintf(`,"num-genomes":%d`, n_genomes))
  io.WriteString(w,`}`)

}

func (ctx *LanternContext) APITileLibraryTagSetsIdTilePositionsIdLocus(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

  tagset_id := param.ByName("tagset_id")
  tilepos := param.ByName("tilepos_id")

  req_assembly_name := param.ByName("assembly-name")
  req_assembly_pdh := param.ByName("assembly-pdh")

  if _,ok := ctx.Config.O["tagset"].O[tagset_id] ; !ok {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  if ctx.VerboseFlag {
    log.Printf("APITileLibraryTagSetsIdTilePositionsIdLocus %v %v (assembly name %v, pdh %v)\n", tagset_id, tilepos, req_assembly_name, req_assembly_pdh)
  }

  if len(tilepos) == 0 {
    send_error_bad_request(w, "invalid tagset id")
    return
  }

  tile_parts := strings.Split(tilepos, ".")
  if len(tile_parts)!=3 {
    send_error_bad_request(w, "invalid tilepos")
    return
  }

  _req_tagset_name,e := strconv.ParseInt(tile_parts[0], 16, 64)
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


  req_tagset_name := int(_req_tagset_name) ; _ = req_tagset_name
  req_path := int(_req_path)
  req_step := int(_req_step)

  if req_path<0 || req_path>=len(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L) {
    send_error_bad_request(w, "path out of range")
    return
  }

  if req_step<0 || req_step>=int(ctx.Config.O["tagset"].O[tagset_id].O["step_per_path"].L[req_path].P+0.5) {
    send_error_bad_request(w, "step out of range")
    return
  }

  tagset_name := int(ctx.Config.O["tagset"].O[tagset_id].O["id"].P+0.5)
  tilepos_str := fmt.Sprintf("%02x.%04x.%04x", tagset_name, req_path, req_step) ; _ = tilepos_str

  //TODO: fill in values
  // these are only placeholders now
  chrom := "Unk"
  start_position := -1
  end_position := -1


  io.WriteString(w,`[`)
  for i:=0; i<1; i++ {
    io.WriteString(w,`{`)
    io.WriteString(w, fmt.Sprintf(`"assembly-name":"%s"`, req_assembly_name))
    io.WriteString(w, fmt.Sprintf(`,"assembly-pdh":"%s"`, req_assembly_pdh))
    io.WriteString(w, fmt.Sprintf(`,"chromosome-name":"%s"`, chrom))
    io.WriteString(w, fmt.Sprintf(`,"indexing":"%d"`, 0))
    io.WriteString(w, fmt.Sprintf(`,"start-position":"%d"`, start_position))
    io.WriteString(w, fmt.Sprintf(`,"end-position":"%d"`, end_position))
    io.WriteString(w,`}`)
  }
  io.WriteString(w,`]`)

}


//================================
//================================
//================================
//================================
//================================

func handle_tile_library_tag_sets(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
  if gVerboseFlag { log.Printf("tile-library/tag-sets") }

  count := 0
  io.WriteString(w,"[")
  for pdh := range gConfig.O["tagset"].O {

    if count>0 { io.WriteString(w,",") }
    io.WriteString(w, fmt.Sprintf(`"%s"`, gConfig.O["tagset"].O[pdh].O["pdh"].S))
    count++

  }
  io.WriteString(w,"]")

}

//====

func handle_tile_library_tag_sets_id(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  pdh := params.ByName("tagset_id")

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v", pdh) }

  var ok bool

  if _,ok = gConfig.O["tagset"].O[pdh] ; !ok {
    io.WriteString(w,`{}`)
    return
  }

  io.WriteString(w, fmt.Sprintf(`{"tag-set-identifier":"%s","tag-set-integer":%d}`,
    gConfig.O["tagset"].O[pdh].O["pdh"].S,
    int(gConfig.O["tagset"].O[pdh].O["id"].P + 0.5)))
}

//====

func handle_tile_library_tag_sets_id_paths(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  pdh := params.ByName("tagset_id")

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v", pdh) }

  //var idx int
  var ok bool

  if _,ok = gConfig.O["tagset"].O[pdh] ; !ok {
    io.WriteString(w,`[]`)
    return
  }

  io.WriteString(w,`[`)
  n := len(gConfig.O["tagset"].O[pdh].O["step_per_path"].L)
  for i:=0; i<n; i++ {
    if i>0 { io.WriteString(w,",") }
    io.WriteString(w,fmt.Sprintf(`%d`, i))
  }
  io.WriteString(w,`]`)

}

//====

func handle_tile_library_tag_sets_id_paths_id(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  pdh := params.ByName("tagset_id")
  path_str := params.ByName("path_id")

  _path,e := strconv.ParseInt(path_str, 16, 64)
  if e!=nil {
    io.WriteString(w,`{}`)
    return
  }
  path:=int(_path)

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/paths-id/%x", pdh, path) }

  var ok bool

  if _,ok = gConfig.O["tagset"].O[pdh] ; !ok {
    io.WriteString(w,`[]`)
    return
  }

  step_per_path := gConfig.O["tagset"].O[pdh].O["step_per_path"].L

  if (path<0) || (path>=len(step_per_path)) {
    io.WriteString(w,`{}`)
    return
  }

  io.WriteString(w, fmt.Sprintf(`{"path":%d,"num-positions":%d}`,
    path, int(step_per_path[path].P+0.5) ))
}

//====

func handle_tile_library_tag_sets_id_tile_positions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  pdh := params.ByName("tagset_id")

  var ok bool
  if _,ok = gConfig.O["tagset"].O[pdh] ; !ok {
    io.WriteString(w,`[]`)
    return
  }

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-positions", pdh) }

  tagset_id := int(gConfig.O["tagset"].O[pdh].O["id"].P+0.5)

  count:=0

  io.WriteString(w,`[`)
  step_per_path := gConfig.O["tagset"].O[pdh].O["step_per_path"].L
  n := len(step_per_path)
  for i:=0; i<n; i++ {

    spp := int(step_per_path[i].P+0.5)

    for j:=0; j<spp; j++ {
      if count>0 { io.WriteString(w,",") }
      io.WriteString(w,fmt.Sprintf(`"%02x.%04x.%04x"`, tagset_id, i,j))
      count++
    }
  }
  io.WriteString(w,`]`)

}

//===

func handle_tile_library_tag_sets_id_tile_positions_id(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  pdh := params.ByName("tagset_id")
  tilepos := params.ByName("tilepos_id")

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-positions/%v", pdh, tilepos) }

  var ok bool
  if _,ok = gConfig.O["tagset"].O[pdh] ; !ok {
    io.WriteString(w,`{}`)
    return
  }

  parts := strings.Split(tilepos, `.`)
  if len(parts)!=3 {
    io.WriteString(w,`{}`)
    return
  }

  _tagset,e := strconv.ParseInt(parts[0], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  _path,e := strconv.ParseInt(parts[1], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  _step,e := strconv.ParseInt(parts[2], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  tagset := int(_tagset) ; _ = tagset
  path := int(_path) ; _ = path
  step := int(_step) ; _ = step

  step_per_path := gConfig.O["tagset"].O[pdh].O["step_per_path"].L

  if (path<0) || (path>=len(step_per_path)) { io.WriteString(w,`{}`) ; return ; }

  //DEBUG
  n_tile_variant := 23
  n_tile_well_sequenced := 11
  n_genome := 3
  //DEBUG

  r_tilepos := fmt.Sprintf("%02x.%04x.%04x", tagset, path, step)

  io.WriteString(w,`{`)
  io.WriteString(w,
    fmt.Sprintf(`"tile-position":"%s","total-tile-variants":%d,"well-sequenced-tile-variants":%d,"num-genomes":%d`,
      r_tilepos,
      n_tile_variant,
      n_tile_well_sequenced,
      n_genome))
  io.WriteString(w,`}`)

}

//===

func handle_tile_library_tag_sets_id_tile_positions_id_locus(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"
  tilepos := "00.2c5.001c"
  //DEBUG


  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-positions/%v/locus", pdh, tilepos) }

  parts := strings.Split(tilepos, `.`)
  if len(parts)!=3 {
    io.WriteString(w,`{}`)
    return
  }

  _tagset,e := strconv.ParseInt(parts[0], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  _path,e := strconv.ParseInt(parts[0], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  _step,e := strconv.ParseInt(parts[0], 16, 64)
  if e!=nil { io.WriteString(w,`{}`) ; return ; }

  tagset := int(_tagset) ; _ = tagset
  path := int(_path) ; _ = path
  step := int(_step) ; _ = step

  if (path<0) || (path>=len(gCGF)) { io.WriteString(w,`{}`) ; return ; }

  //DEBUG
  test_str := `[
    {
      "assembly-name": "hg19",
      "assembly-pdh": "dad94936d4144f5e0a289244d8be93e9+5735",
      "chromosome-name": "13",
      "indexing": 0,
      "start-position": 32199976,
      "end-position": 32200225
     }
  ]`

  io.WriteString(w, test_str)

}

//===

func handle_tile_library_tag_sets_id_tile_variants(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"
  //DEBUG


  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-variants", pdh) }


  //DEBUG
  test_str := `[
    "00.247.0000.830003ac103a97d8f7992e09594ac68e",
    "00.247.0000.455577ff6b0d38188477ee2bfb2f0ea8",
    "00.247.1bfb.c95325c08a449529143776e18561db71",
    "00.2c5.0000.1948117b4a56e4ad73d36dce185110fd",
    "00.2c5.30ae.bc952f709d7419f7e103daa2b7e469a9"
  ]`

  io.WriteString(w,test_str)
}

//===

func handle_tile_library_tag_sets_id_tile_variants_id(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"
  var tile_id string = "00.2c5.30ae.bc952f709d7419f7e103daa2b7e469a9"
  //DEBUG

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-variants/%v", pdh, tile_id) }

  //DEBUG
  test_str := `{
    "tile-variant":"00.2c5.30ae.bc952f709d7419f7e103daa2b7e469a9",
    "tag-length": 24,
    "start-tag": "gccaaggagttttaaaactactga",
    "end-tag": "",
    "is-start-of-path": false,
    "is-end-of-path": true,
    "sequence" : "gccaaggagttttaaaactactgatgcccacctcccacacccaaaagtctgattaattgatctagggtatggcctgagcttcaagagtttttaaagcatccaggtgattacaatgtgtagtgaagtttgagagccactgcacaacattaataattgttgggagaaagactgtggctttagctagggagagctgtccagaagatctgaatgtcaggagagagactagtgagagatttggaaaccatcaacatattgatggtaactgaagccacagaagtggacaacactgccttaggagaagatgccaaataacaagagagtagatacaaagacattttgacataacaaagtatggttacagaaatattttcaggtggaaaggaagttgaaggga",
    "md5sum": "bc952f709d7419f7e103daa2b7e469a9",
    "length": 394,
    "number-of-positions-spanned": 1,
    "population-frequency": 0.5,
    "population-count": 150,
    "population-total": 300
  }`

  io.WriteString(w,test_str)

}

//===

func handle_tile_library_tag_sets_id_tile_variants_id_locus(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"
  var tile_id string = "00.2c5.30ae.bc952f709d7419f7e103daa2b7e469a9"
  //DEBUG

  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-variants/%v/locus", pdh, tile_id) }

  //DEBUG
  test_str := `[
  {
    "assembly-name": "hg19",
    "assembly-pdh": "dad94936d4144f5e0a289244d8be93e9+5735",
    "chromosome-name": "13",
    "indexing": 0,
    "start-position": 32199976,
    "end-position": 32200225
  }
]`
  //DEBUG


  io.WriteString(w,test_str)

}

//===

func handle_tile_library_tag_sets_id_tile_variants_id_annotations(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  //DEBUG
  var pdh string = "1c20dd595e9fd3d8eefb281e314709ec+67"
  var tile_id string = "00.2c5.30ae.bc952f709d7419f7e103daa2b7e469a9"
  //DEBUG


  if gVerboseFlag { log.Printf("tile-library/tag-sets/%v/tile-variants/%v/annotations", pdh, tile_id) }


  //DEBUG
  test_str := `["annotation0","annotation3"]`
  //DEBUG


  io.WriteString(w,test_str)

}
