package main

import (
  "fmt"
  "os"
  "time"
  "bytes"
)

const NANOSEC_PER_SEC = 1000000000

// Types

type TimeServer struct {
  IP string `json: "ip"`
  Port int `json: "port"`
}

type TimeServerArray []TimeServer

// Functions

func (ts TimeServer) String() string {
  var ByteBuffer bytes.Buffer
  fmt.Fprintf(&ByteBuffer, "%s:%d", ts.IP, ts.Port)
  return ByteBuffer.String()
}

func CheckError(err error, willExit bool) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    if willExit {
      os.Exit(1)
    }
  }
}

func GetHumanReadableTime(seconds float64) time.Time {
  return time.Unix(0, int64(seconds*NANOSEC_PER_SEC))
}

