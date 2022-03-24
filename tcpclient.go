package main

import (
	"bufio"
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const UID string = "1e028f50770445658114f05ba2b8ced5:"

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
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

func sendMessageBytes(b []byte, con net.Conn) {
	if _, err := con.Write(b); err != nil {
		log.Printf("failed to respond to client: %v\n", err)
	}
}

func main() {
	arguments := os.Args
	if len(arguments) != 4 {
		fmt.Println("<buffer size> <file src> <file dst>")
		return
	}

	bufferSize, err := strconv.Atoi(os.Args[1])
	checkError(err)

	src, err := os.Open(os.Args[2])
	checkError(err)
	defer src.Close()

	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	dst := os.Args[3]
	crcTable := crc64.MakeTable(crc64.ISO)
	offset := 0

	serverReader := bufio.NewReader(con)

	//end init

	for {

		srcData := readBlock(src, bufferSize, offset)
		if srcData == nil {
			return //end of source file
		}

		crc := crc64.Checksum(srcData, crcTable)
		request := fmt.Sprintf("%s%s:%d:%d:%d:", UID, dst, len(srcData), offset, crc)
		fmt.Println(request)
		sendMessage(request, con)

		serverRequest, err := serverReader.ReadString('\n')
		switch err {
		case nil:
			if strings.TrimSpace(serverRequest) == "crc:false" {
				sendMessageBytes(srcData, con)
				offset += bufferSize
				break
			}

			if strings.TrimSpace(serverRequest) == "crc:true" {
				offset += bufferSize
				break
			}

		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}

	}
}
