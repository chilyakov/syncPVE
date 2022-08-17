package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
)

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
		return
	}
}

func checkMigrate() {
	command := "/opt/syncvm/check_migrate.sh > /dev/null 2>&1"

	cmd := exec.Command("bash", "-c", command)
	err := cmd.Start()
	checkError(err)

}

func checkMigrateCancel() {
	command := "/opt/syncvm/check_migrate_cancel.sh > /dev/null 2>&1"

	cmd := exec.Command("bash", "-c", command)
	err := cmd.Start()
	checkError(err)

}

func main() {

	//flog, err := os.OpenFile("syncvm.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer flog.Close()

	//errorLog := log.New(flog, "ERROR\t", log.Ldate|log.Ltime)
	checkMigrate()
	checkMigrateCancel()

	host, err := net.ResolveTCPAddr("tcp4", "0.0.0.0"+":7011")
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", host)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	//infoLog := log.New(flog, "INFO\t", log.Ldate|log.Ltime)
	//infoLog.Printf("* * * Start server * * *")

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {
	defer con.Close()

	//infoLog := log.New(flog, "INFO\t", log.Ldate|log.Ltime)
	//infoLog.Printf("Connect from %s", con.RemoteAddr())

	readBuffer := make([]byte, 16)
	_, err := con.Read(readBuffer)

	switch err {
	case nil:

		data := strings.Split(string(readBuffer[:]), ":")
		//infoLog.Printf("VMID: %s", data[0])

		vmid := " " + data[0]
		status := data[1]
		command := ""

		if status == "start" {
			command = "/opt/syncvm/start_sync.sh" + vmid

		} else if status == "stop" {
			command = "/opt/syncvm/stop_sync.sh" + vmid

		} else {
			fmt.Printf("Unknown status recieved (%s). Program close connection.\n", status)
			return
		}

		cmd := exec.Command("bash", "-c", command)
		err := cmd.Start()
		checkError(err)

		err = cmd.Wait()
		fmt.Printf("Command finished with error: %v", err)

	case io.EOF:
		return
	default:
		return
	}
}
