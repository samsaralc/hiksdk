/**
 * @file hiksdk_wrapper.h
 * @brief 海康威视 SDK CGO 包装头文件
 * @author HikSDK Go Wrapper
 * @date 2024
 * 
 * 本文件为 CGO 调用海康 C SDK 提供必要的类型定义和函数声明
 * 支持 Windows 和 Linux 平台
 * 
 * 特性：
 * - 跨平台类型定义（Windows/Linux）
 * - 完整的SDK函数声明
 * - CGO友好的接口封装
 */

#ifndef HIKSDK_WRAPPER_H
#define HIKSDK_WRAPPER_H

// ========================================================================
// 平台检测与基础类型定义
// ========================================================================

// 平台相关头文件和类型定义
#if defined(_WIN32) || defined(_WIN64) || defined(__WIN32__) || defined(__WINDOWS__)
    // ============ Windows 平台 ============
    #ifndef WIN32_LEAN_AND_MEAN
        #define WIN32_LEAN_AND_MEAN
    #endif
    #include <windows.h>
    
    // 导出/导入声明
    #define HIKSDK_API __declspec(dllimport)
    // 调用约定
    #define CALLBACK __stdcall
    #define HIKSDK_CALL __stdcall
    
    // Windows 已经定义了所有需要的类型
    
#elif defined(__linux__) || defined(__unix__) || defined(__APPLE__)
    // ============ Linux/Unix/macOS 平台 ============
    #include <stdint.h>
    #include <stdbool.h>
    #include <string.h>
    
    // 导出/导入声明（Linux上无需）
    #define HIKSDK_API
    // 调用约定（Linux上无需）
    #define CALLBACK
    #define HIKSDK_CALL
    
    // Windows 类型在 Linux 上的等价定义
    typedef void* HWND;              // 窗口句柄（Linux上不使用，设为void*）
    typedef uint8_t BYTE;            // 8位无符号整数
    typedef uint16_t WORD;           // 16位无符号整数
    typedef uint32_t DWORD;          // 32位无符号整数
    typedef int32_t LONG;            // 32位有符号整数
    typedef int BOOL;                // 布尔类型（0=false, 非0=true）
    typedef void* LPVOID;            // 通用指针
    typedef const void* LPCVOID;     // 常量通用指针
    typedef char* LPSTR;             // 字符串指针
    typedef const char* LPCSTR;      // 常量字符串指针
    
#else
    #error "不支持的操作系统平台。仅支持 Windows 和 Linux/Unix。"
#endif

