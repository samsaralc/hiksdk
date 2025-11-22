package main

import (
	"fmt"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

// PTZ高级控制示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - PTZ高级控制示例")
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

	channel := 1

	// ==================== 云台移动控制 ====================
	fmt.Println("\n【云台移动控制】")
	demonstrateMovement(session.LoginID, channel)

	// ==================== 相机控制 ====================
	fmt.Println("\n【相机控制】")
	demonstrateCamera(session.LoginID, channel)

	// ==================== 辅助设备控制 ====================
	fmt.Println("\n【辅助设备控制】")
	demonstrateAuxiliary(session.LoginID, channel)

	fmt.Println("\n========================================")
	fmt.Println("示例完成!")
	fmt.Println("========================================")
}

// 云台移动控制演示
func demonstrateMovement(loginID int, channel int) {
	// 创建移动控制器
	movement := ptz.NewMovementController(loginID, channel)

	fmt.Println("\n[1] 基础方向移动:")

	// 向右移动2秒
	fmt.Println("  • 向右移动2秒...")
	if err := movement.Right(5, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 向上移动2秒
	fmt.Println("  • 向上移动2秒...")
	if err := movement.Up(5, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("\n[2] 组合方向移动:")

	// 右上斜向移动
	fmt.Println("  • 右上斜向移动2秒...")
	if err := movement.UpRight(4, 2*time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("\n[3] 自动扫描:")

	// 启动自动扫描
	fmt.Println("  • 启动自动扫描...")
	if err := movement.AutoScan(3); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	} else {
		time.Sleep(5 * time.Second)
		fmt.Println("  • 停止自动扫描...")
		if err := movement.StopAutoScan(); err != nil {
			fmt.Printf("    ✗ 失败: %v\n", err)
		}
	}
}

// 相机控制演示
func demonstrateCamera(loginID int, channel int) {
	// 创建相机控制器
	camera := ptz.NewCameraController(loginID, channel)

	fmt.Println("\n[1] 焦距控制:")

	// 焦距放大
	fmt.Println("  • 焦距放大（拉近）1秒...")
	if err := camera.ZoomIn(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// 焦距缩小
	fmt.Println("  • 焦距缩小（拉远）1秒...")
	if err := camera.ZoomOut(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("\n[2] 焦点控制:")

	// 焦点前调
	fmt.Println("  • 焦点前调（聚焦近处）1秒...")
	if err := camera.FocusNear(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("\n[3] 光圈控制:")

	// 光圈扩大
	fmt.Println("  • 光圈扩大（变亮）1秒...")
	if err := camera.IrisOpen(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// 光圈缩小
	fmt.Println("  • 光圈缩小（变暗）1秒...")
	if err := camera.IrisClose(1 * time.Second); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}
}

// 辅助设备控制演示
func demonstrateAuxiliary(loginID int, channel int) {
	// 创建辅助设备控制器
	aux := ptz.NewAuxiliaryController(loginID, channel)

	fmt.Println("\n[1] 灯光控制:")

	// 开启灯光
	fmt.Println("  • 开启灯光...")
	if err := aux.LightOn(); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	} else {
		time.Sleep(2 * time.Second)

		// 关闭灯光
		fmt.Println("  • 关闭灯光...")
		if err := aux.LightOff(); err != nil {
			fmt.Printf("    ✗ 失败: %v\n", err)
		}
	}

	fmt.Println("\n[2] 雨刷控制:")

	// 开启雨刷
	fmt.Println("  • 开启雨刷...")
	if err := aux.WiperOn(); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	} else {
		time.Sleep(2 * time.Second)

		// 关闭雨刷
		fmt.Println("  • 关闭雨刷...")
		if err := aux.WiperOff(); err != nil {
			fmt.Printf("    ✗ 失败: %v\n", err)
		}
	}

	fmt.Println("\n💡 说明:")
	fmt.Println("  • 辅助设备功能需要硬件支持")
	fmt.Println("  • 如果设备不支持某些功能，会返回错误码23（不支持该操作）")
}
