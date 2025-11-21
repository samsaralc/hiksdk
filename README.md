# HikSDK - 海康威视 Go SDK

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

> 🎥 海康威视设备的 Go 语言封装，覆盖登录、设备信息、PTZ 云台控制、视频预览、报警监听等常用功能。

## 📖 简介

HikSDK 是海康威视官方 C SDK 的 Go 语言封装，通过 CGO 调用底层 SDK，提供简洁易用的 Go API。支持网络摄像机（IPC）、网络视频录像机（NVR）、数字视频录像机（DVR）等全系列海康设备。

## ✨ 功能特性

- ✅ **设备管理**：设备登录/登出，获取设备信息，通道名称读取
- ✅ **PTZ 控制**：云台旋转、变焦、焦点调节、预置点管理（配套 `consts` 常量包）
- ✅ **视频预览**：实时视频预览，输出 PS 流，配合 `Receiver` 管道进行消费
- ✅ **报警监听**：报警回调、布防/撤防、报警输出控制
- ✅ **错误处理**：统一的 `HKError` 错误类型，带错误码和说明
- ✅ **跨平台支持**：兼容 Windows / Linux amd64（内置海康官方 SDK 动态库）

## 🌍 跨平台兼容性

本项目已完美适配 **Windows** 和 **Linux** 平台，经过深度优化，确保最佳性能和稳定性：

### 平台支持
- ✅ **Windows x64**：使用 MinGW-w64 编译，支持 Windows 10/11
- ✅ **Linux x64**：使用 GCC 编译，支持主流发行版（Ubuntu、Debian、CentOS、Arch等）
- ✅ **自动平台检测**：通过 CGO 构建标签自动选择正确的库和类型定义
- ✅ **类型兼容性**：完整处理所有 Windows/Linux 类型差异（HWND、句柄类型、调用约定等）

### 技术特性
- 🔧 **智能链接顺序**：优化了库链接顺序，确保依赖关系正确
- 🧵 **线程安全**：使用互斥锁和生命周期标记管理 SDK 的初始化与清理，可在 `Cleanup()` 后安全重新初始化
- 📦 **资源管理**：使用 `cgo.Handle` 正确管理 Go 与 C 之间的资源传递
- 🎯 **零拷贝优化**：在回调函数中使用零拷贝技术提高性能
- 🛡️ **安全边界检查**：所有 C 字符串操作都有边界检查，防止缓冲区溢出

> 💡 **技术说明**：本项目使用 CGO 调用海康 C SDK，采用了多项跨平台最佳实践，确保代码在不同平台上的一致性和可靠性。

## 📥 安装配置

### 前置要求：CGO 环境准备

本项目使用 CGO 调用海康 C SDK，因此需要确保 CGO 已启用并安装 C 编译器。

#### 1. 启用 CGO

CGO 在 Go 中默认是启用的。如果之前手动禁用过，需要重新启用：

```bash
# Linux / macOS
export CGO_ENABLED=1

# Windows PowerShell
$env:CGO_ENABLED=1
```

**永久启用（推荐）：**

- **Linux / macOS**：添加到 `~/.bashrc` 或 `~/.zshrc`
  ```bash
  echo 'export CGO_ENABLED=1' >> ~/.bashrc
  source ~/.bashrc
  ```

- **Windows**：在系统环境变量中添加 `CGO_ENABLED=1`

#### 2. 安装 C 编译器

##### Windows

需要安装 **MinGW-w64**（GCC for Windows）：

**方式一：使用 MSYS2（推荐）**

