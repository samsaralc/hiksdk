package auth

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
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"unsafe"

	"github.com/samsaralc/hiksdk/core/utils"
)

var (
	// sdkMutex 保护SDK初始化和清理操作的互斥锁
	sdkMutex sync.Mutex
	// sdkInitialized 标记SDK是否已初始化
	sdkInitialized bool
)

// initSDK 初始化SDK（私有方法）
// 登录前会自动调用
func initSDK() error {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()

	// 如果已经初始化，直接返回
	if sdkInitialized {
		return nil
	}

	// 执行初始化
	result := C.NET_DVR_Init()
	if result != C.TRUE {
		return fmt.Errorf("SDK初始化失败")
	}

	// 设置连接超时参数
	C.NET_DVR_SetConnectTime(2000, 5) // 连接超时2秒，重试5次
	// 设置重连参数
	C.NET_DVR_SetReconnect(10000, 1) // 重连间隔10秒，启用重连

	sdkInitialized = true
	log.Println("✓ 海康SDK初始化成功")
	return nil
}

// Cleanup 清理SDK资源
// 通常在程序退出时调用
func Cleanup() error {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()

	if !sdkInitialized {
		return nil // 已经清理
	}

	// 执行清理
	result := C.NET_DVR_Cleanup()
	if result != C.TRUE {
		return fmt.Errorf("SDK清理失败")
	}

	sdkInitialized = false
	log.Println("✓ 海康SDK已清理")
	return nil
}

// SetLogConfig 配置SDK日志
// 参数：
//   - level: 日志级别（0-关闭, 1-错误, 2-警告, 3-信息, 4-调试）
//   - logDir: 日志存储目录（空字符串表示默认）
//   - autoDelete: 是否自动删除超量日志
//
// 返回值：
//   - error: 错误信息，成功时为nil
func SetLogConfig(level int, logDir string, autoDelete bool) error {
	// 确保SDK已初始化
	if err := initSDK(); err != nil {
		return err
	}

	cLogDir := C.CString(logDir)
	defer C.free(unsafe.Pointer(cLogDir))

	deleteFlag := C.BOOL(0)
	if autoDelete {
		deleteFlag = C.BOOL(1)
	}

	result := C.NET_DVR_SetLogToFile(
		C.DWORD(level),
		cLogDir,
		deleteFlag,
	)

	if result != C.TRUE {
		return fmt.Errorf("设置日志配置失败")
	}

	log.Printf("✓ SDK日志已配置 - 级别: %d, 目录: %s", level, logDir)
	return nil
}

// ==================== 认证功能 ====================

// Credentials 登录凭据
type Credentials struct {
	IP       string // 设备IP地址
	Port     int    // 设备端口
	Username string // 用户名
	Password string // 密码
}

// SessionInfo 会话信息
type SessionInfo struct {
	LoginID      int    // 登录ID
	SerialNumber string // 设备序列号
	ChannelNum   int    // 通道数量
}

