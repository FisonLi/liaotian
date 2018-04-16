package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	LOGFILE = "./mylog.text"
)

var onlineConns = make(map[string]net.Conn)

func CheckError(err error) {

	file, err1 := os.OpenFile(LOGFILE, os.O_RDWR|os.O_CREATE, 0)
	if err1 != nil {
		fmt.Println("open file wrong")
	}
	defer file.Close()

	logger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println("dd对对对ddd")

	if err != nil {
		panic(err)
	}

}

var mgQue = make(chan string, 1000)
var quitChan = make(chan bool)

func ProcessInfo(conn net.Conn) {

	//finally close conn
	defer func(conn net.Conn) {
		//first : delete the coon from sets, then close.
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlineConns, addr)
		conn.Close()

		for i := range onlineConns {
			fmt.Println("当前用户列表。。。", i)
		}
		if len(onlineConns) == 0 {
			fmt.Println("当前没有用户。。。")
		}
	}(conn)

	buf := make([]byte, 1024)

	//foreach reading
	for {
		//when its wrong , break
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		//reading from client coon,put coon's message into messageQueue
		if n != 0 {
			mg := string(buf[0:n])
			mgQue <- mg
		}

	}

}

//foreach reading messagequeue
func ConsumeMessage() {
	for {
		select {
		//if channal running, process data from client writing
		case mg := <-mgQue:

			//,when running,message writing into desc coon.
			//doProcessMessage(mg)
			doProcessMessage1(mg)
		case <-quitChan:
			break
		}
	}
}
func doProcessMessage1(s string) {
	fmt.Println(s)
}
func doProcessMessage(s string) {
	//con writing process

	split := strings.Split(s, "#")
	if len(split) > 1 {
		addr := split[0]
		sendMg := strings.Join(split[1:], "#")
		time.Now()
		sendMg = sendMg + "-time-" + (time.Now().String())

		addr = strings.Trim(addr, "")

		//if sets have the conn,writing message for the conn
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMg))
			if err != nil {
				fmt.Println("conn writing wrong", err)
			}
		}
	} else {
		split := strings.Split(s, "*")
		if strings.ToUpper(split[1]) == "LIST" {
			var ips string = ""
			for i := range onlineConns {
				ips = ips + i + "|"
			}
			if conn, ok := onlineConns[split[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("conn writing list wrong", err)
				}
			}
		}
	}
}

func main() {

	listen_socket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()

	fmt.Println("聊天服务开启。。。")

	//foreach reading ,channal not runing at defult
	go ConsumeMessage()

	//foreach listening

	for {
		conn, err := listen_socket.Accept()
		CheckError(err)

		addr := fmt.Sprintf("%s", conn.RemoteAddr())

		//when a coon coming,put it into sets.
		onlineConns[addr] = conn

		for i := range onlineConns {
			fmt.Println("用户列表...", i)
		}

		//process conn's message
		go ProcessInfo(conn)

	}
}
