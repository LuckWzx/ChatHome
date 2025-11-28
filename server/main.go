package main

import (
	"bufio"
	"chatroom/common"
	"fmt"
	"net"
	"strings"
	"sync"
)

func main() {
	// 创建服务器
	server := &common.Server{
		//Clients: make(map[string]*common.Client),
		Clients: sync.Map{},
	}

	// 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("服务器启动失败:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务器已启动，等待客户端连接...")

	// 接收客户端连接
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			fmt.Println("接受连接失败:", err2)
			continue
		}

		// 读取用户名
		reader := bufio.NewReader(conn)
		username, err3 := reader.ReadString('\n')
		if err3 != nil {
			fmt.Println("读取用户名失败:", err3)
			conn.Close()
			continue
		}

		username = strings.TrimSpace(username)

		// 添加客户端
		if !server.AddClient(conn, username) {
			// 用户名重复，连接已在addClient中关闭
			continue
		}
		//var count int = 0
		//server.Clients.Range(func(key, value interface{}) bool {
		//	count++
		//	return true
		//})

		// 在clients map中查找客户端
		//server.Mutex.Lock()
		//client, exists := server.Clients[username]
		//server.Mutex.Unlock()
		value, ok := server.Clients.Load(username)
		if !ok {
			continue
		}

		//if !exists {
		//	conn.Close()
		//	continue
		//}
		client := value.(*common.Client)

		// 启动处理客户端的协程
		go server.ProcessClient(client)
	}
}

// Client 代表一个客户端连接
//type Client struct {
//	conn net.Conn
//	name string
//}

//Server 代表聊天服务器
//type Server struct {
//	clients map[string]*common.Client
//	mutex   sync.Mutex
//}

//// broadcast 在调用前必须持有锁
//func (s *Server) broadcast(message string, excludeClient *common.Client) {
//	// 注意：这个方法假设调用者已经持有 s.mutex 锁
//	for _, client := range s.clients {
//		if excludeClient != nil && client == excludeClient {
//			continue // 跳过发送消息的客户端
//		}
//
//		_, err := client.Conn.Write([]byte(message + "\n"))
//		if err != nil {
//			fmt.Printf("向客户端 %s 广播消息失败: %v\n", client.Name, err)
//			// 不在这里移除客户端，由上层处理
//		}
//	}
//}
//
//func (s *Server) addClient(conn net.Conn, name string) bool {
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//
//	// 检查用户名是否已存在
//	if _, exists := s.clients[name]; exists {
//		conn.Write([]byte("用户名已存在，请重新输入\n"))
//		conn.Close()
//		return false
//	}
//
//	// 创建新客户端
//	client := &common.Client{
//		Conn: conn,
//		Name: name,
//	}
//
//	// 添加到客户端列表
//	s.clients[name] = client
//
//	// 发送欢迎消息
//	welcomeMsg := fmt.Sprintf("欢迎 %s 加入聊天室！当前在线人数: %d\n", name, len(s.clients))
//	conn.Write([]byte(welcomeMsg))
//
//	// 广播新用户加入，使用同一个锁
//	s.broadcast(fmt.Sprintf("%s 加入了聊天室", name), nil)
//	return true
//}
//
//func (s *Server) removeClient(client *common.Client) {
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//
//	// 从客户端列表中移除
//	if _, exists := s.clients[client.Name]; exists {
//		delete(s.clients, client.Name)
//		client.Conn.Close()
//	}
//}
//
//func (s *Server) processClient(client *common.Client) {
//	defer func() {
//		// 先广播离开消息
//		s.mutex.Lock()
//		s.broadcast(fmt.Sprintf("%s 离开了聊天室", client.Name), nil)
//		s.mutex.Unlock()
//
//		// 然后移除客户端
//		s.removeClient(client)
//	}()
//
//	reader := bufio.NewReader(client.Conn)
//	for {
//		message, err := reader.ReadString('\n')
//		if err != nil {
//			// 正常断开连接
//			return
//		}
//
//		content := strings.TrimSpace(message)
//		if content == "exit" {
//			return
//		}
//
//		// 广播消息前获取锁
//		s.mutex.Lock()
//		s.broadcast(fmt.Sprintf("%s: %s", client.Name, content), client)
//		s.mutex.Unlock()
//	}
//}
