package main

import (
	"hash/crc64"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
//	"os/signal"
//	"syscall"
)

const UID string = "1e028f50770445658114f05ba2b8ced5:"

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

func main() {

	flog, err := os.OpenFile("/opt/blcp/blcps.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()

	//errorLog := log.New(flog, "ERROR\t", log.Ldate|log.Ltime)

    host, err := net.ResolveTCPAddr("tcp4", "0.0.0.0"+":7231")
    if err != nil {
        log.Fatalln(err)
    }

    listener, err := net.ListenTCP("tcp", host)
    if err != nil {
        log.Fatalln(err)
    }
    defer listener.Close()

	infoLog := log.New(flog, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf("* * * Start server * * *")

//	signalChanel := make(chan os.Signal, 1)
//	signal.Notify(signalChanel,
//		syscall.SIGHUP,
//		syscall.SIGINT,
//		syscall.SIGTERM,
//		syscall.SIGQUIT)

//	exit_chan := make(chan int)
//	go func() {
//		for {
//			s := <-signalChanel
//			switch s {
			// kill -SIGHUP XXXX [XXXX - идентификатор процесса для программы]
//			case syscall.SIGHUP:
//				infoLog.Println("Signal hang up triggered.")

				// kill -SIGINT XXXX или Ctrl+c  [XXXX - идентификатор процесса для программы]
//			case syscall.SIGINT:
//				infoLog.Println("Signal interrupt triggered.")
//				exit_chan <- 1

				// kill -SIGTERM XXXX [XXXX - идентификатор процесса для программы]
//			case syscall.SIGTERM:
//				infoLog.Println("Signal terminte triggered.")
//				exit_chan <- 1

				// kill -SIGQUIT XXXX [XXXX - идентификатор процесса для программы]
//			case syscall.SIGQUIT:
//				infoLog.Println("Signal quit triggered.")
//				exit_chan <- 0
//
//			default:
//				infoLog.Println("Unknown signal.")
//				exit_chan <- 1
//			}
//		}
//	}()
//	exitCode := <-exit_chan
//	os.Exit(exitCode)

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con, flog)
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

func sendMessage(s string, con net.Conn, errorLog *log.Logger) {
	if _, err := con.Write([]byte(s)); err != nil {
		errorLog.Printf("failed to respond to client: %v\n", err)
	}
}

func handleClientRequest(con net.Conn, flog *os.File) {
	defer con.Close()

	infoLog := log.New(flog, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(flog, "ERROR\t", log.Ldate|log.Ltime)

	infoLog.Printf("Connect from %s", con.RemoteAddr())

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
					infoLog.Printf("%s sending file %s", con.RemoteAddr(), dst.Name())
				}
				dstData := readBlock(dst, blockSize, offset)

				if crc64.Checksum(dstData, crcTable) != crc {
					sendMessage("crc:false\n", con, errorLog)
				} else {
					sendMessage("crc:true\n", con, errorLog)
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
						//log.Println("debug line 148")

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
								sendMessage("crc:false\n", con, errorLog)
							} else {
								sendMessage("crc:true\n", con, errorLog)
							}
						} else {
							errorLog.Fatal("error line 181 (detect request packet)")
						}
					}

				} else {
					errorLog.Println("debug line 187 (unkown buffer size?)")
					blockOffset = 0
				}
			}

		case io.EOF:
			//log.Println("max buffer size:", maxBuffer)
			infoLog.Printf("%s closed file %s by EOF. %d bytes was recorded", con.RemoteAddr(), dst.Name(), bytesRec)
			maxBuffer = 0
			bytesRec = 0
			dst.Close()
			return
		default:
			errorLog.Printf("Connection %s error: %v\n", con.RemoteAddr(), err)
			return
		}
	}
}
