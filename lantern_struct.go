package main

import "time"
import gouuid "github.com/nu7hatch/gouuid"

const gAPIVersionString = "0.1.0"

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
