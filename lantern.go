package main

import "fmt"
import "os"
import "io"
import "runtime"
import "runtime/pprof"

import "github.com/abeconnelly/autoio"
import "github.com/codegangsta/cli"

//import "net"
import "net/http"

import "syscall"
import "os/signal"
import "time"

import "log"

import "io/ioutil"
import "github.com/abeconnelly/sloppyjson"

import "github.com/julienschmidt/httprouter"

import "strings"
import "strconv"

var VERSION_STR string = "0.2.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "lantern.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "lantern.mprof"

var gPortStr string = ":8080"

var gLanternStat LanternStatStruct

var gConfig *sloppyjson.SloppyJSON

func _skip_space(b string) string {
  n := len(b)
  var x byte

  for i:=0; i<n; i++ {
    x = b[i]
    if x=='\t' || x==' ' || x=='\n' || x=='\r' { continue }
    return b[i:]
  }
  return ""
}

func _load_json_config(ctx *LanternContext, config_fn string) error {
  var e error

  raw_str,err := ioutil.ReadFile(config_fn)
  if err!=nil { return err }

  //gConfig,e = sloppyjson.Loads(string(raw_str))
  ctx.Config,e = sloppyjson.Loads(string(raw_str))
  if e!=nil { return e }

  log.Printf("config loaded\n")

  //gConfig.Printr(0,2)

  return nil
}

func index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
  io.WriteString(w, `{"message":"lantern API server"}`)
}

func load_Assembly(ctx *LanternContext, tagset_pdh, assembly_pdh string) error {
  assembly_fn := ctx.Config.O["tagset"].O[tagset_pdh].O["assembly"].O[assembly_pdh].O["gz"].S
  fp,e := autoio.OpenReadScanner(assembly_fn)
  if e!=nil { return e }
  defer fp.Close()

  if ctx.Assembly == nil { ctx.Assembly = make(map[string]map[int][]int) }
  ctx.Assembly[assembly_pdh] = make(map[int][]int)

  if ctx.AssemblyChrom == nil { ctx.AssemblyChrom = make(map[string]map[int]string) }
  ctx.AssemblyChrom[assembly_pdh] = make(map[int]string)

  path := 0
  for fp.ReadScan() {
    l := fp.ReadText()

    if len(l) == 0 { continue }
    if l[0] == '\n' { continue }

    if l[0] == '>' {
      parts  := strings.Split(l[1:], ":")
      name   := parts[0] ; _ = name
      chrom  := parts[1] ; _ = chrom
      path_s := parts[2]

      _path,e := strconv.ParseInt(path_s, 16, 64)
      if e!=nil { return e }
      path = int(_path)

      ctx.Assembly[assembly_pdh][path] = make([]int, 0, 1024)
      ctx.AssemblyChrom[assembly_pdh][path] = chrom
      continue
    }

    _step,e := strconv.ParseInt(l[0:4], 16, 64)
    if e!=nil {return e}
    step := int(_step) ; _ = step

    z := _skip_space(l[5:])
    _ref_pos,e := strconv.ParseInt(z, 10, 64)
    if e!=nil {return e}
    ref_pos := int(_ref_pos)

    ctx.Assembly[assembly_pdh][path] = append(ctx.Assembly[assembly_pdh][path], ref_pos)
  }

  return nil
}

