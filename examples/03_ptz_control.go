package main

import (
	"fmt"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

// PTZ云台控制示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - PTZ控制示例")
	fmt.Println("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64",
		Port:     8000,
		Username: "admin",
		Password: "asdf234.",
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

	// 选择通道
	channel := 1

	// ==================== 方向控制 ====================
	fmt.Println("\n[1] 方向控制示例:")

	// 创建移动控制器
	movement := ptz.NewMovementController(session.LoginID, channel)

	// 向右转动2秒
	fmt.Println("  向右转动2秒...")
	if err := movement.Right(7, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 向上转动2秒
	fmt.Println("  向上转动2秒...")
	if err := movement.Up(7, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 右上斜向移动
	fmt.Println("  右上斜向移动2秒...")
	if err := movement.UpRight(5, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// ==================== 变焦控制 ====================
	fmt.Println("\n[2] 变焦控制示例:")

	// 创建相机控制器
	camera := ptz.NewCameraController(session.LoginID, channel)

	// 放大
	fmt.Println("  焦距放大（拉近）1秒...")
	if err := camera.ZoomIn(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// 缩小
	fmt.Println("  焦距缩小（拉远）1秒...")
	if err := camera.ZoomOut(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// ==================== 预置点操作 ====================
	fmt.Println("\n[3] 预置点操作:")

	// 创建预置点控制器
	preset := ptz.NewPresetManager(session.LoginID, channel)

	// 设置预置点1
	fmt.Println("  设置预置点1...")
	if err := preset.SetPreset(1); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 移动到其他位置
	fmt.Println("  移动到其他位置...")
	if err := movement.Left(4, 3*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 转到预置点1
	fmt.Println("  转到预置点1...")
	if err := preset.GotoPreset(1); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}
	time.Sleep(3 * time.Second)

	fmt.Println("\n示例完成!")
}
