package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"strconv"
//	"fmt"
	"hash/crc64"
	"os"
)

const UID string = "1e028f50770445658114f05ba2b8ced5:"

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

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

func readBlock(f *os.File, size, offset int) []byte {
	buffer := make([]byte, size)

	n, err := f.ReadAt(buffer, int64(offset))
	if err == io.EOF {
		if n > 0 {
			return buffer[0:n]
		} else {
			return nil
		}
	}

	checkError(err)
	return buffer[0:n]
}

func sendMessage(s string, con net.Conn) {
	if _, err := con.Write([]byte(s)); err != nil {
		log.Printf("failed to respond to client: %v\n", err)
	}
}


func handleClientRequest(con net.Conn) {
	defer con.Close()

	crcTable := crc64.MakeTable(crc64.ISO)
	var offset int
	var dst *os.File
	defer dst.Close()

	clientReader := bufio.NewReader(con)
	readBuffer := make([]byte, 512)

	for {

		_, err := clientReader.Read(readBuffer)

		switch err {
		case nil:
			clientRequest := string(readBuffer)
			//log.Println(readBuffer)

			if strings.HasPrefix(clientRequest, UID) {
				clientRequest = strings.TrimPrefix(clientRequest, UID)

				data := strings.Split(clientRequest, ":")
				fileName := data[0]

				bufferSize, err := strconv.Atoi(data[1])
				checkError(err)
				readBuffer = make([]byte, bufferSize)

				offset, err = strconv.Atoi(data[2])
				checkError(err)

				crc, err := strconv.ParseUint(data[3], 0, 64)
				checkError(err)

				dst, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
				checkError(err)

				dstData := readBlock(dst, bufferSize, offset)

				if crc64.Checksum(dstData, crcTable) != crc {
					sendMessage("crc:false\n", con)
				} else {
					sendMessage("crc:true\n", con)
				}

				log.Printf("%s:%d:%d:%d\n", fileName, bufferSize, offset, crc)
			} else {
				//if strings.HasPrefix(clientRequest, "data:") {
				//clientRequest = strings.TrimPrefix(clientRequest, "data:")
				//data := []byte(clientRequest)

				_, err := dst.WriteAt(readBuffer, int64(offset))
				checkError(err)
				//log.Println(string(readBuffer))
				log.Printf("%d bytes recorded\n", len(readBuffer))
                //sendMessage("data:true\n", con)
			}

		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
