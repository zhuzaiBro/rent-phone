package wx_api

import "github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"

type WxSessionKeyDto struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type WxPhoneDto struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

type RequestDto struct {
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
	Code          string `json:"code"`
}

type UserInfoResponce struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
	Id int64 `json:"id"`
}

type PhoneResponce struct {
	Phone       string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

type WxPayRequest struct {
	OrderId string `json:"order_id"`
}

// 微信支付API ========================================================

type WechatPayNotify struct {
	Id           string `json:"id"`
	CreateTime   string `json:"create_time"`
	ResourceType string `json:"resource_type"`
	EventType    string `json:"event_type"`
	Summary      string `json:"summary"`
	Resource     struct {
		OriginalType   string `json:"original_type"`
		Algorithm      string `json:"algorithm"`
		Cipher         string `json:"ciphertext"`
		AssociatedData string `json:"associated_data"`
		Nonce          string `json:"nonce"`
	} `json:"resource"`
}

type Payer struct {
	Openid string `json:"openid"`
}
type WxPayCallback struct {
	TradeState    string       `json:"trade_state"`
	TransIdo      string       `json:"out_trade_no"`
	TransactionId string       `json:"transaction_id"`
	SuccessTime   string       `json:"success_time"`
	Amount        jsapi.Amount `json:"amount"`
	PayUser       Payer        `json:"payer"`
}