1. 下载并安装 [MSYS2](https://www.msys2.org/)
2. 打开 MSYS2 终端，执行：
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```
3. 将 MinGW bin 目录添加到系统 PATH（通常是 `C:\msys64\mingw64\bin`）

**方式二：独立安装包**

1. 从 [MinGW-w64 下载页面](https://github.com/niXman/mingw-builds-binaries/releases) 下载最新版本
2. 解压到目录，如 `C:\mingw64`
3. 将 `C:\mingw64\bin` 添加到系统 PATH

**验证安装：**
```bash
gcc --version
```

##### Linux

安装 GCC 编译器：

**Ubuntu / Debian：**
```bash
sudo apt-get update
sudo apt-get install -y build-essential gcc
```

**CentOS / RHEL / Fedora：**
```bash
sudo yum groupinstall "Development Tools"
sudo yum install gcc
```

**Arch Linux：**
```bash
sudo pacman -S base-devel gcc
```

**验证安装：**
```bash
gcc --version
```

##### macOS

安装 Xcode Command Line Tools：

```bash
xcode-select --install
```

**验证安装：**
```bash
gcc --version
```

---

### 步骤 1：安装 SDK

```bash
go get github.com/samsaralc/hiksdk
```

### 步骤 2：配置动态库路径（⭐ 一次性配置，永久生效）

由于使用了 CGO，需要让系统能找到海康 SDK 的动态库。**配置一次后，之后就完全开箱即用！**

#### Windows

**复制粘贴以下命令到 PowerShell（管理员权限）：**

```powershell
# 自动查找 SDK 库路径并永久添加到 PATH
$hiksdkPath = Get-ChildItem -Path "$env:GOPATH\pkg\mod\github.com\samsaralc" -Filter "hiksdk@*" -Directory | Select-Object -First 1
if ($hiksdkPath) {
    $libPath = "$($hiksdkPath.FullName)\lib\Windows"
    [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$libPath", "User")
    Write-Host "✓ 库路径已永久添加: $libPath" -ForegroundColor Green
    Write-Host "✓ 重启终端后生效" -ForegroundColor Yellow
} else {
    Write-Host "✗ 未找到 hiksdk，请先运行 go get github.com/samsaralc/hiksdk" -ForegroundColor Red
}
```

**说明：**
- `hiksdk@*` 中的 `@*` 是通配符，匹配任何版本（如 `hiksdk@v1.0.0`）
- 配置后无论升级到哪个版本都有效

#### Linux / macOS

**复制粘贴以下命令到终端：**

```bash
# 自动查找 SDK 库路径并永久添加到 ~/.bashrc
HIKSDK_LIB=$(find $GOPATH/pkg/mod/github.com/samsaralc -name "hiksdk@*" -type d 2>/dev/null | head -1)/lib/Linux
if [ -n "$HIKSDK_LIB" ]; then
    echo "export LD_LIBRARY_PATH=\$LD_LIBRARY_PATH:$HIKSDK_LIB" >> ~/.bashrc
    source ~/.bashrc
    echo "✓ 库路径已永久添加: $HIKSDK_LIB"
    echo "✓ 已生效，可以直接使用"
else
    echo "✗ 未找到 hiksdk，请先运行 go get github.com/samsaralc/hiksdk"
fi
```

**说明：**
- `hiksdk@*` 中的 `@*` 是版本号通配符（如 `hiksdk@v1.0.0`）
- 如果使用 zsh，将 `~/.bashrc` 改为 `~/.zshrc`

### 步骤 3：验证配置 ✅

重启终端，然后运行：

```bash
# 克隆本项目测试（或直接在你的项目中使用）
git clone https://github.com/samsaralc/hiksdk.git
cd hiksdk

# 修改 examples/01_login_methods.go 中的 IP、用户名、密码
# 然后运行
go run examples/01_login_methods.go
```

**如果能看到输出，配置成功！** 🎉

之后在任何项目中使用 SDK 都无需再配置，**真正的开箱即用！**

## 🔧 系统要求

- **Go 版本**：1.25 或更高版本
- **操作系统**：Windows 或 Linux（支持 amd64 架构）
- **CGO 支持**：必须启用 CGO（`CGO_ENABLED=1`）
- **C 编译器**：
  - Windows：MinGW-w64 (gcc)
  - Linux：gcc（通过 build-essential 或 Development Tools 安装）
  - macOS：Xcode Command Line Tools
- **硬件**：海康威视网络设备（IPC、NVR、DVR 等）
- **网络**：设备需要在网络可达范围内

> 💡 **重要提示**：由于本项目使用 CGO 调用海康 C SDK，必须确保已安装 C 编译器并启用 CGO。详见 [安装配置](#-安装配置) 部分。

## 🚀 快速开始

> 💡 **提示**：完成上方"安装配置"两步后，就可以在任何地方直接使用了！

### 1. 基础使用

#### 最简示例（v1.4.0+ 推荐方式）

```go
package main

import (
	"fmt"
	"github.com/samsaralc/hiksdk/core"
)

func main() {
	// 直接创建设备实例（SDK会自动初始化，无需手动调用）
	deviceInfo := core.DeviceInfo{
		IP:       "192.168.1.64",
		Port:     8000,
		UserName: "admin",
		Password: "password",
	}
	dev := core.NewHKDevice(deviceInfo) // ✨ 自动初始化SDK

	// 登录设备（推荐使用V40版本）
	loginId, err := dev.LoginV40()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	defer dev.Logout()
	fmt.Printf("登录成功，ID: %d\n", loginId)

	// 获取设备信息
	info, _ := dev.GetDeviceInfo()
	fmt.Printf("设备名称: %s\n", info.DeviceName)
	fmt.Printf("通道数量: %d\n", info.ByChanNum)
	
	// 程序结束时清理SDK（可选）
	defer core.Cleanup()
}
```

> ✨ **设计说明**:
> - SDK 会在第一次调用 `NewHKDevice()` 时自动初始化，进程结束前只需调用一次 `Cleanup()` 即可清理资源
> - 内部使用互斥锁和状态标记保证初始化/清理的线程安全
> - 不再需要单独的 `InitHikSDK()`，示例中的用法已经反映了最新的生命周期设计

#### 完整示例（带错误处理）

```go
package main

import (
	"fmt"
	"github.com/samsaralc/hiksdk/core"
	"os"
)

func main() {
	// 1. 配置设备
	defer core.Cleanup()

	// 2. 创建设备
	dev := core.NewHKDevice(core.DeviceInfo{
		IP:       "192.168.1.64",
		Port:     8000,
		UserName: "admin",
		Password: "password",
	})

	// 3. 登录
	loginId, err := dev.LoginV30()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		os.Exit(1)
	}
	defer dev.Logout()
	fmt.Printf("✓ 登录成功 (ID: %d)\n", loginId)

	// 4. 获取信息
	info, err := dev.GetDeviceInfo()
	if err != nil {
		fmt.Printf("获取信息失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 设备名称: %s\n", info.DeviceName)
	fmt.Printf("✓ 序列号: %s\n", info.DeviceID)
	fmt.Printf("✓ 通道数: %d\n", info.ByChanNum)

	// 5. 获取通道列表
	channels, err := dev.GetChannelName()
	if err == nil {
		for id, name := range channels {
			fmt.Printf("  - 通道 %d: %s\n", id, name)
		}
	}
}
```

### 2. PTZ 云台控制

#### 使用预定义常量

```go
import (
	"time"
	"github.com/samsaralc/hiksdk/consts"
	"github.com/samsaralc/hiksdk/core"
)

// SDK 已经提供了所有 PTZ 命令常量，直接使用
// consts.PAN_RIGHT, consts.TILT_UP, consts.ZOOM_IN 等

// 云台右转（推荐使用 PTZControlWithSpeed_Other）
success, err := dev.PTZControlWithSpeed_Other(
	1,                 // 通道号
	consts.PAN_RIGHT,  // PTZ命令：右转
	consts.PTZ_START,  // 开始动作
	4,                 // 速度：0-7
)
if err == nil && success {
	time.Sleep(2 * time.Second)
	// 停止
	dev.PTZControlWithSpeed_Other(1, consts.PAN_RIGHT, consts.PTZ_STOP, 4)
}

// 变焦放大（不需要速度参数）
dev.PTZControl_Other(1, consts.ZOOM_IN, consts.PTZ_START)
time.Sleep(1 * time.Second)
dev.PTZControl_Other(1, consts.ZOOM_IN, consts.PTZ_STOP) // 停止
```

#### PTZ 控制完整示例

```go
import (
	"time"

	"github.com/samsaralc/hiksdk/consts"
	"github.com/samsaralc/hiksdk/core"
)

// 云台移动示例
func ptzMoveExample(dev *core.HKDevice) {
	channelId := 1
	
	// 右转 2 秒
	dev.PTZControlWithSpeed_Other(channelId, consts.PAN_RIGHT, consts.PTZ_START, 4)
	time.Sleep(2 * time.Second)
	dev.PTZControlWithSpeed_Other(channelId, consts.PAN_RIGHT, consts.PTZ_STOP, 4)
	
	// 上仰 2 秒
	dev.PTZControlWithSpeed_Other(channelId, consts.TILT_UP, consts.PTZ_START, 4)
	time.Sleep(2 * time.Second)
	dev.PTZControlWithSpeed_Other(channelId, consts.TILT_UP, consts.PTZ_STOP, 4)
	
	// 右上斜向移动
	dev.PTZControlWithSpeed_Other(channelId, consts.UP_RIGHT, consts.PTZ_START, 3)
	time.Sleep(2 * time.Second)
	dev.PTZControlWithSpeed_Other(channelId, consts.UP_RIGHT, consts.PTZ_STOP, 3)
}

// 变焦和焦点控制示例
func zoomFocusExample(dev *core.HKDevice) {
	channelId := 1
	
	// 焦距放大（拉近）
	dev.PTZControl_Other(channelId, consts.ZOOM_IN, consts.PTZ_START)
	time.Sleep(1 * time.Second)
	dev.PTZControl_Other(channelId, consts.ZOOM_IN, consts.PTZ_STOP)
	
	// 焦点前调（聚焦）
	dev.PTZControl_Other(channelId, consts.FOCUS_NEAR, consts.PTZ_START)
	time.Sleep(1 * time.Second)
	dev.PTZControl_Other(channelId, consts.FOCUS_NEAR, consts.PTZ_STOP)
}

// 预置点使用示例
func presetExample(dev *core.HKDevice) {
	channelId := 1
	presetId := 1
	
	// 设置预置点1
	dev.PTZControl_Other(channelId, consts.SET_PRESET, presetId)
	
	// 移动云台到其他位置
	dev.PTZControlWithSpeed_Other(channelId, consts.PAN_RIGHT, consts.PTZ_START, 4)
	time.Sleep(3 * time.Second)
	dev.PTZControlWithSpeed_Other(channelId, consts.PAN_RIGHT, consts.PTZ_STOP, 4)
	
	// 转到预置点1
	dev.PTZControl_Other(channelId, consts.GOTO_PRESET, presetId)
	time.Sleep(2 * time.Second) // 等待云台移动到位
}
```

### 3. 视频预览

```go
// 创建接收器
receiver := &core.Receiver{}
receiver.Start()

// 启动实时预览
realHandle, err := dev.RealPlay_V40(channelId, receiver)
if err != nil {
	fmt.Printf("预览失败: %v\n", err)
	return
}
defer dev.StopRealPlay()

// 接收视频数据
for data := range receiver.PSMouth {
	// 处理视频数据包
	fmt.Printf("收到数据包，大小: %d bytes\n", len(data))
}
```

### 4. 报警监听

```go
// 设置报警回调
err := dev.SetAlarmCallBack()
if err != nil {
	fmt.Printf("设置回调失败: %v\n", err)
	return
}

// 启动报警监听
err = dev.StartListenAlarmMsg()
if err != nil {
	fmt.Printf("启动监听失败: %v\n", err)
	return
}
defer dev.StopListenAlarmMsg()

// 报警事件将通过回调函数触发
// 等待报警事件...
```

## 📁 项目结构

```
hiksdk/
├── core/                   # SDK 核心包
│   ├── device.go          # 设备管理和初始化
│   ├── login.go           # 登录认证和动态IP解析
│   ├── config.go          # 设备配置管理
│   ├── video.go           # 视频预览功能
│   ├── ptz.go             # PTZ云台控制
│   ├── ptz_commands.go    # PTZ命令常量（63个命令）
│   ├── alarm.go           # 报警监听功能
│   ├── errors.go          # 错误处理
│   ├── helpers.go         # 工具函数
│   ├── transceiver.go     # PS流数据收发器
│   └── hiksdk_wrapper.h   # CGO头文件
├── examples/               # 可运行的示例代码
│   ├── 01_login_methods.go    # 登录方式示例
│   ├── 02_device_info.go      # 设备信息示例
│   ├── 03_ptz_control.go      # PTZ控制示例
│   ├── 04_video_preview.go    # 视频预览示例
│   ├── 05_alarm_listen.go     # 报警监听示例
│   └── README.md              # 示例说明
├── docs/                   # 文档目录
│   ├── 用户.md            # 官方接口文档
│   ├── 预置点.md          # 预置点说明
│   ├── 错误代码及说明.md  # 错误代码参考
│   └── LOGIN_MODES.md     # 登录方式说明
├── include/                # C 头文件
│   ├── HCNetSDK.h         # 海康SDK主头文件
│   ├── DataType.h         # 数据类型定义
│   ├── DecodeCardSdk.h    # 解码卡SDK
│   └── plaympeg4.h        # MPEG4播放
├── lib/                    # 动态链接库
│   ├── Windows/           # Windows平台DLL
│   └── Linux/             # Linux平台SO
├── go.mod                  # Go模块定义
├── LICENSE                 # MIT许可证
├── CONTRIBUTING.md         # 贡献指南
└── README.md              # 本文件
```

## 📚 API 文档

### 核心概念

#### 两种句柄（Handle）

| 句柄类型 | 中文名 | 获取方式 | 作用域 | 用途 |
|---------|--------|---------|--------|------|
| **loginId** | 登录句柄 | `LoginV30()`/`LoginV40()` | 设备级别 | 设备配置、PTZ控制（配合通道号） |
| **lRealHandle** | 预览句柄 | `RealPlay_V40()` | 视频流级别 | 视频流控制、PTZ控制（当前预览通道） |

**详细说明：** 查看 [docs/HANDLE_EXPLANATION.md](docs/HANDLE_EXPLANATION.md)

---

### 设备管理

#### 1. 初始化与清理

```go
// v1.4.0+ 新版本（推荐）：自动初始化
func main() {
	// 直接创建设备，SDK自动初始化
	dev := core.NewHKDevice(deviceInfo)
	
	// 可选：程序退出时清理资源
	defer core.Cleanup()
	// ... 你的代码
}

// 兼容旧版本：手动初始化（仍然支持）
func main() {
	core.InitHikSDK()  // 手动初始化
	defer core.Cleanup() // 清理资源
	// ... 你的代码
}
```

#### 2. 设备登录

```go
// 创建设备实例
deviceInfo := core.DeviceInfo{
	IP:       "192.168.1.64",  // 设备IP
	Port:     8000,             // 端口（默认8000）
	UserName: "admin",          // 用户名
	Password: "password",       // 密码
}
dev := core.NewHKDevice(deviceInfo)

// 登录方式1：使用 Login_V30（兼容旧设备）
loginId, err := dev.LoginV30()
if err != nil {
	// 处理登录失败
	fmt.Printf("登录失败: %v\n", err)
}

// 登录方式2：使用 Login_V40（推荐，性能更好）
loginId, err := dev.LoginV40()

// 登出（释放连接）
err := dev.Logout()
```

#### 3. 获取设备信息

```go
// 获取设备详细信息
info, err := dev.GetDeviceInfo()

// 设备信息字段
info.DeviceName  // 设备名称
info.DeviceID    // 设备序列号
info.ByChanNum   // 通道数量
info.IP          // IP地址
info.Port        // 端口
info.UserName    // 用户名

// 获取所有通道名称
channels, err := dev.GetChannelName()
// 返回: map[int]string - 通道ID → 通道名称
// 示例: {1: "前门", 2: "后门", 3: "车库"}
```

---

### PTZ 云台控制

#### PTZ 控制方法对比

| 方法 | 参数 | 需要预览 | 推荐度 | 说明 |
|------|------|---------|--------|------|
| `PTZControlWithSpeed_Other` | 通道号+命令+速度 | ❌ 不需要 | ⭐⭐⭐⭐⭐ | **推荐使用** |
| `PTZControl_Other` | 通道号+命令 | ❌ 不需要 | ⭐⭐⭐⭐ | 无速度参数版本 |
| `PTZControlWithSpeed` | 命令+速度 | ✅ 需要 | ⭐⭐ | 需先调用RealPlay_V40 |
| `PTZControl` | 命令 | ✅ 需要 | ⭐⭐ | 需先调用RealPlay_V40 |

#### 1. 云台移动控制（推荐方法）

```go
// PTZControlWithSpeed_Other - 最灵活，不需要预览（推荐✅）
success, err := dev.PTZControlWithSpeed_Other(
	channelId,        // 通道号（1开始）
	core.PAN_RIGHT,    // PTZ命令
	core.PTZ_START,    // 0=开始，1=停止
	4,                // 速度：0-7
)

// PTZControl_Other - 无速度参数
success, err := dev.PTZControl_Other(
	channelId,         // 通道号
	consts.ZOOM_IN,    // PTZ命令
	consts.PTZ_START,  // 0=开始，1=停止
)

// PTZControlWithSpeed - 需要先启动预览
receiver := &core.Receiver{}
receiver.Start()
lRealHandle, _ := dev.RealPlay_V40(1, receiver)
// 现在可以使用（控制当前预览的通道）
success, err := dev.PTZControlWithSpeed(consts.PAN_RIGHT, consts.PTZ_START, 4)
```

#### 2. PTZ 命令常量（已内置 63 个命令）

```go
// ========== 基本移动（需要速度） ==========
consts.TILT_UP    = 21  // 云台上仰
consts.TILT_DOWN  = 22  // 云台下俯
consts.PAN_LEFT   = 23  // 云台左转
consts.PAN_RIGHT  = 24  // 云台右转

// ========== 组合移动（需要速度） ==========
consts.UP_LEFT    = 25  // 上仰+左转
consts.UP_RIGHT   = 26  // 上仰+右转
consts.DOWN_LEFT  = 27  // 下俯+左转
consts.DOWN_RIGHT = 28  // 下俯+右转

// ========== 焦距控制 ==========
consts.ZOOM_IN    = 11  // 焦距变大（拉近）
consts.ZOOM_OUT   = 12  // 焦距变小（拉远）

// ========== 焦点控制 ==========
consts.FOCUS_NEAR = 13  // 焦点前调
consts.FOCUS_FAR  = 14  // 焦点后调

// ========== 光圈控制 ==========
consts.IRIS_OPEN  = 15  // 光圈扩大（变亮）
consts.IRIS_CLOSE = 16  // 光圈缩小（变暗）

// ========== 预置点操作 ==========
consts.SET_PRESET  = 8   // 设置预置点
consts.CLE_PRESET  = 9   // 清除预置点
consts.GOTO_PRESET = 39  // 转到预置点

// ========== 辅助设备 ==========
consts.LIGHT_PWRON  = 2  // 接通灯光
consts.WIPER_PWRON  = 3  // 接通雨刷

// ========== 自动扫描 ==========
consts.PAN_AUTO   = 29  // 左右自动扫描
consts.PAN_CIRCLE = 50  // 圆周扫描

// ========== 巡航和轨迹 ==========
consts.RUN_SEQ         = 37  // 开始巡航
consts.STOP_SEQ        = 38  // 停止巡航
consts.RUN_CRUISE      = 36  // 开始轨迹
consts.STOP_CRUISE     = 44  // 停止轨迹

// ========== 组合控制（移动+变焦） ==========
consts.TILT_DOWN_ZOOM_IN  = 58  // 下俯+放大
consts.PAN_LEFT_ZOOM_IN   = 60  // 左转+放大
consts.PAN_RIGHT_ZOOM_IN  = 62  // 右转+放大
// ... 更多组合命令，共63个

// 查看完整命令列表：consts/ptz.go
```

#### 3. 辅助常量

```go
// 动作控制
consts.PTZ_START = 0  // 开始动作
consts.PTZ_STOP  = 1  // 停止动作

// 速度控制
consts.PTZ_SPEED_MIN     = 0  // 最小速度
consts.PTZ_SPEED_MAX     = 7  // 最大速度
consts.PTZ_SPEED_DEFAULT = 4  // 默认速度

// 获取命令名称（调试用）
name := consts.GetPTZCommandName(consts.PAN_RIGHT)
// 返回: "云台右转"
```

#### 4. 获取 PTZ 位置

```go
// 获取指定通道的 PTZ 当前位置
dev.GetChannelPTZ(channelId)
// 会打印：水平角度、垂直角度、变焦倍数
```

---

### 视频预览

#### 基本使用

```go
// 1. 创建接收器
receiver := &core.Receiver{}
err := receiver.Start()
if err != nil {
	fmt.Printf("接收器启动失败: %v\n", err)
	return
}

// 2. 启动实时预览
lRealHandle, err := dev.RealPlay_V40(channelId, receiver)
if err != nil {
	fmt.Printf("预览失败: %v\n", err)
	return
}
fmt.Printf("预览句柄: %d\n", lRealHandle)

// 3. 停止预览（重要！）
defer dev.StopRealPlay()

// 4. 接收视频数据（PS流格式）
for data := range receiver.PSMouth {
	// data 是 PS 流数据包
	fmt.Printf("收到数据包: %d bytes\n", len(data))
	
	// 你可以：
	// - 保存为文件
	// - 解析 PS 流提取 H.264/H.265
	// - 推送到流媒体服务器
	// - 进行视频分析
}
```

#### 完整的视频预览示例

```go
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/samsaralc/hiksdk/core"
)

func videoPreviewExample() {
	core.InitHikSDK()
	defer core.Cleanup()

	dev := core.NewHKDevice(deviceInfo)
	dev.LoginV40()
	defer dev.Logout()

	// 启动接收器
	receiver := &core.Receiver{}
	receiver.Start()

	// 启动预览
	lRealHandle, err := dev.RealPlay_V40(1, receiver)
	if err != nil || lRealHandle < 0 {
		fmt.Printf("预览失败: %v\n", err)
		return
	}
	defer dev.StopRealPlay()

	fmt.Println("视频预览已启动，按 Ctrl+C 退出")

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 统计数据
	count := 0
	startTime := time.Now()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n退出预览...")
			return

		case data, ok := <-receiver.PSMouth:
			if !ok {
				fmt.Println("接收通道已关闭")
				return
			}
			count++
			
			if count%100 == 0 {
				elapsed := time.Since(startTime)
				rate := float64(count) / elapsed.Seconds()
				fmt.Printf("\r已接收 %d 包 (%.2f pkt/s)    ", count, rate)
			}
		}
	}
}
```

---

### 报警监听

#### 基本使用

```go
// 1. 设置报警回调函数
err := dev.SetAlarmCallBack()
if err != nil {
	fmt.Printf("设置回调失败: %v\n", err)
	return
}

