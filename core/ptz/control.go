package ptz

/*
#include <stdio.h>
#include <stdlib.h>
#include "../hiksdk_wrapper.h"
*/
import "C"
import (
	"fmt"
	"log"
	"time"
)

// ==================== 云台移动命令常量（来自官方文档表 5.10）====================
const (
	// ========== 基本云台移动 ==========
	// TILT_UP 云台上仰
	TILT_UP = 21
	// TILT_DOWN 云台下俯
	TILT_DOWN = 22
	// PAN_LEFT 云台左转
	PAN_LEFT = 23
	// PAN_RIGHT 云台右转
	PAN_RIGHT = 24

	// ========== 组合移动 ==========
	// UP_LEFT 云台上仰和左转
	UP_LEFT = 25
	// UP_RIGHT 云台上仰和右转
	UP_RIGHT = 26
	// DOWN_LEFT 云台下俯和左转
	DOWN_LEFT = 27
	// DOWN_RIGHT 云台下俯和右转
	DOWN_RIGHT = 28

	// ========== 自动扫描 ==========
	// PAN_AUTO 云台左右自动扫描
	PAN_AUTO = 29
)

// ==================== 相机控制命令常量（来自官方文档表 5.10）====================
const (
	// ========== 焦距控制 ==========
	// ZOOM_IN 焦距变大(倍率变大)
	ZOOM_IN = 11
	// ZOOM_OUT 焦距变小(倍率变小)
	ZOOM_OUT = 12

	// ========== 焦点控制 ==========
	// FOCUS_NEAR 焦点前调
	FOCUS_NEAR = 13
	// FOCUS_FAR 焦点后调
	FOCUS_FAR = 14

	// ========== 光圈控制 ==========
	// IRIS_OPEN 光圈扩大
	IRIS_OPEN = 15
	// IRIS_CLOSE 光圈缩小
	IRIS_CLOSE = 16
)

// ==================== 辅助设备命令常量（来自官方文档表 5.10）====================
const (
	// LIGHT_PWRON 接通灯光电源
	LIGHT_PWRON = 2
	// WIPER_PWRON 接通雨刷开关
	WIPER_PWRON = 3
	// FAN_PWRON 接通风扇开关
	FAN_PWRON = 4
	// HEATER_PWRON 接通加热器开关
	HEATER_PWRON = 5
	// AUX_PWRON1 接通辅助设备开关1
	AUX_PWRON1 = 6
	// AUX_PWRON2 接通辅助设备开关2
	AUX_PWRON2 = 7
)

// 控制参数常量（来自官方文档）
const (
	// PTZ_START 开始动作
	PTZ_START = 0
	// PTZ_STOP 停止动作
	PTZ_STOP = 1

	// 速度范围：1-7（根据官方文档 5.6.3 和 5.6.4）
	MinSpeed     = 1
	MaxSpeed     = 7
	DefaultSpeed = 4
)

// ==================== MovementController 云台移动控制器 ====================

// MovementController 云台移动控制器
// 封装云台的方向移动操作
type MovementController struct {
	userID  int // 登录句柄
	channel int // 通道号
}

// NewMovementController 创建云台移动控制器
func NewMovementController(userID int, channel int) *MovementController {
	return &MovementController{
		userID:  userID,
		channel: channel,
	}
}

// Up 云台上仰
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) Up(speed int, duration time.Duration) error {
	return m.move(TILT_UP, speed, duration)
}

// Down 云台下俯
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) Down(speed int, duration time.Duration) error {
	return m.move(TILT_DOWN, speed, duration)
}

// Left 云台左转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) Left(speed int, duration time.Duration) error {
	return m.move(PAN_LEFT, speed, duration)
}

// Right 云台右转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) Right(speed int, duration time.Duration) error {
	return m.move(PAN_RIGHT, speed, duration)
}

// UpLeft 云台上仰并左转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) UpLeft(speed int, duration time.Duration) error {
	return m.move(UP_LEFT, speed, duration)
}

// UpRight 云台上仰并右转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) UpRight(speed int, duration time.Duration) error {
	return m.move(UP_RIGHT, speed, duration)
}

// DownLeft 云台下俯并左转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) DownLeft(speed int, duration time.Duration) error {
	return m.move(DOWN_LEFT, speed, duration)
}

// DownRight 云台下俯并右转
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (m *MovementController) DownRight(speed int, duration time.Duration) error {
	return m.move(DOWN_RIGHT, speed, duration)
}

