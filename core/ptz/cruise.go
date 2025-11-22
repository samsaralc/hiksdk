package ptz

/*
#cgo CFLAGS: -I../../include
#cgo CFLAGS: -I..

// Linux 平台的链接配置
#cgo linux LDFLAGS: -L../../lib/Linux -lhcnetsdk -lhpr -lHCCore

// Windows 平台的链接配置
#cgo windows LDFLAGS: -L../../lib/Windows -lHCNetSDK -lHCCore

#include <stdio.h>
#include <stdlib.h>
#include "../hiksdk_wrapper.h"
*/
import "C"
import (
	"fmt"
	"log"
)

// 巡航命令常量（来自官方文档表 5.12）
const (
	// FILL_PRE_SEQ 将预置点加入巡航序列
	FILL_PRE_SEQ = 30
	// SET_SEQ_DWELL 设置巡航点停顿时间
	SET_SEQ_DWELL = 31
	// SET_SEQ_SPEED 设置巡航速度
	SET_SEQ_SPEED = 32
	// CLE_PRE_SEQ 将预置点从巡航序列中删除
	CLE_PRE_SEQ = 33
	// RUN_SEQ 开始巡航
	RUN_SEQ = 37
	// STOP_SEQ 停止巡航
	STOP_SEQ = 38
	// DEL_SEQ 删除巡航路径
	DEL_SEQ = 43
)

// 巡航参数限制（来自官方文档）
const (
	MaxCruiseRoutes = 32  // 最多支持32条路径
	MaxCruisePoints = 32  // 最多支持32个点
	MaxPresetID     = 255 // 预置点最大值
	MaxDwellTime    = 255 // 停顿时间最大值（秒）
	MaxCruiseSpeed  = 40  // 巡航速度最大值
)

// CruiseManager 巡航控制器
// 封装了云台巡航的所有操作，提供简化的API
type CruiseManager struct {
	userID  int // 登录句柄（NET_DVR_Login_V40 的返回值）
	channel int // 通道号
}

// NewCruiseManager 创建巡航控制器
// 参数：
//   - userID: 设备登录ID (dev.GetLoginID())
//   - channel: 通道号
//
// 返回：
//   - *CruiseManager: 巡航控制器实例
func NewCruiseManager(userID int, channel int) *CruiseManager {
	return &CruiseManager{
		userID:  userID,
		channel: channel,
	}
}

// AddPresetToCruise 将预置点加入巡航路径
// 对应官方命令：FILL_PRE_SEQ
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//   - pointIndex: 巡航点编号（1-32）
//   - presetID: 预置点编号（1-255）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) AddPresetToCruise(routeIndex, pointIndex, presetID int) error {
	// 参数验证
	if err := c.validateRoutePoint(routeIndex, pointIndex); err != nil {
		return err
	}
	if presetID < 1 || presetID > MaxPresetID {
		return fmt.Errorf("预置点编号超出范围：%d（有效范围：1-%d）", presetID, MaxPresetID)
	}

	if err := c.control(FILL_PRE_SEQ, routeIndex, pointIndex, presetID); err != nil {
		return fmt.Errorf("添加预置点%d到巡航路径%d点%d失败: %w", presetID, routeIndex, pointIndex, err)
	}

	log.Printf("✓ 添加预置点%d到巡航路径%d点%d成功", presetID, routeIndex, pointIndex)
	return nil
}

// RemovePresetFromCruise 从巡航路径中删除预置点
// 对应官方命令：CLE_PRE_SEQ
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//   - pointIndex: 巡航点编号（1-32）
//   - presetID: 要删除的预置点编号（1-255）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) RemovePresetFromCruise(routeIndex, pointIndex, presetID int) error {
	// 参数验证
	if err := c.validateRoutePoint(routeIndex, pointIndex); err != nil {
		return err
	}
	if presetID < 1 || presetID > MaxPresetID {
		return fmt.Errorf("预置点编号超出范围：%d（有效范围：1-%d）", presetID, MaxPresetID)
	}

	if err := c.control(CLE_PRE_SEQ, routeIndex, pointIndex, presetID); err != nil {
		return fmt.Errorf("从巡航路径%d点%d删除预置点%d失败: %w", routeIndex, pointIndex, presetID, err)
	}

	log.Printf("✓ 从巡航路径%d点%d删除预置点%d成功", routeIndex, pointIndex, presetID)
	return nil
}

// SetCruiseSpeed 设置巡航点速度
// 对应官方命令：SET_SEQ_SPEED
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//   - pointIndex: 巡航点编号（1-32）
//   - speed: 速度（1-40，根据官方文档）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) SetCruiseSpeed(routeIndex, pointIndex, speed int) error {
	// 参数验证
	if err := c.validateRoutePoint(routeIndex, pointIndex); err != nil {
		return err
	}
	if speed < 1 || speed > MaxCruiseSpeed {
		return fmt.Errorf("速度超出范围：%d（有效范围：1-%d）", speed, MaxCruiseSpeed)
	}

	if err := c.control(SET_SEQ_SPEED, routeIndex, pointIndex, speed); err != nil {
		return fmt.Errorf("设置巡航路径%d点%d速度为%d失败: %w", routeIndex, pointIndex, speed, err)
	}

	log.Printf("✓ 设置巡航路径%d点%d速度为%d成功", routeIndex, pointIndex, speed)
	return nil
}

