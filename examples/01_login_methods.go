package main

import (
	"fmt"
	"os"

	"github.com/samsaralc/hiksdk/core/auth"
)

// 两种登录方式示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 登录方式示例")
	fmt.Println("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64", // 替换为你的设备IP
		Port:     8000,           // 替换为你的端口
		Username: "admin",        // 替换为你的用户名
		Password: "password",     // 替换为你的密码
	}

	fmt.Println("\n设备连接信息:")
	fmt.Printf("  - IP地址: %s\n", cred.IP)
	fmt.Printf("  - 端口: %d\n", cred.Port)
	fmt.Printf("  - 用户名: %s\n", cred.Username)

	// ==================== 方式1: LoginV40 (推荐) ====================
	fmt.Println("\n========================================")
	fmt.Println("方式1: 使用 LoginV40（推荐）")
	fmt.Println("========================================")

	fmt.Println("\n[1] 使用LoginV40登录...")
	session1, err := auth.LoginV40(cred)
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("  1. 设备不在线或网络不可达")
		fmt.Println("  2. 用户名或密码错误")
		fmt.Println("  3. 设备端口配置错误")
	} else {
		fmt.Printf("✓ 登录成功\n")
		fmt.Printf("  登录ID: %d\n", session1.LoginID)
		fmt.Printf("  设备序列号: %s\n", session1.SerialNumber)
		fmt.Printf("  通道数: %d\n", session1.ChannelNum)

		// 登出
		fmt.Println("\n[2] 登出设备...")
		if err := auth.Logout(session1.LoginID); err != nil {
			fmt.Printf("✗ 登出失败: %v\n", err)
		}
	}

	// ==================== 方式2: LoginV30 (兼容旧设备) ====================
	fmt.Println("\n========================================")
	fmt.Println("方式2: 使用 LoginV30（兼容旧设备）")
	fmt.Println("========================================")

	fmt.Println("\n[1] 使用LoginV30登录...")
	session2, err := auth.LoginV30(cred)
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
	} else {
		fmt.Printf("✓ 登录成功\n")
		fmt.Printf("  登录ID: %d\n", session2.LoginID)
		fmt.Printf("  设备序列号: %s\n", session2.SerialNumber)
		fmt.Printf("  通道数: %d\n", session2.ChannelNum)

		// 登出
		fmt.Println("\n[2] 登出设备...")
		if err := auth.Logout(session2.LoginID); err != nil {
			fmt.Printf("✗ 登出失败: %v\n", err)
		}
	}

	// ==================== 对比说明 ====================
	fmt.Println("\n========================================")
	fmt.Println("两种登录方式对比")
	fmt.Println("========================================")

	fmt.Println("\nLoginV40():")
	fmt.Println("  ✓ 推荐使用")
	fmt.Println("  ✓ 支持更多功能")
	fmt.Println("  ✓ 性能更好")
	fmt.Println("  ✓ 设备信息更详细")

	fmt.Println("\nLoginV30():")
	fmt.Println("  ✓ 兼容旧设备")
	fmt.Println("  ✓ 简单直接")

	fmt.Println("\n💡 建议:")
	fmt.Println("  1. 优先使用 LoginV40()")
	fmt.Println("  2. 如果失败，可尝试 LoginV30()")
	fmt.Println("  3. 登录后务必调用 Logout() 释放资源")

	// 程序结束时清理SDK
	defer auth.Cleanup()
	os.Exit(0)
}
