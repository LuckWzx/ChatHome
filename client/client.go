package main

import (
	"bufio"
	"chatroom/utils"
	"fmt"
	"net"
	"os"
)

func main() {
	// 获取用户名
	var flag = true
	var username string
	for flag {
		fmt.Print("请输入您的网名: ")
		username = utils.ReadLine()
		if username == "" {
			fmt.Println("未输入用户名，请从新输入！")
		} else {
			flag = false
		}
	}

	// 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 发送用户名到服务器
	_, err = conn.Write([]byte(username + "\n"))
	if err != nil {
		fmt.Println("发送用户名失败:", err)
		os.Exit(1)
	}

	// 启动接收消息的协程
	go receiveMessages(conn)

	// 发送消息
	fmt.Println("连接到服务器！输入 'exit' 退出，输入check可以显示展示在线用户")
	fmt.Println("输入private进行用户私聊")
	for {
		fmt.Print("> ")
		message := utils.ReadLine()

		if message == "exit" {
			break
		}

		// 发送消息
		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("发送消息失败:", err)
			break
		}
	}

	// 发送离开消息
	conn.Write([]byte("exit\n"))
	fmt.Println("你已退出聊天室")
}

// receiveMessages 从服务器接收消息并打印
func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("与服务器连接断开")
			os.Exit(0)
		}
		fmt.Print(utils.TrimNewLine(message) + "\n")
	}
}
