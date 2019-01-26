package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
)

//pks不全操作
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密参数，定义为常亮
//var iv  = []byte("1234567887654321")

//aes加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	//加密秘钥使用输入文本密码的md5
	md5key:=md5.Sum(key)
	iv:= md5key[:]

	block, err := aes.NewCipher(md5key[:])
	//block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block,iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}


//aes解密
func AesDecrypt(crypted []byte, key []byte) ([]byte, error) {
	defer handler()
	md5key:=md5.Sum(key)
	iv:= md5key[:]

	block, err := aes.NewCipher(md5key[:])
	if err != nil {
		return nil, err
	}

	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//捕获异常
func handler()  {
	recover()
	/*
	if err := recover(); err != nil {

		fmt.Println("recover msg: ", err)

	} else {

		fmt.Println("recover ok")
	}*/
}