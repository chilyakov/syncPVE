package main

import (
	//	"bufio"
	"hash/crc64"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	//	"io/ioutil"
)

const UID string = "1e028f50770445658114f05ba2b8ced5:"

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

func main() {
	host, err := net.ResolveTCPAddr("tcp4", "0.0.0.0"+":7231")
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", host)
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

func readBlock(f *os.File, size int, offset uint64) []byte {
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
	var offset, blockOffset uint64
	var blockSize, maxBuffer, bytesRec int
	var dst *os.File
	defer dst.Close()

	readBuffer := make([]byte, 512)

	for {

		bytes, err := con.Read(readBuffer)
		if bytes > maxBuffer {
			maxBuffer = bytes
			//log.Println(maxBuffer)
		}

		switch err {
		case nil:

			if string(readBuffer[0:33]) == UID {
				data := strings.Split(string(readBuffer[33:]), ":")
				fileName := data[0]

				blockSize, err = strconv.Atoi(data[1])
				checkError(err)
				readBuffer = make([]byte, blockSize)
				blockOffset = 0

				offset, err = strconv.ParseUint(data[2], 0, 64)
				checkError(err)

				crc, err := strconv.ParseUint(data[3], 0, 64)
				checkError(err)

				_, err = dst.Stat()
				if err != nil {
					dst, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
					checkError(err)
				}
				dstData := readBlock(dst, blockSize, offset)

				if crc64.Checksum(dstData, crcTable) != crc {
					sendMessage("crc:false\n", con)
				} else {
					sendMessage("crc:true\n", con)
				}

				//log.Printf("%s:%d:%d:%d\n", fileName, blockSize, offset, crc)
				//log.Println(bytesCount, blockSize)
			} else {
				//log.Println(bytes, blockSize, offset, blockOffset)

				if blockOffset < uint64(blockSize) {
					n, err := dst.WriteAt(readBuffer[:bytes], int64(offset+blockOffset))
					checkError(err)
					if n > 0 {
						bytesRec += n
						//log.Printf("write %d bytes, %d offset\n", n, offset+blockOffset)
					}
					blockOffset += uint64(bytes)

					// если в конце буфера оказался пакет со следующим запросом от клиента
					if blockOffset > uint64(blockSize) {
						//log.Println("debug line 131")

						tmp := blockOffset - uint64(blockSize)
						blck := readBuffer[bytes-int(tmp):]

						if string(blck[:33]) == UID {
							data := strings.Split(string(blck[33:]), ":")
							fileName := data[0]

							blockSize, err = strconv.Atoi(data[1])
							checkError(err)
							readBuffer = make([]byte, blockSize)
							blockOffset = 0

							offset, err = strconv.ParseUint(data[2], 0, 64)
							checkError(err)

							crc, err := strconv.ParseUint(data[3], 0, 64)
							checkError(err)

							_, err = dst.Stat()
							if err != nil {
								dst, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
								checkError(err)
							}
							dstData := readBlock(dst, blockSize, offset)

							if crc64.Checksum(dstData, crcTable) != crc {
								sendMessage("crc:false\n", con)
							} else {
								sendMessage("crc:true\n", con)
							}
						} else {
							log.Fatal("error line 164 (detect request packet)")
						}
					}

				} else {
					log.Println("debug line 141")
					blockOffset = 0
				}
			}

		case io.EOF:
			//log.Println("max buffer size:", maxBuffer)
			log.Printf("client closed the connection by EOF. %d bytes was recorded", bytesRec)
			maxBuffer = 0
			bytesRec = 0
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
