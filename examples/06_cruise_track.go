package main

import (
	"fmt"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

// 巡航和轨迹控制示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 巡航与轨迹示例")
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

	channel := 1 // 通道1

	// ==================== 巡航示例 ====================
	fmt.Println("\n【巡航功能示例】")
	demonstrateCruise(session.LoginID, channel)

	// ==================== 轨迹示例 ====================
	fmt.Println("\n【轨迹功能示例】")
	demonstrateTrack(session.LoginID, channel)

	fmt.Println("\n示例完成!")
}

// 巡航功能演示
func demonstrateCruise(loginID int, channel int) {
	// 创建巡航控制器
	cruise := ptz.NewCruiseManager(loginID, channel)

	// 定义巡航路径
	routeIndex := 1 // 使用路径1

	fmt.Printf("\n[1] 配置巡航路径 %d:\n", routeIndex)

	// 添加预置点到巡航路径
	fmt.Println("  • 添加预置点1到路径1点1...")
	if err := cruise.AddPresetToCruise(routeIndex, 1, 1); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("  • 添加预置点2到路径1点2...")
	if err := cruise.AddPresetToCruise(routeIndex, 2, 2); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("  • 添加预置点3到路径1点3...")
	if err := cruise.AddPresetToCruise(routeIndex, 3, 3); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 设置每个点的停顿时间
	fmt.Println("\n[2] 设置停顿时间:")
	fmt.Println("  • 点1停顿5秒...")
	if err := cruise.SetCruiseDwellTime(routeIndex, 1, 5); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("  • 点2停顿3秒...")
	if err := cruise.SetCruiseDwellTime(routeIndex, 2, 3); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("  • 点3停顿5秒...")
	if err := cruise.SetCruiseDwellTime(routeIndex, 3, 5); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 设置巡航速度
	fmt.Println("\n[3] 设置巡航速度:")
	fmt.Println("  • 点1速度设为20...")
	if err := cruise.SetCruiseSpeed(routeIndex, 1, 20); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	fmt.Println("  • 点2速度设为25...")
	if err := cruise.SetCruiseSpeed(routeIndex, 2, 25); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}

	// 开始巡航
	fmt.Println("\n[4] 开始巡航路径1...")
	if err := cruise.StartCruise(routeIndex); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	} else {
		fmt.Println("    ✓ 巡航已启动，云台将按路径自动移动")
		time.Sleep(20 * time.Second) // 运行20秒
	}

	// 停止巡航
	fmt.Println("\n[5] 停止巡航...")
	if err := cruise.StopCruise(routeIndex); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
	}
}

// 轨迹功能演示
func demonstrateTrack(loginID int, channel int) {
	// 创建轨迹控制器
	track := ptz.NewTrackManager(loginID, channel)

	// 开始记录轨迹
	fmt.Println("\n[1] 开始记录轨迹...")
	if err := track.StartRecordTrack(); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
		return
	}
	fmt.Println("    ✓ 正在记录，请手动控制云台移动...")

	// 模拟云台移动（实际使用时这里应该是真实的云台操作）
	fmt.Println("\n[2] 模拟云台移动（记录中）...")
	time.Sleep(5 * time.Second)

	// 停止记录
	fmt.Println("\n[3] 停止记录轨迹...")
	if err := track.StopRecordTrack(); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
		return
	}
	fmt.Println("    ✓ 轨迹已保存")

	// 等待一会
	time.Sleep(2 * time.Second)

	// 执行记录的轨迹
	fmt.Println("\n[4] 执行记录的轨迹...")
	if err := track.RunTrack(); err != nil {
		fmt.Printf("    ✗ 失败: %v\n", err)
		return
	}
	fmt.Println("    ✓ 轨迹执行中，云台将按记录的路径移动")
	time.Sleep(10 * time.Second)

	fmt.Println("\n💡 说明:")
	fmt.Println("  • 巡航：基于预置点的自动移动路径")
	fmt.Println("  • 轨迹：录制云台的移动轨迹并回放")
	fmt.Println("  • 巡航路径最多32条，每条最多32个点")
	fmt.Println("  • 轨迹通常用于复杂的移动模式")
}
