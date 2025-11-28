package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//消息处理

// ReadLine 从标准输入读取一行
func ReadLine() string {

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("用户退出")
		//fmt.Println("Error reading input:")
		return ""
	}
	return strings.TrimSpace(line)
}

// TrimNewLine 去掉字符串末尾的换行符
func TrimNewLine(s string) string {
	return strings.TrimRight(s, "\n")
}
