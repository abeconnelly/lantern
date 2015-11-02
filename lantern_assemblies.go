package main

import "io"
import "fmt"
import "net/http"

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

var g_assemblies []APIAssembliesStruct
var g_assemblies_idx_map map[string]int

func api_assemblies_init(assembly_info []APIAssembliesStruct) {
  g_assemblies = make([]APIAssembliesStruct,0,27)
  g_assemblies_idx_map = make(map[string]int)
  for i:=0; i<len(assembly_info); i++ {
    g_assemblies = append(g_assemblies, assembly_info[i])
    g_assemblies_idx_map[assembly_info[i].PDH] = i
  }
}


func handle_assemblies_id(w http.ResponseWriter, r *http.Request) {
  var body_reader io.Reader = r.Body ; _ = body_reader

  assembly_name := "123"

  var idx int
  var ok bool
  if idx,ok = g_assemblies_idx_map[assembly_name]; !ok {
    io.WriteString(w,"[]")
    return
  }

  io.WriteString(w,"[")
  for i:=0; i<len(g_assemblies[idx].Locus); i++ {
    if i>0 { io.WriteString(w,",") }
    io.WriteString(w,"{")
    io.WriteString(w, fmt.Sprintf(`"assembly-name":"%s",`, g_assemblies[idx].Name))
    io.WriteString(w, fmt.Sprintf(`"assembly-pdh":"%s",`, g_assemblies[idx].PDH))
    io.WriteString(w, fmt.Sprintf(`"chromosome-name":"%s",`, g_assemblies[idx].Locus[i].ChromosomeName))
    io.WriteString(w, fmt.Sprintf(`"indexing":%d,`, g_assemblies[idx].Locus[i].Indexing))
    io.WriteString(w, fmt.Sprintf(`"start-position":%d,`, g_assemblies[idx].Locus[i].StartPosition))
    io.WriteString(w, fmt.Sprintf(`"end-position":%d`, g_assemblies[idx].Locus[i].EndPosition))
    io.WriteString(w,"}")
  }
  io.WriteString(w,"]")


  /*
  fmt.Printf("handle_assemblies_id:\n")
  fmt.Printf("  RequestURI: %s\n", r.RequestURI)
  fmt.Printf("  Host: %s\n", r.Host)
  fmt.Printf("  Header: %s\n", r.Header)
  */
}

func handle_assemblies(w http.ResponseWriter, r *http.Request) {
  var body_reader io.Reader = r.Body ; _ = body_reader

  io.WriteString(w,"[")
  for i:=0; i<len(g_assemblies); i++ {
    if i>0 { io.WriteString(w,",") }
    io.WriteString(w,"{")
    io.WriteString(w, fmt.Sprintf(`"assembly-name":"%s",`, g_assemblies[i].Name))
    io.WriteString(w, fmt.Sprintf(`"assembly-pdh":"%s"`, g_assemblies[i].PDH))
    io.WriteString(w,"}")
  }
  io.WriteString(w,"]")


  /*
  fmt.Printf("handle_assemblies:\n")
  fmt.Printf("  RequestURI: %s\n", r.RequestURI)
  fmt.Printf("  Host: %s\n", r.Host)
  fmt.Printf("  Header: %s\n", r.Header)
  */
}
