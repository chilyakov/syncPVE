package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		//rq := make([]byte, 1)
		//_, e := con.Read(rq)
		//if e != nil {
		//	log.Println(e)
		//	return
		//}
		//log.Println(int(rq[0]))
		//var res bool
		clientRequest, err := clientReader.ReadString('\n')
		switch err {
		case nil:
//			clientRequest := strings.TrimSpace(clientRequest)
//			res = strings.HasPrefix(clientRequest, "req:")
//			log.Println(res)

			if strings.HasPrefix(clientRequest, "req:") {
				clientRequest = strings.TrimPrefix(clientRequest, "req:")

				data := strings.Split(clientRequest, ":")
				fileName := data[0]
				var bufferSize, offset int
				var crc int64
				if i, err := strconv.Atoi(data[1]); err == nil {
					bufferSize = i
				}
                if i, err := strconv.Atoi(data[2]); err == nil {
                    offset = i
                }
                if i, err := strconv.ParseInt(data[3], 10, 64); err == nil {
                    crc = i
                }

//				if readFile(fileName, bufferSize, offset) != crc

				if _, err = con.Write([]byte("req:true\n")); err != nil {
					log.Printf("failed to respond to client: %v\n", err)
				}
				log.Printf("filename: %s, bufferSize: %d, offset: %d, crc: %d", fileName, bufferSize, offset, crc)

			} else if strings.HasPrefix(clientRequest, "data:") {
				clientRequest = strings.TrimPrefix(clientRequest, "data:")
				data := []byte(clientRequest)



                if _, err = con.Write([]byte("data:true\n")); err != nil {
                    log.Printf("failed to respond to client: %v\n", err)
                }
				log.Println(data)
				log.Println(string(data))

			} else {
				if _, err = con.Write([]byte("err:true\n")); err != nil {
					log.Printf("failed to respond to client: %v\n", err)
				}
				log.Fatalln("unknown error!")
				return
			}

//			if clientRequest == ":QUIT" {
//				log.Println("client requested server to close the connection so closing")
//				return
//			} else {
//				log.Println(clientRequest)
//			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		// Responding to the client request
//		if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
//			log.Printf("failed to respond to client: %v\n", err)
//		}
	}
}
