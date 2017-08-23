package payments

import "github.com/xy02/utils"

//Client 支付客户端
type Client interface {
	CreateTrade(content CreateContent) (Trade, error)
}

//CreateContent 创建交易需要的内容
type CreateContent struct {
	OutTradeNo  string //交易号
	Subject     string //付款主题
	TotalAmount int64  //总金额
}

//Trade 交易
type Trade interface {
	isTrade()
	GetStatus() Status
}

//Payment 支付方式
type Payment int32

//Status 支付状态
type Status int32

const (
	//Alipay 阿里统一收单线下交易预创建
	Alipay Payment = iota
	//Weixin 微信扫码支付
	Weixin
)

//Unkown 未知支付状态
const (
	Unknown   Status = iota
	Precreate        //预支付
	Wait
	Success
	Finished
	Closed
)

//GetClient 获取支付客户端
func GetClient(payment Payment) (Client, error) {
	client := paymentsMap[payment]
	if client == nil {
		return nil, errInvalidPayment
	}
	return client, nil
}

//AlipayConfig 支付宝配置
type AlipayConfig struct {
	PriKey       string
	PriKeyPwd    string
	AlipayPubKey string
	AlipayAppID  string
	NotifyURL    string
}

//ConfigAlipayClient 配置支付宝客户端
func ConfigAlipayClient(config AlipayConfig) error {
	pk, err := utils.DecryptPrivateKey(config.PriKey, []byte(config.PriKeyPwd))
	if err != nil {
		return err
	}
	pubKey, err := utils.RetrievePublicKey(config.AlipayPubKey)
	if err != nil {
		return err
	}
	alipayClient = &AlipayClient{
		Key:       pk,
		PubKey:    pubKey,
		AppID:     config.AlipayAppID,
		NotifyURL: config.NotifyURL,
	}
	return nil
}

var alipayClient *AlipayClient
var paymentsMap = map[Payment]Client{
	Alipay: alipayClient,
}
