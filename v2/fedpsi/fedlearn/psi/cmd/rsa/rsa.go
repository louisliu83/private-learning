package main

import (
	"fmt"

	"github.com/wumansgy/goEncrypt/rsa"
)

func main() {
	msg := "床前明月光，疑是地上霜，举头望明月，低头思故乡"
	rsaBase64Key, err := rsa.GenerateRsaKeyBase64(1024)
	if err != nil {
		fmt.Println(err)
		return
	}
	base64Text, err := rsa.RsaEncryptToBase64([]byte(msg), rsaBase64Key.PublicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("rsa加密后的base64密文为:\n%s\n", base64Text)
	plaintext, err := rsa.RsaDecryptByBase64(base64Text, rsaBase64Key.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("rsa解密后:\n%s\n", string(plaintext))

	fmt.Printf("公钥:\n%s\n", rsaBase64Key.PublicKey)
	fmt.Printf("私钥:\n%s\n", rsaBase64Key.PrivateKey)

}
