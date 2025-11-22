package utils

import (
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// Strcpy 将Go字符串安全地复制到C字符数组
// 自动处理字符串长度截断和空终止符，避免缓冲区溢出
// 参数：
//   - dst: 目标C字符数组的指针
//   - src: 源Go字符串
//   - dstLen: 目标数组的长度（包括空终止符）
//
// 注意：
//   - 如果源字符串长度超过 dstLen-1，将被截断
//   - 自动添加空终止符
func Strcpy(dst unsafe.Pointer, src string, dstLen int) {
	if dstLen <= 0 {
		return // 无效的目标长度
	}

	srcBytes := []byte(src)
	copyLen := len(srcBytes)

	// 确保不会溢出，留出空间给空终止符
	if copyLen >= dstLen {
		copyLen = dstLen - 1
	}

	// 复制数据
	dstSlice := (*[1 << 30]byte)(dst)
	copy(dstSlice[:copyLen], srcBytes[:copyLen])

	// 添加空终止符
	dstSlice[copyLen] = 0
}

// GBKToUTF8 将GBK编码转换为UTF-8
// 用于处理设备返回的中文字符串
// 参数：
//   - b: GBK编码的字节数组
//
// 返回值：
//   - string: UTF-8编码的字符串
//   - error: 转换错误
func GBKToUTF8(b []byte) (string, error) {
	r, err := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	return string(r), err
}

// UTF8ToGBK 将UTF-8编码转换为GBK
// 用于向设备发送中文字符串
// 参数：
//   - s: UTF-8编码的字符串
//
// 返回值：
//   - []byte: GBK编码的字节数组
//   - error: 转换错误
func UTF8ToGBK(s string) ([]byte, error) {
	return simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
}
