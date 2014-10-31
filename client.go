package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
  "encoding/json"
	"time"
  "io/ioutil"
  flags "github.com/jessevdk/go-flags"
)

// Consts

const (
	RHO    = 15 // local clock drift rate = 15 ms/min
	DELTA  = 1  // drift tolerance = 1 ms
)

// Types

type TimeServerResponseInfo struct {
	localSendTime, serverTime, rtt float64
  success bool
}

// Functions

func GetConfig(serverSettingsFilename string) []TimeServer {
  readBytes, err := ioutil.ReadFile(serverSettingsFilename)
  CheckError(err, true)

  var jsonObject map[string]TimeServerArray
  err = json.Unmarshal(readBytes, &jsonObject)
  CheckError(err, true)

  return jsonObject["servers"]
}

func RequestTimeSync(server TimeServer) (info TimeServerResponseInfo) {
	var message = "\n"
	var buf [256]byte
	addr := server.String()

	serverAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    info.success = false
    return info
  }

	info.localSendTime = float64(time.Now().UnixNano()) / NANOSEC_PER_SEC

	conn, err := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()
	if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    info.success = false
    return info
  }

	_, err = conn.Write([]byte(message))
	if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    info.success = false
    return info
  }

	numBytesRead, err := conn.Read(buf[0:])
	if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
    info.success = false
    return info
  }

	localRecvTime := float64(time.Now().UnixNano()) / NANOSEC_PER_SEC
	info.rtt = localRecvTime - info.localSendTime
	info.serverTime, err = strconv.ParseFloat(string(buf[0:numBytesRead]), 64)
  info.success = true

	return info
}

func Run(trial int, servers []TimeServer) {
	numServers := len(servers)

	offsets := make(IntervalOffsetArray, numServers*2)
	info := make([]TimeServerResponseInfo, numServers)

	var updatedTime, estimationError float64 = 0, 0
	var minLocalSendTime float64 = 0.0

	for serverIndex, server := range servers {
    response := RequestTimeSync(server)
    if !response.success {
      offsets[2*serverIndex] = IntervalOffset{serverIndex: serverIndex, offset: INVALID_OFFSET, offsetType: 0}
      offsets[2*serverIndex] = IntervalOffset{serverIndex: serverIndex, offset: INVALID_OFFSET, offsetType: 0}
      continue
    }

    info[serverIndex] = response

		// Time Delta shift for pretending that the requests were all sent at the same time, minLocalSendTime
		if minLocalSendTime == 0.0 {
			minLocalSendTime = info[serverIndex].localSendTime
		} else {
			var delta = info[serverIndex].localSendTime - minLocalSendTime
			info[serverIndex].localSendTime = minLocalSendTime
			info[serverIndex].serverTime -= delta
		}

		offsets[2*serverIndex] = IntervalOffset{serverIndex: serverIndex,
			                                      offset: info[serverIndex].serverTime - info[serverIndex].rtt,
			                                      offsetType: BEGIN}

		offsets[2*serverIndex+1] = IntervalOffset{serverIndex: serverIndex,
			                                        offset: info[serverIndex].serverTime + info[serverIndex].rtt,
			                                        offsetType: END}
	}

	maxIntersections, maxIntersectionStart, maxIntersectionEnd := GetMaxIntersection(offsets)

	fmt.Printf("\nTrial #%d:\n", trial)
	fmt.Printf("current time = %s\n", GetHumanReadableTime(minLocalSendTime))

	for serverIndex := 0; serverIndex < numServers; serverIndex++ {
    if info[serverIndex].success {
      fmt.Printf("[%s:%d] server time = %s, rtt = %fs\n", servers[serverIndex].IP,
                                                          servers[serverIndex].Port,
                                                          GetHumanReadableTime(info[serverIndex].serverTime),
                                                          info[serverIndex].rtt)
    }
	}

	if maxIntersections > 1 {
		updatedTime = (maxIntersectionStart.offset + maxIntersectionEnd.offset) / 2
		estimationError = (maxIntersectionEnd.offset - maxIntersectionStart.offset) / 2
		fmt.Printf("updated time = %s, error = %fs\n", GetHumanReadableTime(updatedTime), estimationError)

	} else {
		fmt.Printf("No Intersection !")
	}
}

func ClientMain(serverSettingsFilename string) {
	servers := GetConfig(serverSettingsFilename)

	timeToSleep := time.Duration((DELTA*60)/(2*RHO)) * time.Second
	for trial := 0; ; trial++ {
		Run(trial, servers)
		time.Sleep(timeToSleep)
	}
}

func main() {
	opts := new(
		struct {
			Settings string `short:"s" long:"settings" default:"conf/servers.json" description:"Server Settings Filename"`
		})

	parser := flags.NewParser(opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}

	ClientMain(opts.Settings)
}