#ifdef __cplusplus
extern "C" {
#endif

/* ========================================================================
 * 基础常量定义
 * ======================================================================== */

// 布尔值定义（兼容C和C++）
#ifndef TRUE
    #define TRUE  1
#endif
#ifndef FALSE
    #define FALSE 0
#endif

// 字符串和数组最大长度定义
#define MAX_SERIALNO_LEN        48   // 序列号最大长度
#define MAX_USERNAME_LEN        64   // 用户名最大长度
#define MAX_PASSWORD_LEN        64   // 密码最大长度
#define MAX_DEVICE_ADDR_LEN     129  // 设备地址最大长度（IP或域名）
#define MAX_CHANNEL_NAME_LEN    32   // 通道名称最大长度
#define MAX_DEVICE_NAME_LEN     32   // 设备名称最大长度

// 网络参数默认值
#define DEFAULT_DEVICE_PORT     8000 // 设备默认端口
#define DEFAULT_CONNECT_TIMEOUT 5000 // 默认连接超时（毫秒）
#define DEFAULT_RECONNECT_INTERVAL 10000 // 默认重连间隔（毫秒）

// 常用错误码
#define NET_DVR_NOERROR             0   // 没有错误
#define NET_DVR_PASSWORD_ERROR      1   // 用户名密码错误
#define NET_DVR_NOENOUGHPRI         2   // 权限不足
#define NET_DVR_NOINIT              3   // 没有初始化
#define NET_DVR_CHANNEL_ERROR       4   // 通道号错误
#define NET_DVR_OVER_MAXLINK        5   // 连接到DVR的客户端个数超过最大
#define NET_DVR_VERSIONNOMATCH      6   // 版本不匹配
#define NET_DVR_NETWORK_FAIL_CONNECT 7  // 连接服务器失败
#define NET_DVR_NETWORK_SEND_ERROR  8   // 向服务器发送失败
#define NET_DVR_NETWORK_RECV_ERROR  9   // 从服务器接收数据失败
#define NET_DVR_NETWORK_RECV_TIMEOUT 10 // 从服务器接收数据超时

// 预览流类型（云台控制需要指定）
#define STREAM_TYPE_MAIN            0    // 主码流
#define STREAM_TYPE_SUB             1    // 子码流

/* ========================================================================
 * 数据结构定义 - 登录相关
 * ======================================================================== */

// 设备信息 V30 版本
typedef struct tagNET_DVR_DEVICEINFO_V30 {
    BYTE  sSerialNumber[MAX_SERIALNO_LEN];  // 序列号
    BYTE  byAlarmInPortNum;                 // 报警输入个数
    BYTE  byAlarmOutPortNum;                // 报警输出个数
    BYTE  byDiskNum;                        // 硬盘个数
    BYTE  byDVRType;                        // 设备类型
    BYTE  byChanNum;                        // 模拟通道个数
    BYTE  byStartChan;                      // 起始通道号，从1开始
    BYTE  byAudioChanNum;                   // 语音通道数
    BYTE  byIPChanNum;                      // 最大数字通道个数，低8位
    BYTE  byZeroChanNum;                    // 零通道编码个数
    BYTE  byMainProto;                      // 主码流传输协议类型
    BYTE  bySubProto;                       // 子码流传输协议类型
    BYTE  bySupport;                        // 能力集扩展，位与结果为0表示不支持，1表示支持
    BYTE  bySupport1;                       // 能力集扩展
    BYTE  bySupport2;                       // 能力集扩展
    WORD  wDevType;                         // 设备型号
    BYTE  bySupport3;                       // 能力集扩展
    BYTE  byMultiStreamProto;               // 是否支持多码流
    BYTE  byStartDChan;                     // 起始数字通道号
    BYTE  byStartDTalkChan;                 // 起始数字对讲通道号
    BYTE  byHighDChanNum;                   // 数字通道个数，高8位
    BYTE  bySupport4;                       // 能力集扩展
    BYTE  byLanguageType;                   // 语言类型
    BYTE  byVoiceInChanNum;                 // 音频输入通道数
    BYTE  byStartVoiceInChanNo;             // 音频输入起始通道号
    BYTE  byRes3[2];                        // 保留
    BYTE  byMirrorChanNum;                  // 镜像通道个数
    WORD  wStartMirrorChanNo;               // 起始镜像通道号
    BYTE  byRes2[2];                        // 保留
} NET_DVR_DEVICEINFO_V30, *LPNET_DVR_DEVICEINFO_V30;

// 设备信息 V40 版本（扩展V30）
typedef struct tagNET_DVR_DEVICEINFO_V40 {
    NET_DVR_DEVICEINFO_V30 struDeviceV30;   // V30设备信息
    BYTE  byRes1[4];                        // 保留
    BYTE  byRes2[12];                       // 保留
    BYTE  bySupport5;                       // 能力集扩展
    BYTE  byLanguageTypeEx;                 // 语言类型扩展
    BYTE  byRes3[54];                       // 保留
} NET_DVR_DEVICEINFO_V40, *LPNET_DVR_DEVICEINFO_V40;

// 用户登录信息
typedef struct tagNET_DVR_USER_LOGIN_INFO {
    BYTE  sDeviceAddress[MAX_DEVICE_ADDR_LEN]; // 设备地址（IP或域名）
    BYTE  byUseAsynLogin;                      // 是否异步登录：0-否，1-是
    BYTE  byProxyType;                         // 代理类型：0-不使用代理，1-使用标准代理，2-使用EHome代理
    BYTE  byUseUTCTime;                        // 是否使用UTC时间：0-不使用，1-使用
    BYTE  byLoginMode;                         // 登录模式：0-Private，1-ISAPI，2-自适应
    BYTE  byHttps;                             // 是否使用HTTPS：0-不使用，1-使用
    LONG  iProxyID;                           // 代理服务器序号
    BYTE  byVerifyMode;                       // 认证方式：0-不认证，1-双向认证，2-单向认证
    BYTE  byRes3[119];                        // 保留
    WORD  wPort;                              // 设备端口号
    BYTE  sUserName[MAX_USERNAME_LEN];        // 登录用户名
    BYTE  sPassword[MAX_PASSWORD_LEN];        // 登录密码
    void  (CALLBACK *fLoginResultCallBack)(LONG lUserID, DWORD dwResult, LPNET_DVR_DEVICEINFO_V30 lpDeviceInfo, void *pUser);
    void  *pUser;                             // 用户数据
} NET_DVR_USER_LOGIN_INFO, *LPNET_DVR_USER_LOGIN_INFO;

/* ========================================================================
 * 数据结构定义 - 报警相关
 * ======================================================================== */

// 报警设备信息
typedef struct tagNET_DVR_ALARMER {
    BYTE  sSerialNumber[MAX_SERIALNO_LEN];    // 序列号
    BYTE  byUserIDValid;                      // UserID是否有效
    BYTE  bySerialValid;                      // 序列号是否有效
    BYTE  byVersionValid;                     // 版本号是否有效
    BYTE  byDeviceNameValid;                  // 设备名字是否有效
    BYTE  byMacAddrValid;                     // MAC地址是否有效
    BYTE  byLinkPortValid;                    // 连接端口是否有效
    BYTE  byDeviceIPValid;                    // 设备IP是否有效
    BYTE  bySocketIPValid;                    // Socket IP是否有效
    LONG  lUserID;                           // 用户ID
    char  sDeviceName[MAX_DEVICE_NAME_LEN];   // 设备名字
    BYTE  byMacAddr[6];                       // MAC地址
    BYTE  byRes[2];                           // 保留
    DWORD dwDeviceIP;                         // 设备IP
    DWORD dwSocketIP;                         // Socket IP
    WORD  wLinkPort;                          // 连接端口
    BYTE  sSerialDevice[64];                  // 序列设备
    BYTE  byIpProtocol;                       // IP协议
    BYTE  byRes2[3];                          // 保留
} NET_DVR_ALARMER, *LPNET_DVR_ALARMER;

// 报警设置参数
typedef struct tagNET_DVR_SETUPALARM_PARAM {
    DWORD dwSize;                             // 结构体大小
    BYTE  byLevel;                            // 布防等级：0-一级，1-二级
    BYTE  byAlarmInfoType;                    // 报警信息类型：0-老报警信息，1-新报警信息
    BYTE  byRes1[2];                          // 保留
    DWORD dwSubScribeType;                    // 订阅类型
    BYTE  byRes[60];                          // 保留
} NET_DVR_SETUPALARM_PARAM, *LPNET_DVR_SETUPALARM_PARAM;

/* ========================================================================
 * 数据结构定义 - 预览相关
 * ======================================================================== */

// 预览信息（云台控制需要通过预览句柄操作）
typedef struct tagNET_DVR_PREVIEWINFO {
    LONG  lChannel;                           // 通道号
    DWORD dwStreamType;                       // 码流类型：0-主码流，1-子码流，2-三码流
    DWORD dwLinkMode;                         // 连接方式：0-TCP，1-UDP，2-多播，3-RTP，4-RTP/RTSP，5-RSTP/HTTP
    HWND  hPlayWnd;                           // 播放窗口句柄
    BOOL  bBlocked;                           // 是否阻塞取流：0-否，1-是
    BOOL  bPassbackRecord;                    // 是否启用录像回传：0-不启用，1-启用
    BYTE  byPreviewMode;                      // 预览模式：0-正常，1-延迟
    BYTE  byStreamID[32];                     // 流ID
    BYTE  byProtoType;                        // 应用层协议：0-私有协议，1-RTSP协议
    BYTE  byRes1;                             // 保留
    BYTE  byVideoCodingType;                  // 码流数据编码格式
    DWORD dwDisplayBufNum;                    // 播放库播放缓冲区最大缓冲帧数
    BYTE  byNPQMode;                          // NPQ模式
    BYTE  byRecvMetaData;                     // 是否接收元数据
    BYTE  byDataType;                         // 数据类型
    BYTE  byRes[69];                          // 保留
} NET_DVR_PREVIEWINFO, *LPNET_DVR_PREVIEWINFO;

/* ========================================================================
 * 数据结构定义 - 设备配置相关（文档未涉及，保留注释供参考）
 * ======================================================================== */

// 注：设备配置相关结构（NET_DVR_DEVICECFG, NET_DVR_PICCFG）已删除
// 当前封装不涉及设备配置功能

/* ========================================================================
 * 数据结构定义 - PTZ相关（文档未涉及，保留注释供参考）
 * ======================================================================== */

// 注：PTZ位置信息结构（NET_DVR_PTZPOS, NET_DVR_PTZSCOPE）已删除
// 当前封装仅使用基本PTZ控制命令，不涉及位置信息获取

/* ========================================================================
 * 数据结构定义 - 其他（文档未涉及，保留注释供参考）
 * ======================================================================== */

// 注：其他辅助结构（如 NET_DVR_TIME）已删除
// 当前封装不涉及时间相关配置

/* ========================================================================
 * 回调函数类型定义
 * ======================================================================== */

// 报警回调函数
typedef void(CALLBACK *MSGCallBack)(LONG lCommand, NET_DVR_ALARMER *pAlarmer, char *pAlarmInfo, DWORD dwBufLen, void* pUser);

// 实时数据回调函数（预览使用，云台控制需要）
typedef void(CALLBACK *REALDATACALLBACK)(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, void *pUser);

// 登录结果回调函数（异步登录）
typedef void(CALLBACK *fLoginResultCallBack)(LONG lUserID, DWORD dwResult, LPNET_DVR_DEVICEINFO_V30 lpDeviceInfo, void *pUser);

/* ========================================================================
 * Go回调函数声明（供CGO使用）
 * ======================================================================== */

// 登录结果回调（异步登录使用）
extern void GoLoginResultCallback(LONG lUserID, DWORD dwResult, LPNET_DVR_DEVICEINFO_V30 lpDeviceInfo, void *pUser);

// 实时数据回调（预览使用，云台控制需要）
extern void GoRealDataCallback(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, uintptr_t handle);

/* ========================================================================
 * SDK函数声明 - 初始化和清理
 * ======================================================================== */

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_Init();                                      // SDK初始化
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_Cleanup();                                   // SDK清理
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_SetConnectTime(DWORD dwWaitTime, DWORD dwTryTimes); // 设置连接超时
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_SetReconnect(DWORD dwInterval, BOOL bEnableRecon);  // 设置重连
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_SetLogToFile(DWORD nLogLevel, char * strLogDir, BOOL bAutoDel); // 设置日志

/* ========================================================================
 * SDK函数声明 - 用户登录
 * ======================================================================== */

// 用户登录（V30版本，兼容旧设备，云台控制常用）
HIKSDK_API LONG HIKSDK_CALL NET_DVR_Login_V30(
    char *sDVRIP,                             // 设备IP地址
    WORD wDVRPort,                           // 设备端口号
    char *sUserName,                         // 登录用户名
    char *sPassword,                         // 登录密码
    LPNET_DVR_DEVICEINFO_V30 lpDeviceInfo    // 设备信息
);

// 用户登录（V40版本，推荐使用）
HIKSDK_API LONG HIKSDK_CALL NET_DVR_Login_V40(
    LPNET_DVR_USER_LOGIN_INFO pLoginInfo,    // 登录参数
    LPNET_DVR_DEVICEINFO_V40 lpDeviceInfo    // 设备信息
);

// 用户登出
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_Logout(LONG lUserID);

// 动态IP解析
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_GetDVRIPByResolveSvr_EX(
    char *sServerIP,                         // 解析服务器地址
    WORD wServerPort,                        // 解析服务器端口
    BYTE *sDVRName,                          // 设备名称
    WORD wDVRNameLen,                        // 设备名称长度
    BYTE *sDVRSerialNumber,                  // 设备序列号
    WORD wDVRSerialLen,                      // 序列号长度
    char* sGetIP,                            // 获取到的IP地址
    DWORD *dwPort                            // 获取到的端口号
);

/* ========================================================================
 * SDK函数声明 - 参数配置（文档未涉及，保留注释供参考）
 * ======================================================================== */

// 注：设备配置函数（NET_DVR_GetDVRConfig, NET_DVR_SetDVRConfig）已删除
// 当前封装不涉及设备配置功能

/* ========================================================================
 * SDK函数声明 - 报警功能
 * ======================================================================== */

// 设置报警消息回调函数
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_SetDVRMessageCallBack_V30(
    MSGCallBack fMessageCallBack,            // 报警回调函数
    void* pUser                              // 用户数据
);

// 建立报警上传通道（布防）
HIKSDK_API LONG HIKSDK_CALL NET_DVR_SetupAlarmChan_V41(
    LONG lUserID,                            // 用户ID
    LPNET_DVR_SETUPALARM_PARAM lpSetupParam  // 布防参数
);

// 关闭报警上传通道
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_CloseAlarmChan_V30(
    LONG lAlarmHandle                        // 报警句柄
);

// 启动监听（接收设备主动上传的报警信息）
HIKSDK_API LONG HIKSDK_CALL NET_DVR_StartListen_V30(
    char *sLocalIP,                          // PC机本地IP地址（可以为NULL）
    WORD wLocalPort,                         // PC本地监听端口号
    MSGCallBack DataCallback,                // 回调函数
    void* pUserData                          // 用户数据
);

// 停止监听
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_StopListen_V30(
    LONG lListenHandle                       // 监听句柄
);

/* ========================================================================
 * SDK函数声明 - 实时预览（云台控制需要预览句柄）
 * ======================================================================== */

// 实时预览（云台控制需要通过此函数获取预览句柄）
HIKSDK_API LONG HIKSDK_CALL NET_DVR_RealPlay_V40(
    LONG lUserID,                            // 用户ID
    LPNET_DVR_PREVIEWINFO lpPreviewInfo,     // 预览参数
    REALDATACALLBACK fRealDataCallBack_V30,  // 数据回调函数
    void* pUser                              // 用户数据
);

// 停止实时预览
HIKSDK_API BOOL HIKSDK_CALL NET_DVR_StopRealPlay(
    LONG lRealHandle                         // 预览句柄
);

/* ========================================================================
 * SDK函数声明 - PTZ控制
 * ======================================================================== */

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZControl(
    LONG lRealHandle,                        // 预览句柄
    DWORD dwPTZCommand,                      // PTZ控制命令
    DWORD dwStop                             // 是否停止：0-开始，1-停止
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZControl_Other(
    LONG lUserID,                            // 用户ID
    LONG lChannel,                           // 通道号
    DWORD dwPTZCommand,                      // PTZ控制命令
    DWORD dwStop                             // 是否停止：0-开始，1-停止
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZControlWithSpeed(
    LONG lRealHandle,                        // 预览句柄
    DWORD dwPTZCommand,                      // PTZ控制命令
    DWORD dwStop,                            // 是否停止：0-开始，1-停止
    DWORD dwSpeed                            // 速度：1-7
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZControlWithSpeed_Other(
    LONG lUserID,                            // 用户ID
    LONG lChannel,                           // 通道号
    DWORD dwPTZCommand,                      // PTZ控制命令
    DWORD dwStop,                            // 是否停止：0-开始，1-停止
    DWORD dwSpeed                            // 速度：1-7
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZPreset(
    LONG lRealHandle,                        // 预览句柄
    DWORD dwPTZPresetCmd,                    // 预置点命令
    DWORD dwPresetIndex                      // 预置点序号
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZPreset_Other(
    LONG lUserID,                            // 用户ID
    LONG lChannel,                           // 通道号
    DWORD dwPTZPresetCmd,                    // 预置点命令
    DWORD dwPresetIndex                      // 预置点序号
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZCruise(
    LONG lRealHandle,                        // 预览句柄
    DWORD dwPTZCruiseCmd,                    // 巡航命令
    BYTE byCruiseRoute,                      // 巡航路径
    BYTE byCruisePoint,                      // 巡航点
    WORD wInput                              // 输入参数
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZCruise_Other(
    LONG lUserID,                            // 用户ID
    LONG lChannel,                           // 通道号
    DWORD dwPTZCruiseCmd,                    // 巡航命令
    BYTE byCruiseRoute,                      // 巡航路径
    BYTE byCruisePoint,                      // 巡航点
    WORD wInput                              // 输入参数
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZTrack(
    LONG lRealHandle,                        // 预览句柄
    DWORD dwPTZTrackCmd                      // 轨迹命令
);

HIKSDK_API BOOL HIKSDK_CALL NET_DVR_PTZTrack_Other(
    LONG lUserID,                            // 用户ID
    LONG lChannel,                           // 通道号
    DWORD dwPTZTrackCmd                      // 轨迹命令
);

/* ========================================================================
 * SDK函数声明 - 其他功能
 * ======================================================================== */

HIKSDK_API DWORD HIKSDK_CALL NET_DVR_GetLastError();     // 获取最后错误码
HIKSDK_API char* HIKSDK_CALL NET_DVR_GetErrorMsg(LONG *pErrorNo); // 获取错误信息
HIKSDK_API DWORD HIKSDK_CALL NET_DVR_GetSDKVersion();    // 获取SDK版本

/* ========================================================================
 * C包装函数（用于简化CGO调用）
 * ======================================================================== */

// 包装的预览函数，自动设置Go回调（云台控制需要）
static inline LONG NET_DVR_RealPlay_V40_WithCallback(LONG lUserID, LPNET_DVR_PREVIEWINFO lpPreviewInfo, uintptr_t handle) {
    return NET_DVR_RealPlay_V40(lUserID, lpPreviewInfo, (REALDATACALLBACK)GoRealDataCallback, (void*)handle);
}

#ifdef __cplusplus
}
#endif

#endif // HIKSDK_WRAPPER_H