// 2. 启动报警监听
err = dev.StartListenAlarmMsg()
if err != nil {
	fmt.Printf("启动监听失败: %v\n", err)
	return
}

// 3. 停止报警监听（重要！）
defer dev.StopListenAlarmMsg()

// 报警事件将通过回调函数触发
// 监听期间程序保持运行
```

#### 支持的报警类型

- 移动侦测报警
- 遮挡报警
- 音频异常报警
- 硬盘满报警
- 硬盘故障报警
- 视频信号丢失报警
- 输入/输出报警
- 智能事件报警

#### 完整示例

参考 `examples/05_alarm_listen.go`

---

## 🧪 开发者测试

> ⚠️ **注意**：以下内容仅适用于**本项目的开发者**，想要运行项目自带的测试用例时使用。
> 
> **如果你只是想使用这个 SDK**，不需要运行测试，直接跳到 [快速开始](#快速开始) 部分即可。

### 1. 设置测试环境变量（仅用于测试）

本项目的测试用例需要连接真实设备，运行测试前需设置环境变量：

```bash
# Windows PowerShell
$env:HIK_IP="192.168.1.64"
$env:HIK_USER="admin"
$env:HIK_PASSWORD="your_password"

# Linux / macOS
export HIK_IP="192.168.1.64"
export HIK_USER="admin"
export HIK_PASSWORD="your_password"
```

### 2. 运行测试

```bash
# 运行所有测试
go test -v ./core/...