// LoginV40 使用V40接口登录设备（推荐）
// 参数：
//   - cred: 登录凭据
//
// 返回值：
//   - *SessionInfo: 会话信息（包含loginID等）
//   - error: 错误信息，成功时为nil
func LoginV40(cred *Credentials) (*SessionInfo, error) {
	// 确保SDK已初始化（登录前必须调用）
	if err := initSDK(); err != nil {
		return nil, err
	}

	var deviceInfoV40 C.NET_DVR_DEVICEINFO_V40
	var userLoginInfo C.NET_DVR_USER_LOGIN_INFO

	// 设置登录参数
	utils.Strcpy(unsafe.Pointer(&userLoginInfo.sDeviceAddress[0]), cred.IP, len(userLoginInfo.sDeviceAddress))
	userLoginInfo.wPort = C.WORD(cred.Port)
	utils.Strcpy(unsafe.Pointer(&userLoginInfo.sUserName[0]), cred.Username, len(userLoginInfo.sUserName))
	utils.Strcpy(unsafe.Pointer(&userLoginInfo.sPassword[0]), cred.Password, len(userLoginInfo.sPassword))

	// 设置为同步登录模式（0=同步，1=异步）
	userLoginInfo.byUseAsynLogin = 0

	// 调用NET_DVR_Login_V40函数
	loginID := int(C.NET_DVR_Login_V40(&userLoginInfo, (*C.NET_DVR_DEVICEINFO_V40)(unsafe.Pointer(&deviceInfoV40))))

	// 提取设备序列号
	serialNumberBytes := make([]byte, len(deviceInfoV40.struDeviceV30.sSerialNumber))
	for i := range deviceInfoV40.struDeviceV30.sSerialNumber {
		serialNumberBytes[i] = byte(deviceInfoV40.struDeviceV30.sSerialNumber[i])
	}
	serialNumber := strings.Trim(string(serialNumberBytes), "\x00")

	if loginID < 0 {
		errCode := int(C.NET_DVR_GetLastError())
		return nil, fmt.Errorf("登录失败(V40)，错误码: %d", errCode)
	}

	session := &SessionInfo{
		LoginID:      loginID,
		SerialNumber: serialNumber,
		ChannelNum:   int(deviceInfoV40.struDeviceV30.byChanNum),
	}

	log.Printf("✓ 登录成功(V40) - 用户ID: %d, 设备序列号: %s, 通道数: %d",
		loginID, serialNumber, session.ChannelNum)
	return session, nil
}

// LoginV30 使用V30接口登录设备（兼容旧设备）
// 参数：
//   - cred: 登录凭据
//
// 返回值：
//   - *SessionInfo: 会话信息
//   - error: 错误信息，成功时为nil
func LoginV30(cred *Credentials) (*SessionInfo, error) {
	// 确保SDK已初始化（登录前必须调用）
	if err := initSDK(); err != nil {
		return nil, err
	}

	var deviceInfoV30 C.NET_DVR_DEVICEINFO_V30

	// 转换参数
	ip := C.CString(cred.IP)
	usr := C.CString(cred.Username)
	passwd := C.CString(cred.Password)
	defer func() {
		C.free(unsafe.Pointer(ip))
		C.free(unsafe.Pointer(usr))
		C.free(unsafe.Pointer(passwd))
	}()

	// 调用NET_DVR_Login_V30函数
	loginID := int(C.NET_DVR_Login_V30(
		ip,
		C.WORD(cred.Port),
		usr,
		passwd,
		(*C.NET_DVR_DEVICEINFO_V30)(unsafe.Pointer(&deviceInfoV30)),
	))

	// 提取设备序列号
	serialNumberBytes := make([]byte, len(deviceInfoV30.sSerialNumber))
	for i := range deviceInfoV30.sSerialNumber {
		serialNumberBytes[i] = byte(deviceInfoV30.sSerialNumber[i])
	}
	serialNumber := strings.Trim(string(serialNumberBytes), "\x00")

	if loginID < 0 {
		errCode := int(C.NET_DVR_GetLastError())
		return nil, fmt.Errorf("登录失败(V30)，错误码: %d", errCode)
	}

	session := &SessionInfo{
		LoginID:      loginID,
		SerialNumber: serialNumber,
		ChannelNum:   int(deviceInfoV30.byChanNum),
	}

	log.Printf("✓ 登录成功(V30) - 用户ID: %d, 设备序列号: %s, 通道数: %d",
		loginID, serialNumber, session.ChannelNum)
	return session, nil
}

// Logout 登出设备
// 参数：
//   - loginID: 登录ID
//
// 返回值：
//   - error: 错误信息，成功时为nil
func Logout(loginID int) error {
	if loginID < 0 {
		return nil // 未登录，不是错误
	}

	result := C.NET_DVR_Logout(C.LONG(loginID))
	if result == 0 {
		errCode := int(C.NET_DVR_GetLastError())
		return fmt.Errorf("登出失败，错误码: %d", errCode)
	}

	log.Printf("✓ 登出成功（LoginID: %d）", loginID)
	return nil
}