func _main( c *cli.Context ) {
  gLanternStat = LanternStatStruct{}
  gLanternStat.StartTime = time.Now()



  ctx := &LanternContext{}

  //DEBUG
  ctx.VerboseFlag = true
  ctx.PrettyAPIFlag = true

  e := _load_json_config(ctx, c.String("config"))
  if e!=nil { log.Fatal(e) }

  tagset_pdh := "dad7041d432965cd07a4ad8e0aad1b6e"
  assembly_pdh := "aa39590be6f1812f0a792dd4c86678e8+1348"

  if ctx.VerboseFlag { log.Printf("loading assembly %v (tagset %v)\n", tagset_pdh, assembly_pdh) }

  e = load_Assembly(ctx, tagset_pdh, assembly_pdh)
  if e!=nil { log.Fatal(e) }

  if ctx.VerboseFlag { log.Printf("assembly loaded\n") }


  /*
  //DEBUG
  fmt.Printf("testing...\n")
  for pdh := range ctx.Assembly {
    for path := range ctx.Assembly[pdh] {
      for i:=0; i<len(ctx.Assembly[pdh][path]);  i++ {
        fmt.Printf("%v %v %v\n", pdh, path, ctx.Assembly[pdh][path][i])
      }
    }

    os.Exit(0)
  }

  for pdh := range ctx.AssemblyChrom {
    for path := range ctx.AssemblyChrom[pdh] {
      fmt.Printf("%v %v %v\n", pdh, path, ctx.AssemblyChrom[pdh][path])
    }
  }

  fmt.Printf("cp\n")
  */


  /*
  if c.String("input") == "" {
    fmt.Fprintf( os.Stderr, "Input required, exiting\n" )
    cli.ShowAppHelp( c )
    os.Exit(1)
  }

  ain,err := autoio.OpenScanner( c.String("input") ) ; _ = ain
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer ain.Close()
  */


  aout,err := autoio.CreateWriter( c.String("output") ) ; _ = aout
  if err!=nil {
    fmt.Fprintf(os.Stderr, "%v", err)
    os.Exit(1)
  }
  defer func() { aout.Flush() ; aout.Close() }()

  if c.Bool( "pprof" ) {
    gProfileFlag = true
    gProfileFile = c.String("pprof-file")
  }

  if c.Bool( "mprof" ) {
    gMemProfileFlag = true
    gMemProfileFile = c.String("mprof-file")
  }

  gVerboseFlag = c.Bool("Verbose")

  if c.Int("max-procs") > 0 {
    runtime.GOMAXPROCS( c.Int("max-procs") )
  }

  if gProfileFlag {
    prof_f,err := os.Create( gProfileFile )
    if err != nil {
      fmt.Fprintf( os.Stderr, "Could not open profile file %s: %v\n", gProfileFile, err )
      os.Exit(2)
    }

    pprof.StartCPUProfile( prof_f )
    defer pprof.StopCPUProfile()
  }

  // Start server
  //


  //listener,err := net.Listen("tcp", gPortStr)
  //if err!=nil { log.Fatal(err) }

  term := make(chan os.Signal,1)
  go func( sig <-chan os.Signal) {
    s := <-sig
    if gVerboseFlag {
      fmt.Printf("caught signal: %v\n", s)
      //listener.Close()
    }
  }(term)
  signal.Notify(term, syscall.SIGTERM)
  signal.Notify(term, syscall.SIGHUP)

  // Set up routing
  //
  router := httprouter.New()

  router.POST("/", handle_json_req)
  router.GET("/", index)
  router.GET("/status", handle_status)

  //router.GET("/qux", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) { ctx.Qux(w,r, params); } )
  router.GET("/qux", ctx.Qux)

  router.GET("/assemblies", ctx.APIAssemblies)
  router.GET("/assemblies/:id", ctx.APIAssembliesId)

  //router.GET("/callsets", handle_callsets)

  router.GET("/tile-library/tag-sets", handle_tile_library_tag_sets)
  router.GET("/tile-library/tag-sets/:tagset_id", handle_tile_library_tag_sets_id)
  router.GET("/tile-library/tag-sets/:tagset_id/paths", handle_tile_library_tag_sets_id_paths)
  router.GET("/tile-library/tag-sets/:tagset_id/paths/:path_id", handle_tile_library_tag_sets_id_paths_id)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-positions", handle_tile_library_tag_sets_id_tile_positions)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-positions/:tilepos_id", handle_tile_library_tag_sets_id_tile_positions_id)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-positions/:tilepos_id/locus", handle_tile_library_tag_sets_id_tile_positions_id_locus)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-variants", handle_tile_library_tag_sets_id_tile_variants)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-variants/:tilevariant_id", handle_tile_library_tag_sets_id_tile_variants_id)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-variants/:tilevariant_id/locus", handle_tile_library_tag_sets_id_tile_variants_id_locus)
  //router.GET("/tile-library/tag-sets/:tagset_id/tile-variants/:tilevariant_id/subsequence", handle_tile_library_tag_sets_id_tile_variants_id_subsequence)
  router.GET("/tile-library/tag-sets/:tagset_id/tile-variants/:tilevariant_id/annotations", handle_tile_library_tag_sets_id_tile_variants_id_annotations)

  /*
  http.HandleFunc("/", handle_json_req)
  http.HandleFunc("/status", handle_status)
  http.HandleFunc("/assemblies/", handle_assemblies_id)
  http.HandleFunc("/assemblies", handle_assemblies)
  */


  if gVerboseFlag {
    fmt.Printf("listening: %v\n", gPortStr)
  }

  http.ListenAndServe(gPortStr, router)

  if gVerboseFlag {
    fmt.Printf("shutting down\n")
  }

}

func main() {

  app := cli.NewApp()
  app.Name  = "lantern"
  app.Usage = "lantern -C config"
  app.Version = VERSION_STR
  app.Author = "Curoverse, Inc."
  app.Email = "info@curoverse.com"
  app.Action = func( c *cli.Context ) { _main(c) }

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name: "input, i",
      Usage: "INPUT",
    },

    cli.StringFlag{
      Name: "output, o",
      Value: "-",
      Usage: "OUTPUT",
    },

    cli.StringFlag{
      Name: "config, C",
      Value: "~/.config/lantern/config.json",
      Usage: "Config file (default to \"$HOME/.config/lantern/config.json\")",
    },

    cli.IntFlag{
      Name: "max-procs, N",
      Value: -1,
      Usage: "MAXPROCS",
    },

    cli.BoolFlag{
      Name: "Verbose, V",
      Usage: "Verbose flag",
    },

    cli.BoolFlag{
      Name: "pprof",
      Usage: "Profile usage",
    },

    /*
    cli.StringFlag{
      Name: "config, C",
      Value: "config.json",
      Usage: "config file (JSON)",
    },
    */

    cli.StringFlag{
      Name: "pprof-file",
      Value: gProfileFile,
      Usage: "Profile File",
    },

    cli.BoolFlag{
      Name: "mprof",
      Usage: "Profile memory usage",
    },

    cli.StringFlag{
      Name: "mprof-file",
      Value: gMemProfileFile,
      Usage: "Profile Memory File",
    },

  }

  app.Run( os.Args )

  if gMemProfileFlag {
    fmem,err := os.Create( gMemProfileFile )
    if err!=nil { panic(fmem) }
    pprof.WriteHeapProfile(fmem)
    fmem.Close()
  }

}