# 运行设备管理测试
go test -v ./core/ -run TestDevice

# 运行 PTZ 控制测试
go test -v ./core/ -run TestPTZ

# 运行基准测试（性能测试）
go test -v ./core/ -bench=. -benchmem

# 查看测试覆盖率
go test -v ./core/ -cover
```

### 3. 测试说明

- **测试用例仅用于开发者验证 SDK 功能**
- **普通用户使用 SDK 不需要运行测试**
- 如果没有设置环境变量或没有设备，测试会自动跳过：
  ```
  === RUN   TestDeviceLogin
      device_test.go:16: 跳过测试: 请设置环境变量 HIK_IP, HIK_USER, HIK_PASSWORD
  --- SKIP: TestDeviceLogin (0.00s)
  ```

---

## 📖 运行示例

配置好动态库路径后，可以直接运行示例：

```bash
# 运行示例（修改代码中的 IP、用户名、密码）
go run examples/01_login_methods.go
go run examples/03_ptz_control.go
go run examples/04_video_preview.go
go run examples/05_alarm_listen.go
```

### 示例说明

| 示例 | 文件 | 功能演示 |
|------|------|---------|
| 登录方式 | `01_login_methods.go` | V30/V40 登录、动态IP解析 |
| 设备信息 | `02_device_info.go` | 获取设备信息、通道列表 |
| PTZ 控制 | `03_ptz_control.go` | 云台移动、变焦、预置点 |
| 视频预览 | `04_video_preview.go` | 启动预览、接收 PS 流、统计 |
| 报警监听 | `05_alarm_listen.go` | 设置回调、监听报警事件 |

## ❓ 常见问题

### 1. 编译时报错 "C compiler gcc not found" 或 "cannot use nil as _Ctype_HWND value"

**原因**：未安装 C 编译器或未启用 CGO。

**解决方案**：

1. **确保 CGO 已启用**：
   ```bash
   # 检查 CGO 状态
   go env CGO_ENABLED
   
   # 如果输出是 0，需要启用
   export CGO_ENABLED=1  # Linux/macOS
   $env:CGO_ENABLED=1    # Windows PowerShell
   ```

2. **安装 C 编译器**：
   - **Windows**：安装 MinGW-w64，详见 [前置要求](#前置要求cgo-环境准备)
   - **Linux**：运行 `sudo apt-get install build-essential gcc`
   - **macOS**：运行 `xcode-select --install`

3. **验证安装**：
   ```bash
   gcc --version  # 应该显示版本信息
   ```

**关于 HWND 错误**：这是 Windows/Linux 类型差异导致的，已在最新版本中修复。请确保使用最新版本的 SDK。

### 2. 使用 SDK 需要配置设备信息的环境变量吗？

**完全不需要！** 设备信息直接在代码中传参：

```go
dev := core.NewHKDevice(core.DeviceInfo{
    IP:       "192.168.1.64",  // 直接写在代码里
    Port:     8000,
    UserName: "admin",
    Password: "password",
})
```

**只需要一次性配置动态库路径**（见上方"安装"部分），之后就完全开箱即用了。

**HIK_IP 等环境变量仅用于：**
- 运行本项目自带的测试用例 `go test ./pkg/...`
- 普通用户完全不需要

### 3. PTZControlWithSpeed 和 PTZControlWithSpeed_Other 的区别？

| 方法 | 需要预览 | 推荐度 | 说明 |
|------|---------|--------|------|
| `PTZControlWithSpeed_Other` | ❌ | ⭐⭐⭐⭐⭐ | **推荐**，不需要预览视频 |
| `PTZControlWithSpeed` | ✅ | ⭐⭐ | 需要先调用 `RealPlay_V40` |

**推荐使用 `PTZControlWithSpeed_Other`！**

详细说明：[docs/HANDLE_EXPLANATION.md](docs/HANDLE_EXPLANATION.md)

### 4. 为什么 PTZ 控制失败？

可能的原因：
- 设备不支持 PTZ 功能（普通固定摄像头没有云台）
- 通道号不正确（尝试 1, 2, 3... 不同的通道）
- 用户权限不足（检查用户是否有 PTZ 控制权限）
- 云台未连接或未启用（检查设备配置）
- 使用了错误的方法（推荐使用 `PTZControlWithSpeed_Other`）

### 5. 支持哪些设备？

支持海康威视全系列网络设备：
- 网络摄像机（IPC）
- 网络视频录像机（NVR）
- 数字视频录像机（DVR）
- 混合型 NVR/DVR

### 6. lib 目录下的文件很大，可以删除吗？

**不可以！** 这些是海康 SDK 的核心库文件，SDK 运行时必需。

- Windows: 需要 `lib/Windows/` 下的所有 DLL 文件
- Linux: 需要 `lib/Linux/` 下的所有 SO 文件
- 所有文件已被 Git 托管，克隆时会自动下载

### 7. 如何处理视频数据？

SDK 提供 PS 流数据，你需要：
1. 解析 PS 流
2. 提取 H.264/H.265 编码数据
3. 使用解码库（如 FFmpeg）解码
4. 渲染显示

### 8. 为什么连接超时？

检查：
- 设备 IP 地址是否正确
- 网络连接是否正常
- 防火墙是否阻止
- 设备端口是否开放（默认 8000）

## 🏗️ 架构设计

### 核心特性
- **模块化设计**：各功能模块独立（登录、PTZ、视频、报警、配置）
- **优雅的错误处理**：统一的错误类型，包含错误码和详细描述
- **自动资源管理**：使用 defer 模式确保资源正确释放
- **可扩展接口**：清晰的接口设计，易于扩展新功能

### 最佳实践
```go
// 1. 使用 defer 确保资源释放
func main() {
    // SDK 会自动初始化
    dev := core.NewHKDevice(deviceInfo)
    
    // 程序退出时清理 SDK
    defer core.Cleanup()
    
    // 登录设备
    _, err := dev.LoginV40()
    if err != nil {
        log.Fatal(err)
    }
    defer dev.Logout() // 确保登出
    
    // 使用设备...
}

