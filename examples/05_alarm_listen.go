package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/samsaralc/hiksdk/core/alarm"
	"github.com/samsaralc/hiksdk/core/auth"
)

// 报警监听示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 报警监听示例")
	fmt.Println("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64",
		Port:     8000,
		Username: "admin",
		Password: "password",
	}

	// 登录设备
	session, err := auth.LoginV40(cred)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Printf("登录成功 (ID: %d)\n", session.LoginID)
	defer auth.Logout(session.LoginID)
	defer auth.Cleanup()

	// 创建报警监听器
	fmt.Println("\n创建报警监听器...")
	listener := alarm.NewAlarmListener(session.LoginID)

	// 启动报警监听
	fmt.Println("启动报警监听...")
	if err := listener.Start(); err != nil {
		fmt.Printf("启动监听失败: %v\n", err)
		return
	}
	defer listener.Stop()

	fmt.Println("\n监听中... 按 Ctrl+C 退出")
	fmt.Println("等待接收报警消息（移动侦测、遮挡报警等）")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n停止监听...")
	fmt.Println("示例完成!")
}
