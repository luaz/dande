package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	//"math/rand"
	"net"
	"os"
	"strings"
	//"time"
)

var (
	Error *log.Logger
)

func initialize() {
	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func handle_command(conn net.Conn, input chan string) {
	defer close(input)
	defer conn.Close()

	for {
		line := <-input
		line = strings.TrimSpace(line)
		/*
			log.Print(line)
			min := 1000
			max := 5000
			num := rand.Intn(max-min) + min
			time.Sleep(time.Duration(num) * time.Millisecond)
		*/
		i := strings.Index(line, " ")
		var command, message string
		if i != -1 {
			command = line[:i]
			//param := line[i+1:]
		} else {
			command = line
		}

		switch command {
		case "quit":
			message = "Bye!"
			break
		default:
			message = "I don't understand ..."
		}

		send_output(conn, "> "+message+"\n")
		if message == "Bye!" {
			break
		}
	}

	log.Print("connection end")
}

func handle_input(conn net.Conn, input chan string) {
	//defer close(input)
	//defer conn.Close()

	for {
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			//Error.Println("Read error", err)
			break
		}
		input <- string(line)
	}
}

func send_output(conn net.Conn, line string) {
	io.Copy(conn, bytes.NewBufferString(line))
}

func main() {
	initialize()

	log.Print("listening on 5994 ...")
	sock, err := net.Listen("tcp", ":5994")
	if err != nil {
		Error.Println("Listen error", err)
		return
	}

	for {
		conn, err := sock.Accept()
		if err != nil {
			Error.Println("Accept error", err)
			continue
		}

		log.Print("connection start")
		channel := make(chan string, 10)
		go handle_input(conn, channel)
		go handle_command(conn, channel)
	}
}
