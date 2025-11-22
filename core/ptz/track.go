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
)

// 轨迹命令常量（来自官方文档表 5.13）
const (
	// STA_MEM_CRUISE 开始记录花样扫描路径（轨迹）
	STA_MEM_CRUISE = 34
	// STO_MEM_CRUISE 停止记录花样扫描路径（轨迹）
	STO_MEM_CRUISE = 35
	// RUN_CRUISE 开始执行花样扫描路径（轨迹）
	RUN_CRUISE = 36
)

// TrackManager 轨迹控制器
// 封装了云台轨迹（花样扫描路径）的所有操作
type TrackManager struct {
	userID  int // 登录句柄（NET_DVR_Login_V40 的返回值）
	channel int // 通道号
}

// NewTrackManager 创建轨迹控制器
// 参数：
//   - userID: 设备登录ID (dev.GetLoginID())
//   - channel: 通道号
//
// 返回：
//   - *TrackManager: 轨迹控制器实例
func NewTrackManager(userID int, channel int) *TrackManager {
	return &TrackManager{
		userID:  userID,
		channel: channel,
	}
}

// StartRecordTrack 开始记录轨迹
// 对应官方命令：STA_MEM_CRUISE
// 调用后，云台的所有移动操作将被记录，直到调用 StopRecordTrack
//
// 返回：
//   - error: 错误信息，成功时为nil
func (t *TrackManager) StartRecordTrack() error {
	if err := t.control(STA_MEM_CRUISE); err != nil {
		return fmt.Errorf("开始记录轨迹失败: %w", err)
	}

	log.Printf("✓ 开始记录轨迹（通道%d）", t.channel)
	return nil
}

// StopRecordTrack 停止记录轨迹
// 对应官方命令：STO_MEM_CRUISE
// 停止记录，保存当前记录的轨迹
//
// 返回：
//   - error: 错误信息，成功时为nil
func (t *TrackManager) StopRecordTrack() error {
	if err := t.control(STO_MEM_CRUISE); err != nil {
		return fmt.Errorf("停止记录轨迹失败: %w", err)
	}

	log.Printf("✓ 停止记录轨迹（通道%d）", t.channel)
	return nil
}

// RunTrack 开始执行轨迹
// 对应官方命令：RUN_CRUISE
// 执行之前记录的轨迹，云台将按记录的路径运动
//
// 返回：
//   - error: 错误信息，成功时为nil
func (t *TrackManager) RunTrack() error {
	if err := t.control(RUN_CRUISE); err != nil {
		return fmt.Errorf("执行轨迹失败: %w", err)
	}

	log.Printf("✓ 开始执行轨迹（通道%d）", t.channel)
	return nil
}

// control 内部通用控制函数
// 直接调用 NET_DVR_PTZTrack_Other（推荐，不需要预览）
func (t *TrackManager) control(cmd int) error {
	if t.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", t.userID)
	}

	// 调用 C 接口
	ret := C.NET_DVR_PTZTrack_Other(
		C.LONG(t.userID),
		C.LONG(t.channel),
		C.DWORD(cmd),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("轨迹操作失败 [通道:%d Cmd:%d 错误码:%d]", t.channel, cmd, errCode)
	}

	return nil
}

// GetTrackCommandName 获取轨迹命令的名称（用于调试）
func GetTrackCommandName(cmd int) string {
	names := map[int]string{
		STA_MEM_CRUISE: "开始记录花样扫描路径",
		STO_MEM_CRUISE: "停止记录花样扫描路径",
		RUN_CRUISE:     "开始执行花样扫描路径",
	}
	if name, ok := names[cmd]; ok {
		return name
	}
	return fmt.Sprintf("未知命令(%d)", cmd)
}
