package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	//fmt.Println("链接建立成功")

	user := NewUser(conn, this)

	user.Online()

	//接受客户端发送的消息
	go func() {
		// 创建一个4096字节的缓冲区，用于存储从客户端接收的数据
		buf := make([]byte, 4096)
		for {
			// 从conn（一个网络连接对象）中读取数据到buf中，返回读取的字节数和可能发生的错误
			// 这一句时阻塞的
			n, err := conn.Read(buf)

			// 因为上一句时阻塞的，客户端 ctrl+c 关闭的时候会发一个 eof 以及 n==0 的消息过来
			if n == 0 {
				user.Offline()
				return
			}

			// 'err != nil' 用来检查错误对象是否为空，如果不为空，表示在之前的操作中出现了一些错误
			// 'err != io.EOF' 是用来检查错误是否为 'io.EOF'，也就是文件或数据流已经结束的错误。
			// 只有当错误对象不为空且错误类型不是 io.EOF 时，这个条件才会为真。
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			//用户针对msg进行消息处理
			user.DoMessage(msg)
		}
	}()

	//当前handler阻塞
	select {}
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
