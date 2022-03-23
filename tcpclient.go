package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"fmt"
	"hash/crc64"
	"strconv"
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

func getBlockCRC() {
}

func sendMessage(s string, con net.Conn) {
    if _, err := con.Write([]byte(s)); err != nil {
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

//-------------// end init //-----------------//

	srcData := readBlock(src, bufferSize, offset)
	if srcData == nil {
		return //end of source file
	}

 	crc := crc64.Checksum(srcData, crcTable)
	request := fmt.Sprintf("req:%s:%d:%d:%d:", dst,bufferSize,offset,crc)
    fmt.Println(request)
	sendMessage(request, con)



//	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)
	//bf := []byte{001, 33}

	//if _, err = con.Write(bf); err != nil {
	//	log.Printf("failed to send the client request: %v\n", err)
	//}



//	for {
		// Waiting for the client request
		serverRequest, err := serverReader.ReadString('\n')
		switch err {
		case nil:
			if strings.TrimSpace(serverRequest) == "crc:false" {
				res := "data:" + string(srcData)
				log.Println(res)
				sendMessage(res, con)
				offset += bufferSize
			}

            if strings.TrimSpace(serverRequest) == "crc:true" {
				offset += bufferSize
                //send next request
            }
		}

		//check eof
		//read next block
		//send crc
		//get response, check crc
		//read next block or send data

//		switch err {
//		case nil:
//			clientRequest := strings.TrimSpace(clientRequest)
//			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
//				log.Printf("failed to send the client request: %v\n", err)
//			}
//		case io.EOF:
//			log.Println("client closed the connection")
//			return
//		default:
//			log.Printf("client error: %v\n", err)
//			return
//		}

		// Waiting for the server response
//		serverResponse, err := serverReader.ReadString('\n')

//		switch err {
//		case nil:
//			log.Println(strings.TrimSpace(serverResponse))
//		case io.EOF:
//			log.Println("server closed the connection")
//			return
//		default:
//			log.Printf("server error: %v\n", err)
//			return
//		}
//	}
}