// SetCruiseDwellTime 设置巡航点停顿时间
// 对应官方命令：SET_SEQ_DWELL
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//   - pointIndex: 巡航点编号（1-32）
//   - dwellTime: 停顿时间（秒，1-255，根据官方文档）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) SetCruiseDwellTime(routeIndex, pointIndex, dwellTime int) error {
	// 参数验证
	if err := c.validateRoutePoint(routeIndex, pointIndex); err != nil {
		return err
	}
	if dwellTime < 1 || dwellTime > MaxDwellTime {
		return fmt.Errorf("停顿时间超出范围：%d秒（有效范围：1-%d秒）", dwellTime, MaxDwellTime)
	}

	if err := c.control(SET_SEQ_DWELL, routeIndex, pointIndex, dwellTime); err != nil {
		return fmt.Errorf("设置巡航路径%d点%d停顿时间为%d秒失败: %w", routeIndex, pointIndex, dwellTime, err)
	}

	log.Printf("✓ 设置巡航路径%d点%d停顿时间为%d秒成功", routeIndex, pointIndex, dwellTime)
	return nil
}

// StartCruise 开始巡航
// 对应官方命令：RUN_SEQ
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) StartCruise(routeIndex int) error {
	// 参数验证
	if routeIndex < 1 || routeIndex > MaxCruiseRoutes {
		return fmt.Errorf("巡航路径编号超出范围：%d（有效范围：1-%d）", routeIndex, MaxCruiseRoutes)
	}

	// 开始巡航时，Point 和 Input 设为0
	if err := c.control(RUN_SEQ, routeIndex, 0, 0); err != nil {
		return fmt.Errorf("开始巡航路径%d失败: %w", routeIndex, err)
	}

	log.Printf("✓ 开始巡航路径%d", routeIndex)
	return nil
}

// StopCruise 停止巡航
// 对应官方命令：STOP_SEQ
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) StopCruise(routeIndex int) error {
	// 参数验证
	if routeIndex < 1 || routeIndex > MaxCruiseRoutes {
		return fmt.Errorf("巡航路径编号超出范围：%d（有效范围：1-%d）", routeIndex, MaxCruiseRoutes)
	}

	if err := c.control(STOP_SEQ, routeIndex, 0, 0); err != nil {
		return fmt.Errorf("停止巡航路径%d失败: %w", routeIndex, err)
	}

	log.Printf("✓ 停止巡航路径%d", routeIndex)
	return nil
}

// DeleteCruiseRoute 删除巡航路径
// 对应官方命令：DEL_SEQ
// 参数：
//   - routeIndex: 巡航路径编号（1-32）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (c *CruiseManager) DeleteCruiseRoute(routeIndex int) error {
	// 参数验证
	if routeIndex < 1 || routeIndex > MaxCruiseRoutes {
		return fmt.Errorf("巡航路径编号超出范围：%d（有效范围：1-%d）", routeIndex, MaxCruiseRoutes)
	}

	if err := c.control(DEL_SEQ, routeIndex, 0, 0); err != nil {
		return fmt.Errorf("删除巡航路径%d失败: %w", routeIndex, err)
	}

	log.Printf("✓ 删除巡航路径%d成功", routeIndex)
	return nil
}

// validateRoutePoint 验证路径和点的编号范围
func (c *CruiseManager) validateRoutePoint(routeIndex, pointIndex int) error {
	if routeIndex < 1 || routeIndex > MaxCruiseRoutes {
		return fmt.Errorf("巡航路径编号超出范围：%d（有效范围：1-%d）", routeIndex, MaxCruiseRoutes)
	}
	if pointIndex < 1 || pointIndex > MaxCruisePoints {
		return fmt.Errorf("巡航点编号超出范围：%d（有效范围：1-%d）", pointIndex, MaxCruisePoints)
	}
	return nil
}

// control 内部通用控制函数
// 直接调用 NET_DVR_PTZCruise_Other（推荐，不需要预览）
func (c *CruiseManager) control(cmd, route, point, input int) error {
	if c.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", c.userID)
	}

	// 调用 C 接口
	ret := C.NET_DVR_PTZCruise_Other(
		C.LONG(c.userID),
		C.LONG(c.channel),
		C.DWORD(cmd),
		C.BYTE(route),
		C.BYTE(point),
		C.WORD(input),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("巡航操作失败 [通道:%d Cmd:%d Route:%d Point:%d Input:%d 错误码:%d]",
			c.channel, cmd, route, point, input, errCode)
	}

	return nil
}

// GetCommandName 获取巡航命令的名称（用于调试）
func GetCommandName(cmd int) string {
	names := map[int]string{
		FILL_PRE_SEQ:  "将预置点加入巡航序列",
		SET_SEQ_DWELL: "设置巡航点停顿时间",
		SET_SEQ_SPEED: "设置巡航速度",
		CLE_PRE_SEQ:   "将预置点从巡航序列中删除",
		RUN_SEQ:       "开始巡航",
		STOP_SEQ:      "停止巡航",
		DEL_SEQ:       "删除巡航路径",
	}
	if name, ok := names[cmd]; ok {
		return name
	}
	return fmt.Sprintf("未知命令(%d)", cmd)
}
