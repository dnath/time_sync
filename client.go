package main

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

// Consts

const (
	NANOSEC_PER_SEC = 1000000000
	BEGIN           = -1
	END             = 1
)

// Types

type TimeServer struct {
	ip   string
	port int
}

type TimeServerResponseInfo struct {
	localSendTime, serverTime, rtt float64
}

type IntervalOffset struct {
	serverIndex int
	offset      float64
	offsetType  int
}

type IntervalOffsetArray []IntervalOffset

// Functions

func (arr IntervalOffsetArray) Len() int {
	return len(arr)
}

func (arr IntervalOffsetArray) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr IntervalOffsetArray) Less(i, j int) bool {
	if arr[i].offset == arr[j].offset {
		if arr[i].offsetType == arr[j].offsetType {
			return i < j
		}
		return arr[i].offsetType < arr[j].offsetType
	}
	return arr[i].offset < arr[j].offset
}

func checkError(err error, willExit bool) {
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

func GetConfig() []TimeServer {
	return []TimeServer{
		TimeServer{ip: "54.172.168.244", port: 5000},
		TimeServer{ip: "54.169.67.45", port: 5000},
		TimeServer{ip: "54.207.15.207", port: 5000},
		TimeServer{ip: "54.191.73.92", port: 5000},
		TimeServer{ip: "128.111.44.106", port: 12291},
	}
}

func RequestTimeSync(server TimeServer) (info TimeServerResponseInfo) {
	var message = "\n"
	var buf [256]byte

	var byteBuffer bytes.Buffer
	fmt.Fprintf(&byteBuffer, "%s:%d", server.ip, server.port)

	addr := byteBuffer.String()
	// fmt.Printf("server = %s", addr)

	serverAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkError(err, true)

	info.localSendTime = float64(time.Now().UnixNano()) / NANOSEC_PER_SEC

	conn, err := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()
	checkError(err, true)

	_, err = conn.Write([]byte(message))
	checkError(err, true)

	numBytesRead, err := conn.Read(buf[0:])
	checkError(err, true)

	localRecvTime := float64(time.Now().UnixNano()) / NANOSEC_PER_SEC
	// fmt.Printf("localSendTime  = %f\n", info.localSendTime)
	// fmt.Printf("localRecvTime  = %f\n", localRecvTime)

	info.rtt = localRecvTime - info.localSendTime
	// fmt.Printf("rtt            = %f\n", info.rtt)

	info.serverTime, err = strconv.ParseFloat(string(buf[0:numBytesRead]), 64)
	// fmt.Printf("serverTime     = %f\n", info.serverTime)

	return info
}

func GetMaxIntersection(offsets IntervalOffsetArray) (maxIntersections int,
	maxIntersectionStart IntervalOffset,
	maxIntersectionEnd IntervalOffset) {

	var intersectionStart, intersectionEnd IntervalOffset
	var isMaxIntersection bool = false
	var numIntersections int = 0
	var maxIntersectionMidpoint float64 = 0

	sort.Sort(offsets)
	var globalMidpoint = (offsets[0].offset + offsets[len(offsets)-1].offset) / 2

	// fmt.Printf("sorted offsets =\n")
	// for _, offset := range offsets {
	// 	fmt.Printf("%d (%f, %d)\n", offset.serverIndex, offset.offset, offset.offsetType)
	// }

	for _, offset := range offsets {
		if offset.offsetType == BEGIN {
			intersectionStart = offset
		} else {
			intersectionEnd = offset
		}

		numIntersections -= offset.offsetType
		if numIntersections > maxIntersections {
			maxIntersections = numIntersections
			maxIntersectionStart = intersectionStart
			isMaxIntersection = true
		}

    // fmt.Printf("numIntersections = %d isMaxIntersection = %s offset = {%f, %d}\n", numIntersections, 
    //                                                                                isMaxIntersection, 
    //                                                                                offset.offset, 
    //                                                                                offset.offsetType)

		if numIntersections == maxIntersections-1 && !isMaxIntersection && offset.offsetType == END {
      var midpoint = (intersectionStart.offset + offset.offset) / 2
      // fmt.Printf("midpoint = %f\n", midpoint)
      // fmt.Printf("global midpoint = %f\n", globalMidpoint)
			
      if math.Abs(midpoint-globalMidpoint) < math.Abs(midpoint-maxIntersectionMidpoint) {
        maxIntersectionStart = intersectionStart
        isMaxIntersection = true
			}
		}

		if offset.offsetType == END && isMaxIntersection && numIntersections == maxIntersections-1 {
			maxIntersectionEnd = intersectionEnd
			maxIntersectionMidpoint = (maxIntersectionStart.offset + maxIntersectionEnd.offset) / 2
			isMaxIntersection = false
		}
	}

	// fmt.Printf("maxIntersections = %d, start = %f, end = %f\n\n", maxIntersections,
	// 	                                                            maxIntersectionStart.offset,
	// 	                                                            maxIntersectionEnd.offset)

	return maxIntersections, maxIntersectionStart, maxIntersectionEnd
}

func Run(trial int, servers []TimeServer) {
	numServers := len(servers)

	offsets := make(IntervalOffsetArray, numServers*2)
	info := make([]TimeServerResponseInfo, numServers)

	var updatedTime, estimationError float64 = 0, 0
	var minLocalSendTime float64 = 0.0

	for serverIndex, server := range servers {
		info[serverIndex] = RequestTimeSync(server)

		// Time Delta shift for pretending that the requests were all sent at the same time, minLocalSendTime
		if serverIndex == 0 {
			minLocalSendTime = info[serverIndex].localSendTime
		} else {
			var delta = info[serverIndex].localSendTime - minLocalSendTime
			info[serverIndex].localSendTime = minLocalSendTime
			info[serverIndex].serverTime -= delta
		}

		// fmt.Printf("%f %f %f\n\n", info[serverIndex].localSendTime,
		//                            info[serverIndex].serverTime,
		//                            info[serverIndex].rtt)

		offsets[2*serverIndex] = IntervalOffset{serverIndex: serverIndex,
			offset:     info[serverIndex].serverTime - info[serverIndex].rtt,
			offsetType: BEGIN}

		offsets[2*serverIndex+1] = IntervalOffset{serverIndex: serverIndex,
			offset:     info[serverIndex].serverTime + info[serverIndex].rtt,
			offsetType: END}
	}

	// for _, offset := range offsets {
	//  fmt.Printf("%s\n", offset)
	// }

	maxIntersections, maxIntersectionStart, maxIntersectionEnd := GetMaxIntersection(offsets)

	fmt.Printf("\nTrial #%d:\n", trial)
	fmt.Printf("current time = %s\n", GetHumanReadableTime(minLocalSendTime))

	for serverIndex := 0; serverIndex < numServers; serverIndex++ {
		fmt.Printf("[%s:%d] server time = %s, rtt = %fs\n", servers[serverIndex].ip,
			servers[serverIndex].port,
			GetHumanReadableTime(info[serverIndex].serverTime),
			info[serverIndex].rtt)
	}

	if maxIntersections > 1 {
		updatedTime = (maxIntersectionStart.offset + maxIntersectionEnd.offset) / 2
		estimationError = (maxIntersectionEnd.offset - maxIntersectionStart.offset) / 2
		fmt.Printf("updated time = %s, error = %fs\n", GetHumanReadableTime(updatedTime), estimationError)

	} else {
		fmt.Printf("No Intersection !")
	}
}

func ClientMain() {
	servers := GetConfig()
	for trial := 0; ; trial++ {
		Run(trial, servers)
	}
}

func main() {
  ClientMain()
}

