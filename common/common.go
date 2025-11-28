package common

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	Conn net.Conn
	Name string
}
type Server struct {
	//Clients map[string]*Client
	Clients sync.Map
	//Mutex   sync.Mutex
}

// broadcast 在调用前必须持有锁

func (s *Server) Broadcast(message string, excludeClient *Client) {
	// 注意：这个方法假设调用者已经持有 s.mutex 锁
	//for _, client := range s.Clients {
	//	if excludeClient != nil && client == excludeClient {
	//		continue // 跳过发送消息的客户端
	//	}
	//
	//	_, err := client.Conn.Write([]byte(message + "\n"))
	//	if err != nil {
	//		fmt.Printf("客户端 %s离开聊天室， 发送消息失败。\n", client.Name)
	//		// 不在这里移除客户端，由上层处理
	//	}
	//}
	s.Clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)

		if excludeClient != nil && client == excludeClient {
			// 如果是用户本身，就跳过消息广播
			return true
		}

		_, err := client.Conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Printf("客户端 %s离开聊天室， 发送消息失败。\n", client.Name)
		}
		return true
	})
}

// 在 common 包的 Server 结构体中添加这个方法  获取到在线人数

func (s *Server) GetOnlineCount() int {
	count := 0
	s.Clients.Range(func(key, value interface{}) bool {
		count++
		return true // 继续遍历
	})
	return count
}

//添加用户

func (s *Server) AddClient(conn net.Conn, name string) bool {
	// 检查用户是否已存在
	if _, exist := s.Clients.Load(name); exist {
		conn.Write([]byte("网名已存在!!!!!!\n"))
		conn.Close()
		return false
	}

	fmt.Printf("%s进入聊天室。\n", name)
	//创建新的客户端
	client := &Client{
		Conn: conn,
		Name: name,
	}

	// 存储客户列表
	s.Clients.Store(name, client)
	conn.Write([]byte(fmt.Sprintf("欢迎 %s 加入聊天室！当前在线人数: %d\n", name, s.GetOnlineCount())))

	// 广播新用户加入
	s.Broadcast(fmt.Sprintf("%s 加入了聊天室", name), client)
	return true
	//s.Mutex.Lock()
	//defer s.Mutex.Unlock()
	//
	////_, exists := s.Clients[name]
	//// 检查用户名是否已存在
	//if _, exists := s.Clients[name]; exists {
	//	conn.Write([]byte("用户名已存在!!!!!!\n"))
	//	conn.Close()
	//	return false
	//}
	//fmt.Printf("%s进入聊天室。\n", name)
	//// 创建新客户端
	//client := &Client{
	//	Conn: conn,
	//	Name: name,
	//}
	//
	//// 添加到客户端列表
	//s.Clients[name] = client
	//
	//// 发送欢迎消息
	//welcomeMsg := fmt.Sprintf("欢迎 %s 加入聊天室！当前在线人数: %d\n", name, len(s.Clients))
	//conn.Write([]byte(welcomeMsg))
	//
	//// 广播新用户加入，使用同一个锁
	//s.Broadcast(fmt.Sprintf("%s 加入了聊天室", name), client)
	//return true
}

//移除用户

func (s *Server) RemoveClient(client *Client) {
	s.Clients.Delete(client.Name)
	client.Conn.Close()
	client.Conn.Write([]byte(fmt.Sprintf("当前在线人数: %d\n", s.GetOnlineCount())))
	//s.Mutex.Lock()
	//defer s.Mutex.Unlock()
	//
	//// 从客户端列表中移除
	//if _, exists := s.Clients[client.Name]; exists {
	//	delete(s.Clients, client.Name)
	//	client.Conn.Close()
	//}
}

