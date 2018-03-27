package util

import (
    "io"
    "crypto/aes"
    "crypto/rsa"
    "crypto/x509"
    "crypto/cipher"
    "crypto/rand"
    "encoding/pem"
    "encoding/base64"
)

func encodeBase64(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
    data, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        panic(err)
    }
    return data
}

func encrypt(key, text []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    b := encodeBase64(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext
}

func decrypt(key, text []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    if len(text) < aes.BlockSize {
        panic("ciphertext too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    return decodeBase64(string(text))
}

func Encrypt(text, key []byte) (b []byte) {
    if block, _ := pem.Decode(key); block != nil {
        if pi,err:=x509.ParsePKIXPublicKey(block.Bytes);err==nil{
            if pub, ok := pi.(*rsa.PublicKey); ok && pub != nil {
                b, _ = rsa.EncryptPKCS1v15(rand.Reader, pub, text)
            }            
        }
    }
    return
}

func Decrypt(cipher, key []byte) (b []byte) {
    if block, _ := pem.Decode(key); block != nil {
        if priv,err:=x509.ParsePKCS1PrivateKey(block.Bytes);err==nil{
            b, _ = rsa.DecryptPKCS1v15(rand.Reader, priv, cipher)
        }
    }
    return   
}

/*
"io/ioutil"
func main() {
    if b, err := ioutil.ReadFile("service.conf.bat"); err == nil {
        b = encrypt(key, b)
        ioutil.WriteFile("service.conf", b, 0644)
    }   
}
*/