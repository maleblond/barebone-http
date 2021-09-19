package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const PORT = "3000"

func main() {
	l, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	fmt.Println("Listening on port: " + PORT)

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	request := newRequest()

	reader := bufio.NewReader(conn)

	index := 0
	for {
		line, _ := reader.ReadString('\n')
		line = strings.ReplaceAll(line[:len(line)-1], "\r", "")

		if line == "" {
			break
		}

		if index == 0 {
			parseVerbAndPath(line, request)
		} else {
			parseHeader(line, request)
		}
		index++
	}

	contentLength, _ := strconv.Atoi(request.headers["Content-Length"])
	body := make([]byte, contentLength)

	io.ReadFull(reader, body)
	request.body = string(body)

	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\n\nHere is a response with your passed body: %v\n", request.body)))

	conn.Close()
}

type request struct {
	verb    string
	headers map[string]string
	path    string
	body    string
}

func newRequest() *request {
	req := new(request)

	req.headers = make(map[string]string)

	return req
}

func parseVerbAndPath(line string, request *request) {
	regex := regexp.MustCompile(`^(?P<verb>[A-Z]*)\s(?P<path>[^\s]*)\sHTTP\/.*$`)

	matches := regex.FindStringSubmatch(line)

	if matches != nil {
		request.verb = matches[1]
		request.path = matches[2]
	}
}

func parseHeader(line string, request *request) {
	regex := regexp.MustCompile(`^(?P<key>[^:]+):\s*(?P<val>.+)$`)

	matches := regex.FindStringSubmatch(line)

	if matches != nil {
		request.headers[matches[1]] = matches[2]
	}
}
