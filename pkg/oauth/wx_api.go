package oauth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"rentServer/pkg/config"
	. "rentServer/pkg/model/wx_api"
	"time"
)

func GetPhone(req RequestDto) (phone []byte, err error) {
	resp, err := Code2Session(req.Code)
	if err != nil {
		return
	}
	key, _ := base64.StdEncoding.DecodeString(resp.SessionKey)
	iv, _ := base64.StdEncoding.DecodeString(req.Iv)
	ciphertext, _ := base64.StdEncoding.DecodeString(req.EncryptedData)
	plaintext := make([]byte, len(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		// panic(err)
		return
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	phone = PKCS7UnPadding(plaintext)

	return
}

func GetUserProfile(req RequestDto) (res UserInfoResponce, err error) {

	log.Println(req)

	sessionKeyDto, err := Code2Session(req.Code)
	if err != nil {
		println(err.Error())

	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKeyDto.SessionKey)
	if err != nil {
		println(err.Error())

	}
	cipherText, err := base64.StdEncoding.DecodeString(req.EncryptedData)
	if err != nil {
		println(err.Error())
	}
	ivBytes, err := base64.StdEncoding.DecodeString(req.Iv)
	if err != nil {
		println(err.Error())
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		println(err.Error())
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = Pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return
	}

	err = json.Unmarshal(cipherText, &res)
	if err != nil {
		return
	}
	res.OpenID = sessionKeyDto.OpenId
	res.UnionID = sessionKeyDto.UnionId

	return
}

func Code2Session(code string) (sessionKeyDto WxSessionKeyDto, e error) {
	var (
		url = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	)

	conf := config.GetConfig()

	httpState, bytesData := Get(fmt.Sprintf(url, conf.WechatAppConfiguration.AppId, conf.WechatAppConfiguration.AppSecret, code))
	if httpState != 200 {
		// print("获取sessionKey失败,HTTP CODE:%d", httpState)
		e = errors.New("获取sessionKey失败")
		return
	}
	e = json.Unmarshal(bytesData, &sessionKeyDto)
	if e != nil {
		print(e.Error())
	}
	print(sessionKeyDto.SessionKey, sessionKeyDto.OpenId, sessionKeyDto.UnionId)
	return sessionKeyDto, nil
}

func Get(url string) (int, []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		print("error sending GET request, url: %s, %q", url, err)
		return http.StatusInternalServerError, nil
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			print("error decoding response from GET request, url: %s, %q", url, err)
		}
	}
	return resp.StatusCode, result.Bytes()
}

func DecryptPhoneData(phoneData, sessionKey, iv string) (string, error) {
	decrypt, err := AesDecrypt(phoneData, sessionKey, iv)
	print(decrypt)
	if err != nil {
		print("解密数据失败", err.Error())
		return "", err
	}
	var phoneDto = WxPhoneDto{}
	err = json.Unmarshal(decrypt, &phoneDto)
	if err != nil {
		print("解析手机号信息失败", err)
		return "", err
	}
	var phone = phoneDto.PurePhoneNumber
	return phone, nil
}

func AesDecrypt(encryptedData, sessionKey, iv string) ([]byte, error) {
	//Base64解码
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		print("解密数据")
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	cryptData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	origData := make([]byte, len(cryptData))
	//AES
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		println("aes出错")
		return nil, err
	}
	//CBC
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	//解密
	mode.CryptBlocks(origData, cryptData)
	//去除填充位
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}

func Pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	// fmt.Println("kaishi", data, blockSize)
	var (
		ErrInvalidBlockSize    = errors.New("invalid block size")
		ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
		ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
	)

	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		fmt.Print("n == 0 || n > len(data)")
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

//
//func WxPay(c *core.Context) {
//	println("wxpay")
//	var req WxPayRequest
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
//		return
//	}
//	payData,  err := service.JsApi(&req)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	c.JSONOK(payData)
//}