// ResolveDynamicIP 通过解析服务器获取设备的动态IP地址和端口
// 支持通过设备名称或序列号从IPServer/hiDDNS服务器解析设备当前IP
// 参数：
//   - serverIP: 解析服务器的IP地址或域名
//   - serverPort: 解析服务器端口（IPServer: 7071, hiDDNS: 80）
//   - dvrName: 设备名称或域名（可选，与序列号至少提供一个）
//   - serialNumber: 设备序列号（可选，与设备名称至少提供一个）
//
// 返回值：
//   - ip: 解析得到的设备IP地址
//   - port: 解析得到的设备端口号
//   - error: 错误信息
func ResolveDynamicIP(serverIP string, serverPort uint16, dvrName string, serialNumber string) (string, uint32, error) {
	// 验证输入参数
	if serverIP == "" {
		return "", 0, errors.New("解析服务器IP不能为空")
	}
	if dvrName == "" && serialNumber == "" {
		return "", 0, errors.New("设备名称和序列号不能同时为空")
	}

	// 转换参数为C类型
	cServerIP := C.CString(serverIP)
	defer C.free(unsafe.Pointer(cServerIP))

	// 准备设备名称参数
	var cDVRName *C.BYTE
	var dvrNameLen C.WORD
	if dvrName != "" {
		dvrNameBytes := []byte(dvrName)
		cDVRName = (*C.BYTE)(unsafe.Pointer(&dvrNameBytes[0]))
		dvrNameLen = C.WORD(len(dvrNameBytes))
	}

	// 准备序列号参数
	var cSerialNumber *C.BYTE
	var serialLen C.WORD
	if serialNumber != "" {
		serialBytes := []byte(serialNumber)
		cSerialNumber = (*C.BYTE)(unsafe.Pointer(&serialBytes[0]))
		serialLen = C.WORD(len(serialBytes))
	}

	// 准备输出缓冲区
	ipBuffer := make([]byte, 128)
	cGetIP := (*C.char)(unsafe.Pointer(&ipBuffer[0]))
	var dwPort C.DWORD

	// 调用SDK函数
	result := C.NET_DVR_GetDVRIPByResolveSvr_EX(
		cServerIP,
		C.WORD(serverPort),
		cDVRName,
		dvrNameLen,
		cSerialNumber,
		serialLen,
		cGetIP,
		&dwPort,
	)

	if result == 0 {
		errorCode := int(C.NET_DVR_GetLastError())
		return "", 0, fmt.Errorf("解析设备IP失败，错误码: %d", errorCode)
	}

	// 转换结果
	resolvedIP := C.GoString(cGetIP)
	resolvedPort := uint32(dwPort)

	log.Printf("✓ 动态IP解析成功 - IP: %s, 端口: %d", resolvedIP, resolvedPort)
	return resolvedIP, resolvedPort, nil
}

// LoginWithDynamicIP 使用动态IP解析后登录设备（V40版本）
// 先通过解析服务器获取设备当前IP，然后进行登录
// 参数：
//   - serverIP: 解析服务器地址
//   - serverPort: 解析服务器端口
//   - dvrName: 设备名称
//   - serialNumber: 设备序列号
//   - username: 登录用户名
//   - password: 登录密码
//
// 返回值：
//   - *SessionInfo: 会话信息
//   - error: 错误信息
func LoginWithDynamicIP(serverIP string, serverPort uint16, dvrName string, serialNumber string, username string, password string) (*SessionInfo, error) {
	resolvedIP, resolvedPort, err := ResolveDynamicIP(serverIP, serverPort, dvrName, serialNumber)
	if err != nil {
		return nil, err
	}

	cred := &Credentials{
		IP:       resolvedIP,
		Port:     int(resolvedPort),
		Username: username,
		Password: password,
	}

	return LoginV40(cred)
}
