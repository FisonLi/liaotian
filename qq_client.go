package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	LOGFILE1 = "./mylog.text"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
	/*
		file, err1:= os.OpenFile(LOGFILE1, os.O_RDWR|os.O_CREATE, 0)
		if err1!=nil{
			fmt.Println("open file wrong")
		}
		defer file.Close()

		logger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
		logger.Println(e)
		if e!=nil {
			panic(e)
		}*/
}

func MessageSend(conn net.Conn) {
	var input string
	//foreach the conn writing
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		input = string(line)

		if strings.ToUpper(input) == "Q" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("conn writing fail", err)
			conn.Close()
			break
		}
	}
}

func main() {
	conn, e := net.Dial("tcp", "127.0.0.1:8805")
	Check(e)
	defer conn.Close()

	//foreach writing via a goroutine
	go MessageSend(conn)

	//the coon foreach reading via main goroutine

	bytes := make([]byte, 1024)

	for {
		read, err := conn.Read(bytes)
		if err != nil {
			fmt.Println("您已经退出。。。 ")
			os.Exit(0)
		}
		//fmt.Println("从服务器接收到：",string(bytes[0:read]))
		fmt.Println(string(bytes[0:read]))
	}

	// clien exit,print
	//fmt.Println("client end")
}
