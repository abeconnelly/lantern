package main

import "time"
import gouuid "github.com/nu7hatch/gouuid"

type LanternStatStruct struct {
  StartTime time.Time
  Requests int
  Failures int
  AvgRespMs []float64
  AvgRespDt float64
}

/*
type LanterDataStruct struct {
  CGF []cgf.CGF
}
*/

type LanternConnInfo struct {
  UUID *gouuid.UUID
}
