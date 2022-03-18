package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"log"
	"hash/crc64"
)

func checkError(e error) {
    if e != nil {
		log.Fatal(e)
		return
    }
}

func readBlock(f *os.File, offset, size int) []byte {
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

	dst, err := os.OpenFile(os.Args[3], os.O_RDWR|os.O_CREATE, 0644)
    checkError(err)
	defer dst.Close()

	crcTable := crc64.MakeTable(crc64.ISO)
	idx := 0
	writeBytes := 0

	for {

		srcData := readBlock(src, idx, bufferSize)
		if srcData == nil { break }

		dstData := readBlock(dst, idx, bufferSize)
		if crc64.Checksum(srcData, crcTable) != crc64.Checksum(dstData, crcTable) {
			_, err := dst.WriteAt(srcData, int64(idx))
			checkError(err)

			writeBytes += len(srcData)
			//fmt.Println(string(srcData))
		}

		idx += bufferSize

	}

	fmt.Printf("Write %d blocks, %d bytes\n", writeBytes/bufferSize, writeBytes)
}