// AutoScan 云台左右自动扫描
// 参数：
//   - speed: 速度（1-7）
func (m *MovementController) AutoScan(speed int) error {
	if err := m.validateSpeed(speed); err != nil {
		return err
	}

	// 开始自动扫描
	if err := m.controlWithSpeed(PAN_AUTO, PTZ_START, speed); err != nil {
		return fmt.Errorf("启动自动扫描失败: %w", err)
	}

	log.Printf("✓ 启动自动扫描（通道%d，速度%d）", m.channel, speed)
	return nil
}

// StopAutoScan 停止自动扫描
func (m *MovementController) StopAutoScan() error {
	if err := m.controlWithSpeed(PAN_AUTO, PTZ_STOP, DefaultSpeed); err != nil {
		return fmt.Errorf("停止自动扫描失败: %w", err)
	}

	log.Printf("✓ 停止自动扫描（通道%d）", m.channel)
	return nil
}

// move 内部移动函数，自动处理开始和停止
func (m *MovementController) move(cmd, speed int, duration time.Duration) error {
	// 验证速度
	if err := m.validateSpeed(speed); err != nil {
		return err
	}

	// 开始移动
	if err := m.controlWithSpeed(cmd, PTZ_START, speed); err != nil {
		return err
	}

	// 等待指定时间
	time.Sleep(duration)

	// 停止移动
	if err := m.controlWithSpeed(cmd, PTZ_STOP, speed); err != nil {
		return err
	}

	return nil
}

// ==================== CameraController 相机控制器 ====================

// CameraController 相机控制器
// 封装焦距、焦点、光圈等相机参数控制
type CameraController struct {
	userID  int // 登录句柄
	channel int // 通道号
}

// NewCameraController 创建相机控制器
func NewCameraController(userID int, channel int) *CameraController {
	return &CameraController{
		userID:  userID,
		channel: channel,
	}
}

// ZoomIn 焦距放大（拉近）
// 参数：
//   - duration: 持续时间
func (c *CameraController) ZoomIn(duration time.Duration) error {
	return c.adjust(ZOOM_IN, duration, "焦距放大")
}

// ZoomOut 焦距缩小（拉远）
// 参数：
//   - duration: 持续时间
func (c *CameraController) ZoomOut(duration time.Duration) error {
	return c.adjust(ZOOM_OUT, duration, "焦距缩小")
}

// FocusNear 焦点前调（聚焦近处）
// 参数：
//   - duration: 持续时间
func (c *CameraController) FocusNear(duration time.Duration) error {
	return c.adjust(FOCUS_NEAR, duration, "焦点前调")
}

// FocusFar 焦点后调（聚焦远处）
// 参数：
//   - duration: 持续时间
func (c *CameraController) FocusFar(duration time.Duration) error {
	return c.adjust(FOCUS_FAR, duration, "焦点后调")
}

// IrisOpen 光圈扩大（变亮）
// 参数：
//   - duration: 持续时间
func (c *CameraController) IrisOpen(duration time.Duration) error {
	return c.adjust(IRIS_OPEN, duration, "光圈扩大")
}

// IrisClose 光圈缩小（变暗）
// 参数：
//   - duration: 持续时间
func (c *CameraController) IrisClose(duration time.Duration) error {
	return c.adjust(IRIS_CLOSE, duration, "光圈缩小")
}

// adjust 内部调整函数，自动处理开始和停止
func (c *CameraController) adjust(cmd int, duration time.Duration, actionName string) error {
	// 开始调整
	if err := c.control(cmd, PTZ_START); err != nil {
		return fmt.Errorf("%s失败: %w", actionName, err)
	}

	// 等待指定时间
	time.Sleep(duration)

	// 停止调整
	if err := c.control(cmd, PTZ_STOP); err != nil {
		return fmt.Errorf("停止%s失败: %w", actionName, err)
	}

	log.Printf("✓ %s（通道%d，持续%v）", actionName, c.channel, duration)
	return nil
}

// ==================== AuxiliaryController 辅助设备控制器 ====================

// AuxiliaryController 辅助设备控制器
// 封装灯光、雨刷、风扇、加热器等辅助设备控制
type AuxiliaryController struct {
	userID  int // 登录句柄
	channel int // 通道号
}

// NewAuxiliaryController 创建辅助设备控制器
func NewAuxiliaryController(userID int, channel int) *AuxiliaryController {
	return &AuxiliaryController{
		userID:  userID,
		channel: channel,
	}
}

// LightOn 接通灯光电源
func (a *AuxiliaryController) LightOn() error {
	return a.switchDevice(LIGHT_PWRON, true, "灯光")
}

