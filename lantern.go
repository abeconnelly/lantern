package main

import "fmt"
import "os"
import "runtime"
import "runtime/pprof"

import "github.com/abeconnelly/autoio"
import "github.com/codegangsta/cli"

import "net"
import "net/http"

import "syscall"
import "os/signal"
import "time"

import "log"

var VERSION_STR string = "0.2.0"
var gVerboseFlag bool

var gProfileFlag bool
var gProfileFile string = "lantern.pprof"

var gMemProfileFlag bool
var gMemProfileFile string = "lantern.mprof"

var gPortStr string = ":8080"

var gLanternStat LanternStatStruct

func _main( c *cli.Context ) {
  gLanternStat = LanternStatStruct{}
  gLanternStat.StartTime = time.Now()

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



  listener,err := net.Listen("tcp", gPortStr)
  if err!=nil { log.Fatal(err) }

  term := make(chan os.Signal,1)
  go func( sig <-chan os.Signal) {
    s := <-sig
    if gVerboseFlag {
      fmt.Printf("caught signal: %v\n", s)
      listener.Close()
    }
  }(term)
  signal.Notify(term, syscall.SIGTERM)
  signal.Notify(term, syscall.SIGHUP)

  // Set up routing
  //
  http.HandleFunc("/", handle_json_req)

  http.HandleFunc("/status", handle_status)

  http.HandleFunc("/assemblies/", handle_assemblies_id)
  http.HandleFunc("/assemblies", handle_assemblies)

  srv := &http.Server{ Addr: gPortStr }

  if gVerboseFlag {
    fmt.Printf("listening: %v\n", gPortStr)
  }

  //INIT
  z := []APIAssembliesStruct{}
  z = append(z, APIAssembliesStruct{Name:"111", PDH:"1c20dd595e9fd3d8eefb281e314709ec+67", Locus:[]APILocusStruct{}})
  z[0].Locus = append(z[0].Locus, APILocusStruct{ ChromosomeName:"13", Indexing:0, StartPosition:0, EndPosition:0 })
  api_assemblies_init(z)

  api_tile_library_init()

  srv.Serve(listener)

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
