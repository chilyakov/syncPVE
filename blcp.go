package main

import (
	"fmt"
	"hash/crc64"
	"io"
	"log"
	"os"
	"strconv"
	//	"math"
)

var writeBytes, offset, writeBlocks int

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

func readBlock(f *os.File, size int) []byte {
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

func syncFiles(src *os.File, dst *os.File, size int) bool {

	crcTable := crc64.MakeTable(crc64.ISO)
	srcData := readBlock(src, size)
	if srcData == nil {
		return true //end of source file
	}

	dstData := readBlock(dst, size)
	if crc64.Checksum(srcData, crcTable) != crc64.Checksum(dstData, crcTable) {
		_, err := dst.WriteAt(srcData, int64(offset))
		checkError(err)
		writeBytes += len(srcData)
		writeBlocks ++
		fmt.Printf("block %d was recorded\n", writeBlocks)
	}

	offset += size
	return false
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

	for {

		if syncFiles(src, dst, bufferSize) {
			fmt.Printf("total %d blocks, %d bytes was recorded\n", writeBlocks, writeBytes)
			return
		}

	}

}