// 2. 错误处理
loginId, err := dev.LoginV40()
if err != nil {
    if hkErr, ok := err.(*core.HKError); ok {
        log.Printf("错误码: %d, 描述: %s", hkErr.Code, hkErr.Msg)
    }
    return
}

// 3. 资源清理
receiver := &core.Receiver{}
receiver.Start()
defer receiver.Stop() // 确保停止接收器

lRealHandle, _ := dev.RealPlay_V40(1, receiver)
defer dev.StopRealPlay() // 确保停止预览
```

## 注意事项

1. **线程安全**：
   - SDK 初始化使用 `sync.Once`，保证线程安全
   - 设备操作建议在单个 goroutine 中进行
   - 回调函数会在单独的线程中执行

2. **资源管理**：
   - ✅ 务必使用 `defer dev.Logout()` 释放连接
   - ✅ 务必使用 `defer core.Cleanup()` 清理 SDK
   - ✅ 停止预览时会自动清理 cgo.Handle

3. **错误处理**：
   - 所有 API 都返回详细的错误信息
   - 错误类型为 `*HKError`，包含错误码和描述
   - 建议使用类型断言获取详细错误信息

4. **设备限制**：
   - 某些功能受设备型号和固件版本限制
   - 建议在测试环境验证功能可用性
   - 使用 V40 版本的 API 获得更好的兼容性

5. **并发连接**：
   - 同一设备支持的并发连接数有限（通常 128-512）
   - 建议复用连接而非频繁创建/销毁
   - 可以使用连接池模式管理多个设备

## 性能优化

- 使用 `LoginV40` 替代 `LoginV30` 可获得更好的性能
- PTZ 控制建议使用 `PTZControlWithSpeed_Other` 而非 `PTZControlWithSpeed`
- 视频预览时注意处理数据流，避免缓冲区溢出
- 不使用时及时停止预览和监听，释放设备资源

## 依赖

本项目依赖海康威视官方 SDK 动态链接库：
- Windows: `HCNetSDK.dll` 及相关 DLL（位于 `lib/Windows/`）
- Linux: `libhcnetsdk.so` 及相关 SO（位于 `lib/Linux/`）

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题或建议，请创建 Issue。

## 致谢

- 感谢海康威视提供的 SDK 支持。
- 感谢 百川8488 支持。
---

**⭐ 如果这个项目对你有帮助，请给个 Star！**
