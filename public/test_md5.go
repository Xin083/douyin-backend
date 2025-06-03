package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func MD5(params string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(params))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func Base64Md5(params string) string {
	return MD5(base64.StdEncoding.EncodeToString([]byte(params)))
}

func main() {
	// 测试密码
	passwords := []string{
		"123456",
		"11111111",
		"11111111111",
		"88888888",
		"12345678",
	}

	fmt.Println("测试 Base64Md5 加密结果：")
	fmt.Println("----------------------------------------")
	for _, pwd := range passwords {
		encrypted := Base64Md5(pwd)
		fmt.Printf("密码: %s\n加密后: %s\n", pwd, encrypted)
		fmt.Println("----------------------------------------")
	}
}