func (s *Server) ProcessClient(client *Client) {
	defer func() {
		// 广播离开消息
		s.Broadcast(fmt.Sprintf("%s 离开了聊天室\n", client.Name), nil)
		// 移除客户端
		s.RemoveClient(client)
	}()

	reader := bufio.NewReader(client.Conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			// 正常断开连接
			return
		}
		content := strings.TrimSpace(message)
		// 如果输入 exit 则退出程序
		if content == "exit" {
			fmt.Printf("%s离开聊天室", client.Name)
			return
		}
		if content == "check" {
			//查看在线用户
			s.ViewOnlineClient(client)
			//跳过本次循环
			continue
		}
		if content == "private" {
			//进入私聊模式
			s.PrivateChat(client)
			client.Conn.Write([]byte("已返回广播模式！\n"))
			continue
		}
		// 广播消息
		s.Broadcast(fmt.Sprintf("%s: %s", client.Name, content), client)
	}

	//defer func() {
	//	// 先广播离开消息
	//	s.Mutex.Lock()
	//	s.Broadcast(fmt.Sprintf("%s 离开了聊天室", client.Name), nil)
	//	s.Mutex.Unlock()
	//
	//	// 然后移除客户端
	//	s.RemoveClient(client)
	//}()
	//
	//reader := bufio.NewReader(client.Conn)
	//for {
	//	message, err := reader.ReadString('\n')
	//	if err != nil {
	//		// 正常断开连接
	//		return
	//	}
	//
	//	content := strings.TrimSpace(message)
	//	if content == "exit" {
	//		fmt.Printf("%s离开聊天室。", client.Name)
	//		return
	//	}
	//
	//	// 广播消息前获取锁
	//	s.Mutex.Lock()
	//	s.Broadcast(fmt.Sprintf("%s: %s", client.Name, content), client)
	//	s.Mutex.Unlock()
	//}
}

//查看在线用户

func (s *Server) ViewOnlineClient(client *Client) {
	//s.Clients
	var keys string

	s.Clients.Range(func(key, value interface{}) bool {
		// 因为知道 key 是 string，所以直接断言
		if k, ok := key.(string); ok {
			//keys = append(keys, k)
			keys += k + " "
		} else {
			// 防御性编程
			fmt.Printf("警告: key 不是 string 类型: %v (%T)\n", key, key)
		}
		return true // 继续遍历下一个
	})
	client.Conn.Write([]byte("在线用户名展示" + keys + "\n"))

}

//私聊模式

func (s *Server) PrivateChat(client *Client) {
	conn := client.Conn
	reader := bufio.NewReader(client.Conn)

	var targetUsername string
	var targetClient *Client

	// 选择聊天对象
	for {
		conn.Write([]byte("请输入聊天对象网名：\n"))
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		targetUsername = strings.TrimSpace(message)
		if targetUsername == "" {
			conn.Write([]byte("未输入用户名，请重新输入：\n"))
			continue
		}

		// 检查用户是否存在
		found := false
		s.Clients.Range(func(key, value interface{}) bool {
			if k, ok := key.(string); ok && k == targetUsername {
				found = true
				if v, ok := value.(*Client); ok {
					targetClient = v
				}
				return false
			}
			return true
		})

		if !found {
			conn.Write([]byte("用户不在线或不存在，请重新输入：\n"))
			continue
		}

		// 不能和自己私聊
		if targetUsername == client.Name {
			conn.Write([]byte("不能给自己发私聊消息！\n"))
			continue
		}

		// 找到有效目标，退出循环
		break
	}

	conn.Write([]byte("----------------进入私聊模式---------------------\n"))
	conn.Write([]byte(fmt.Sprintf("正在与 %s 私聊 (输入 'exits' 退出私聊)\n", targetUsername)))

	// 私聊消息循环
	for {
		conn.Write([]byte(client.Name + "> "))
		message, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte("读取信息失败，返回公聊模式...\n"))
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue // 忽略空消息
		}

		if message == "exits" {
			conn.Write([]byte("退出私聊，返回公聊模式...\n"))
			return
		}

		// 发送消息给目标用户
		if targetClient != nil {
			_, err = targetClient.Conn.Write([]byte(fmt.Sprintf("[私聊]%s: %s\n", client.Name, message)))
			if err != nil {
				conn.Write([]byte("消息发送失败，对方可能已离线，返回公聊模式...\n"))
				return
			}
		}
	}
}
