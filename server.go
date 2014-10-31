package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
  flags "github.com/jessevdk/go-flags"
)

// Functions

func HandleClientRequest(conn *net.UDPConn) {
	fmt.Println("Processing request from ", conn.LocalAddr().String())
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])

	CheckError(err, false)

	var CurrentTime float64 = float64(time.Now().UnixNano()) / NANOSEC_PER_SEC
	conn.WriteToUDP([]byte(strconv.FormatFloat(CurrentTime, 'f', 6, 64)), addr)
}

func ServerMain(ts TimeServer) {
	UdpAddr := net.UDPAddr{
		Port: ts.Port,
		IP:   net.ParseIP(ts.IP),
	}

	conn, err := net.ListenUDP("udp4", &UdpAddr)
	CheckError(err, true)

	fmt.Println("Server running at", ts.String())

	for {
		HandleClientRequest(conn)
	}
}

func main() {
	opts := new(
		struct {
			Port int    `short:"p" long:"port" default:"6000" description:"Server Port"`
			IP   string `short:"i" long:"ip" default:"127.0.0.1" description:"Server IPv4 Address"`
		})

	parser := flags.NewParser(opts, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}

	ts := TimeServer{IP: opts.IP, Port: opts.Port}
	ServerMain(ts)
}
