package main

import (
	"io"
	"log"
	"net"
)

const DbAddr = "127.0.0.1:3306"
const ProxyAddr = "127.0.0.1:3307"
const ComQuery = byte(0x03)

func main() {
	socket, err := net.Listen("tcp", ProxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listing on %s", ProxyAddr)

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	dbConnect, err := net.Dial("tcp", DbAddr)
	if err != nil {
		log.Printf("Error connecting to DB: %v\n", err)
	}

	defer dbConnect.Close()

	// forward from proxy to db
	go func(src, dst net.Conn) {
		buffer := make([]byte, 4096)
		// see https://www.oreilly.com/library/view/understanding-mysql-internals/0596009577/ch04.html
		for {
			n, _ := src.Read(buffer)
			if n > 5 {
				switch buffer[4] {
				case ComQuery:
					query := string(buffer[5:n])
					log.Printf("original query: %s\n", query)
				}
			}

			// forward
			_, err := dst.Write(buffer[0:n])
			if err != nil {
				log.Printf("Error writing to destination: %v\n", err)
			}
		}
	}(conn, dbConnect)

	// forward everything from db to client until receive EOF
	if _, err := io.Copy(conn, dbConnect); err != nil {
		log.Printf("Error copying data: %v\n", err)
	}
}
