package alarm

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

// 声明Go回调函数
extern void AlarmCallBack(LONG command, NET_DVR_ALARMER *alarm, char *info, DWORD len, void *user);
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

// 报警类型常量（来自官方SDK）
const (
	// COMM_ALARM_RULE 行为分析报警
	COMM_ALARM_RULE = 0x1101
	// COMM_ALARM_V30 V30报警信息（移动侦测、视频丢失、遮挡、IO等）
	COMM_ALARM_V30 = 0x4000
	// COMM_ALARM_V40 V40报警信息（扩展报警类型）
	COMM_ALARM_V40 = 0x4007
)

// MAX_SERIALNO_LEN 序列号最大长度
const MAX_SERIALNO_LEN = 48

// AlarmListener 报警监听器
// 封装设备报警监听的所有操作
type AlarmListener struct {
	loginID     int // 登录句柄
	alarmHandle int // 报警句柄
}

// NewAlarmListener 创建报警监听器
// 参数：
//   - loginID: 设备登录ID
//
// 返回：
//   - *AlarmListener: 报警监听器实例
func NewAlarmListener(loginID int) *AlarmListener {
	return &AlarmListener{
		loginID:     loginID,
		alarmHandle: -1,
	}
}

// AlarmCallBack 报警回调函数
// 由C代码调用，接收设备的报警信息
// 这是一个全局回调函数，会被所有设备共享
//
//export AlarmCallBack
func AlarmCallBack(command C.LONG, alarm *C.NET_DVR_ALARMER, info *C.char, len C.DWORD, user unsafe.Pointer) {
	// 安全检查
	if alarm == nil {
		log.Println("警告: 收到空报警信息")
		return
	}

	// 获取报警设备信息
	var deviceIP string
	if alarm.byDeviceIPValid == 1 {
		// 提取设备IP（需要从DWORD转换为IP字符串）
		ip := alarm.dwDeviceIP
		deviceIP = fmt.Sprintf("%d.%d.%d.%d",
			byte(ip&0xFF),
			byte((ip>>8)&0xFF),
			byte((ip>>16)&0xFF),
			byte((ip>>24)&0xFF))
	}

	// 获取设备序列号
	var serialNo string
	if alarm.bySerialValid == 1 {
		// 手动复制序列号字节
		serialBytes := make([]byte, MAX_SERIALNO_LEN)
		for i := 0; i < MAX_SERIALNO_LEN; i++ {
			serialBytes[i] = byte(alarm.sSerialNumber[i])
		}
		// 去除空字符
		for i, b := range serialBytes {
			if b == 0 {
				serialNo = string(serialBytes[:i])
				break
			}
		}
		if serialNo == "" {
			serialNo = string(serialBytes)
		}
	}

	// 记录报警信息
	log.Printf("收到报警 - 类型: 0x%X, 设备IP: %s, 序列号: %s", command, deviceIP, serialNo)

	// 根据命令类型处理不同的报警
	switch int(command) {
	case COMM_ALARM_RULE:
		log.Println("  → 行为分析报警（规则检测）")
	case COMM_ALARM_V30:
		log.Println("  → 标准报警（移动侦测/视频丢失/遮挡/IO信号等）")
	case COMM_ALARM_V40:
		log.Println("  → 扩展报警V40（增强型报警信息）")
	default:
		log.Printf("  → 其他报警类型: 0x%X\n", command)
	}
}

// Start 启动报警监听
// 建立报警上传通道，开始接收设备的报警信息
// 返回值：
//   - error: 错误信息，成功时为nil
func (a *AlarmListener) Start() error {
	if a.loginID < 0 {
		return fmt.Errorf("无效的登录ID")
	}

	// 设置报警回调函数
	C.NET_DVR_SetDVRMessageCallBack_V30(
		(*[0]byte)(C.AlarmCallBack),
		nil,
	)

	// 设置报警布防参数
	var setupParam C.NET_DVR_SETUPALARM_PARAM
	setupParam.dwSize = C.DWORD(unsafe.Sizeof(setupParam))
	setupParam.byLevel = 1         // 布防等级
	setupParam.byAlarmInfoType = 1 // 上传报警信息类型：0-老报警信息，1-新报警信息

	// 建立报警上传通道
	a.alarmHandle = int(C.NET_DVR_SetupAlarmChan_V41(
		C.LONG(a.loginID),
		&setupParam,
	))

	if a.alarmHandle < 0 {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("建立报警通道失败，错误码: %d", errCode)
	}

	log.Printf("✓ 报警监听启动成功（句柄: %d）", a.alarmHandle)
	return nil
}

// Stop 停止报警监听
// 撤销报警上传通道，停止接收设备的报警信息
// 返回值：
//   - error: 错误信息，成功时为nil
func (a *AlarmListener) Stop() error {
	if a.alarmHandle >= 0 {
		if C.NET_DVR_CloseAlarmChan_V30(C.LONG(a.alarmHandle)) != C.TRUE {
			errCode := int(C.NET_DVR_GetLastError())
			return fmt.Errorf("关闭报警通道失败，错误码: %d", errCode)
		}

		a.alarmHandle = -1
		log.Println("✓ 报警监听已停止")
	}
	return nil
}