// LightOff 关闭灯光电源
func (a *AuxiliaryController) LightOff() error {
	return a.switchDevice(LIGHT_PWRON, false, "灯光")
}

// WiperOn 接通雨刷
func (a *AuxiliaryController) WiperOn() error {
	return a.switchDevice(WIPER_PWRON, true, "雨刷")
}

// WiperOff 关闭雨刷
func (a *AuxiliaryController) WiperOff() error {
	return a.switchDevice(WIPER_PWRON, false, "雨刷")
}

// FanOn 接通风扇
func (a *AuxiliaryController) FanOn() error {
	return a.switchDevice(FAN_PWRON, true, "风扇")
}

// FanOff 关闭风扇
func (a *AuxiliaryController) FanOff() error {
	return a.switchDevice(FAN_PWRON, false, "风扇")
}

// HeaterOn 接通加热器
func (a *AuxiliaryController) HeaterOn() error {
	return a.switchDevice(HEATER_PWRON, true, "加热器")
}

// HeaterOff 关闭加热器
func (a *AuxiliaryController) HeaterOff() error {
	return a.switchDevice(HEATER_PWRON, false, "加热器")
}

// AuxDevice1On 接通辅助设备1
func (a *AuxiliaryController) AuxDevice1On() error {
	return a.switchDevice(AUX_PWRON1, true, "辅助设备1")
}

// AuxDevice1Off 关闭辅助设备1
func (a *AuxiliaryController) AuxDevice1Off() error {
	return a.switchDevice(AUX_PWRON1, false, "辅助设备1")
}

// AuxDevice2On 接通辅助设备2
func (a *AuxiliaryController) AuxDevice2On() error {
	return a.switchDevice(AUX_PWRON2, true, "辅助设备2")
}

// AuxDevice2Off 关闭辅助设备2
func (a *AuxiliaryController) AuxDevice2Off() error {
	return a.switchDevice(AUX_PWRON2, false, "辅助设备2")
}

// switchDevice 内部设备开关函数
func (a *AuxiliaryController) switchDevice(cmd int, turnOn bool, deviceName string) error {
	action := PTZ_START // 0=开启
	actionName := "开启"
	if !turnOn {
		action = PTZ_STOP // 1=关闭
		actionName = "关闭"
	}

	if err := a.control(cmd, action); err != nil {
		return fmt.Errorf("%s%s失败: %w", actionName, deviceName, err)
	}

	log.Printf("✓ %s%s（通道%d）", actionName, deviceName, a.channel)
	return nil
}

// ==================== 通用底层控制函数 ====================

// controlWithSpeed 带速度的云台控制（MovementController 使用）
// 根据官方文档，速度范围是 1-7
func (m *MovementController) controlWithSpeed(cmd, stop, speed int) error {
	if m.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", m.userID)
	}

	ret := C.NET_DVR_PTZControlWithSpeed_Other(
		C.LONG(m.userID),
		C.LONG(m.channel),
		C.DWORD(cmd),
		C.DWORD(stop),
		C.DWORD(speed),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("云台控制失败 [通道:%d Cmd:%d 错误码:%d]", m.channel, cmd, errCode)
	}

	return nil
}

// validateSpeed 验证速度范围
func (m *MovementController) validateSpeed(speed int) error {
	if speed < MinSpeed || speed > MaxSpeed {
		return fmt.Errorf("速度超出范围：%d（有效范围：%d-%d）", speed, MinSpeed, MaxSpeed)
	}
	return nil
}

// control 无速度的云台控制（CameraController 和 AuxiliaryController 使用）
func (c *CameraController) control(cmd, stop int) error {
	if c.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", c.userID)
	}

	ret := C.NET_DVR_PTZControl_Other(
		C.LONG(c.userID),
		C.LONG(c.channel),
		C.DWORD(cmd),
		C.DWORD(stop),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("相机控制失败 [通道:%d Cmd:%d 错误码:%d]", c.channel, cmd, errCode)
	}

	return nil
}

// control 无速度的云台控制（AuxiliaryController 使用）
func (a *AuxiliaryController) control(cmd, stop int) error {
	if a.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", a.userID)
	}

	ret := C.NET_DVR_PTZControl_Other(
		C.LONG(a.userID),
		C.LONG(a.channel),
		C.DWORD(cmd),
		C.DWORD(stop),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("辅助设备控制失败 [通道:%d Cmd:%d 错误码:%d]", a.channel, cmd, errCode)
	}

	return nil
}
