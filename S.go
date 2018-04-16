package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

//针对客户端连接的 map映射
var CoonQuee = make(map[string]net.Conn)

//有缓存的聊天消息队列,
var MsQuee = make(chan string, 1000)
var quitChan1 = make(chan bool)

func main() {
	listener, e := net.Listen("tcp", "127.0.0.1:8805")
	myCheck(e)

	defer listener.Close()
	fmt.Println("service,,,openning")

	//该协程属于消费者
	go DoRes()
	for {

		conn, i := listener.Accept()
		myCheck(i)
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		//为每个客户端连接添加map映射
		CoonQuee[addr] = conn

		for i := range CoonQuee {
			fmt.Println("用户列表...", i)
		}
		//该协程属于生产者，
		go DoReq(conn, addr)

	}

}

func DoRes() {
	for {
		select {
		case ms := <-MsQuee:
			DoMs(ms)
			//Dom(ms)
		/*default:
		break*/
		case <-quitChan1:
			break
		}
	}
}
func Dom(s string) {
	fmt.Println(s)
}

/**
@ms 处理客户端发来的消息
for 中转私聊  和  公共聊天
*/
func DoMs(ms string) {
	//fmt.Println(ms)

	split := strings.Split(ms, "#")

	//如果数组为3，说明私聊    案例  127.0.0.1:8897#我叫xxxx#(myIp):127.0.0.1:8893
	if len(split) > 2 {

		addr := strings.Trim(split[0], " ")
		sendTime := time.Now().Format("2006-01-02 15:04:05") + "FORM " + split[2]
		chatMs := split[1] + " __ " + sendTime

		if desCon, ok := CoonQuee[addr]; ok {
			_, err := desCon.Write([]byte(chatMs))
			myCheck(err)
		}

	} else {
		//公开聊天   你好地球人#(p)127.0.0.1:8897
		for ip, coon := range CoonQuee {
			trimIP := strings.Trim(split[len(split)-1], " ")
			sendTime := time.Now().Format("2006-01-02 15:04:05")
			if trimIP == ip {
				//sendTime := time.Now().Format("2006-01-02 15:04:05")+
				chatMs := split[0] + " __ " + sendTime + " FORM MySelf "
				_, err := coon.Write([]byte(chatMs))
				myCheck(err)
				continue
			}

			chatMs := split[0] + " __ " + sendTime + " FORM " + trimIP
			_, err := coon.Write([]byte(chatMs))
			myCheck(err)
		}

	}

}

func DoReq(conn net.Conn, addr string) {
	defer func() {
		delete(CoonQuee, addr)
		conn.Close()
	}()

	tipLogin := "welcome using qq,your ip:" + addr
	_, err := conn.Write([]byte(tipLogin))
	myCheck(err)

	var ips string = ""
	for i := range CoonQuee {
		if i == addr {
			i = "wo"
		}
		ips = ips + i + "|"
	}
	ips = "all online users :" + ips
	_, err2 := conn.Write([]byte(ips))
	myCheck(err2)

	b := make([]byte, 1024)

	for {
		n, err := conn.Read(b)
		fmt.Println(n)
		if err != nil {
			break
		}
		if n != 0 {
			ms := string(b[0:n]) + "#" + addr
			MsQuee <- ms
		}

	}
}
func myCheck(e error) {
	if e != nil {
		panic(e)
	}
}
