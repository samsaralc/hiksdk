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

// 预置点命令常量（来自官方文档表 5.11）
const (
	// SET_PRESET 设置预置点
	SET_PRESET = 8
	// CLE_PRESET 清除预置点
	CLE_PRESET = 9
	// GOTO_PRESET 转到预置点
	GOTO_PRESET = 39
)

// 预置点参数限制（来自官方文档）
const (
	MinPresetID = 1 // 预置点最小编号
	// MaxPresetID 定义在 cruise.go 中（255），巡航和预置点共享此限制
)

// PresetManager 预置点控制器
// 封装了云台预置点的所有操作，提供简化的API
type PresetManager struct {
	userID  int // 登录句柄（NET_DVR_Login_V40 的返回值）
	channel int // 通道号
}

// NewPresetManager 创建预置点控制器
// 参数：
//   - userID: 设备登录ID (dev.GetLoginID())
//   - channel: 通道号
//
// 返回：
//   - *PresetManager: 预置点控制器实例
func NewPresetManager(userID int, channel int) *PresetManager {
	return &PresetManager{
		userID:  userID,
		channel: channel,
	}
}

// SetPreset 设置预置点
// 将云台当前位置保存为指定编号的预置点
// 对应官方命令：SET_PRESET
// 参数：
//   - presetID: 预置点编号（1-255）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (p *PresetManager) SetPreset(presetID int) error {
	// 参数验证
	if err := p.validatePresetID(presetID); err != nil {
		return err
	}

	if err := p.control(SET_PRESET, presetID); err != nil {
		return fmt.Errorf("设置预置点%d失败: %w", presetID, err)
	}

	log.Printf("✓ 设置预置点%d成功（通道%d）", presetID, p.channel)
	return nil
}

// GotoPreset 转到预置点
// 控制云台移动到指定预置点位置
// 对应官方命令：GOTO_PRESET
// 参数：
//   - presetID: 预置点编号（1-255）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (p *PresetManager) GotoPreset(presetID int) error {
	// 参数验证
	if err := p.validatePresetID(presetID); err != nil {
		return err
	}

	if err := p.control(GOTO_PRESET, presetID); err != nil {
		return fmt.Errorf("转到预置点%d失败: %w", presetID, err)
	}

	log.Printf("✓ 转到预置点%d（通道%d）", presetID, p.channel)
	return nil
}

// DeletePreset 删除预置点
// 清除指定编号的预置点
// 对应官方命令：CLE_PRESET
// 参数：
//   - presetID: 预置点编号（1-255）
//
// 返回：
//   - error: 错误信息，成功时为nil
func (p *PresetManager) DeletePreset(presetID int) error {
	// 参数验证
	if err := p.validatePresetID(presetID); err != nil {
		return err
	}

	if err := p.control(CLE_PRESET, presetID); err != nil {
		return fmt.Errorf("删除预置点%d失败: %w", presetID, err)
	}

	log.Printf("✓ 删除预置点%d成功（通道%d）", presetID, p.channel)
	return nil
}

// validatePresetID 验证预置点编号范围
func (p *PresetManager) validatePresetID(presetID int) error {
	if presetID < MinPresetID || presetID > MaxPresetID {
		return fmt.Errorf("预置点编号超出范围：%d（有效范围：%d-%d）", presetID, MinPresetID, MaxPresetID)
	}
	return nil
}

// control 内部通用控制函数
// 直接调用 NET_DVR_PTZPreset_Other（推荐，不需要预览）
func (p *PresetManager) control(cmd, presetID int) error {
	if p.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", p.userID)
	}

	// 调用 C 接口
	ret := C.NET_DVR_PTZPreset_Other(
		C.LONG(p.userID),
		C.LONG(p.channel),
		C.DWORD(cmd),
		C.DWORD(presetID),
	)

	if ret != C.TRUE {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("预置点操作失败 [通道:%d Cmd:%d PresetID:%d 错误码:%d]",
			p.channel, cmd, presetID, errCode)
	}

	return nil
}

// GetPresetCommandName 获取预置点命令的名称（用于调试）
func GetPresetCommandName(cmd int) string {
	names := map[int]string{
		SET_PRESET:  "设置预置点",
		CLE_PRESET:  "清除预置点",
		GOTO_PRESET: "转到预置点",
	}
	if name, ok := names[cmd]; ok {
		return name
	}
	return fmt.Sprintf("未知命令(%d)", cmd)
}
