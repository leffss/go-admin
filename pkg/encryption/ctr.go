// Package encryption
// 对称加密 AES算法 + CTR分组模式
// key 长度要求：
// 16 字节 - AES-128
// 24 字节 - AES-192
// 32 字节 - AES-256
// 加密后的密文（[]byte）直接使用 string(密文) 会包含乱码
// 需要使用 hex 转 16 进制 hex.EncodeToString(密文) ,返回
// 的字符串可以以 varchar 类型直接存 mysql 数据库，取出后使
// 用 hex.DecodeString(16进制的密文字符) 得到原始密文
package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/leffss/go-admin/pkg/setting"
)

var aesKey []byte
var privateKey string

func Setup() {
	// aesKey 为程序写死密钥+配置文件设置的密钥
	privateKey = "P[;@#5am"
	aesKey = []byte(privateKey + setting.AppSetting.AesKey)
	//fmt.Println("加密测试 ----------------------------------- start")
	//z := "123()12;]你好了潇洒xzc,a@!$#%^&~P12"
	//fmt.Println("原始字符:", z)
	//x, _ := AesCtrEncryptToString([]byte(z))
	//fmt.Println("加密字符:", x)
	//y, _ := AesCtrStringToDecrypt(x)
	//fmt.Println("解密字符:", string(y))
	//if z == string(y) {
	//	fmt.Println("加密解密成功")
	//} else {
	//	fmt.Println("加密解密失败")
	//}
	//fmt.Println("加密测试 ----------------------------------- end")
}

// AesCtrEncrypt 输入明文，输出密文([]byte)
func AesCtrEncrypt(plainText []byte) ([]byte, error) {
	// 第一步：创建aes密码接口
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	//第二步：创建分组模式ctr
	// iv 要与算法长度一致，16字节
	// 使用bytes.Repeat创建一个切片，长度为blockSize()，16个字符"1"
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	//第三步：加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

// AesCtrEncryptToString 输入明文，输出密文(string)
func AesCtrEncryptToString (plainText []byte) (string, error) {
	encryptData, err := AesCtrEncrypt(plainText)
	if err != nil {
		return "", err
	}
	return EncryptDataToString(encryptData), nil
}

// AesCtrDecrypt 输入密文([]byte)，得到明文
func AesCtrDecrypt(encryptData []byte) ([]byte, error) {
	return AesCtrEncrypt(encryptData)
}

// AesCtrStringToDecrypt 输入密文(string)，得到明文
func AesCtrStringToDecrypt(stringEncryptData string) ([]byte, error) {
	encryptData, err := StringToEncryptData(stringEncryptData)
	if err != nil {
		return nil, err
	}
	return AesCtrDecrypt(encryptData)
}

func EncryptDataToString(encryptData []byte) string {
	return hex.EncodeToString(encryptData)
}

func StringToEncryptData(stringEncryptData string) ([]byte, error) {
	return hex.DecodeString(stringEncryptData)
}
