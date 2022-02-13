package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	port    int     = 9094
	session Session = InitSession()
)

func init() {
	flag.IntVar(&port, "port", port, "Port used for listenning connections")
	flag.Parse()

	// Validate flags
	if port < 1024 {
		log.Error("Port number must be > 1024")
		os.Exit(101)
	}
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("Cannot start server, err: %v", err)
		os.Exit(102)
	}
	wg := sync.WaitGroup{}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("Error accepting request from %s, err: %v", conn.RemoteAddr().String(), err)
			conn.Close()
		}
		go func() {
			defer wg.Done()
			defer conn.Close()
			wg.Add(1)
			conn.Write([]byte(handleRequest(conn)))
		}()
	}
}

func handleRequest(conn net.Conn) string {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.WithFields(log.Fields{
			"remote address": conn.RemoteAddr().String(),
			"error":          err,
		}).Error("Error reading request body")
	}
	text := strings.Trim(string(buffer), "\x00")
	text = strings.Replace(text, "\n", "", -1)
	result, err := session.ExecuteCommand(text)
	if err != nil {
		log.WithFields(log.Fields{
			"remote address": conn.RemoteAddr().String(),
			"error":          err,
			"request":        text,
		}).Error("Error while processing request")
		return err.Error()
	}
	return result
}
