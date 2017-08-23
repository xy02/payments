package payments

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	//	alRUL           = "https://openapi.alipaydev.com/gateway.do"
	alRUL          = "https://openapi.alipay.com/gateway.do"
	alCreateMethod = "alipay.trade.precreate"
	alQueryMethod  = "alipay.trade.query"
	alCharset      = "utf-8"
	alSignType     = "RSA"
)

type (
	//AlipayClient 阿里支付客户端
	AlipayClient struct {
		Key       *rsa.PrivateKey //用于发送数据的私钥
		PubKey    *rsa.PublicKey  //阿里账号的公钥
		AppID     string
		NotifyURL string //回调地址
	}
	//ALCreateContent 阿里创建预付的内容
	ALCreateContent struct {
		OutTradeNo  string `json:"out_trade_no"`
		TotalAmount string `json:"total_amount"`
		Subject     string `json:"subject"`
		// TimeoutExpress string `json:"timeout_express"`
	}
	//ALCreateRequest 阿里创建预付的包装数据
	ALCreateRequest struct {
		AppID      string `json:"app_id"`
		Method     string `json:"method"`
		Charset    string `json:"charset"`
		SignType   string `json:"sign_type"`
		Sign       string `json:"sign"`
		Timestamp  string `json:"timestamp"`
		NotifyURL  string `json:"notify_url"`
		BizContent string `json:"biz_content"`
	}
	//ALCreateResult 阿里创建预付的结果
	ALCreateResult struct {
		AlipayTradePrecreateResponse json.RawMessage `json:"alipay_trade_precreate_response"`
		Sign                         string          `json:"sign"`
	}
	//AlipayTradePrecreateResponse 阿里创建预付的结果中的应答
	AlipayTradePrecreateResponse struct {
		Code       string `json:"code" bson:"code"`
		Msg        string `json:"msg" bson:"msg"`
		SubCode    string `json:"sub_code" bson:"sub_code"`
		SubMsg     string `json:"sub_msg" bson:"sub_msg"`
		OutTradeNo string `json:"out_trade_no" bson:"out_trade_no"`
		QrCode     string `json:"qr_code" bson:"qr_code"`
	}
	//ALQueryContent 查询交易的内容
	ALQueryContent struct {
		OutTradeNo string `json:"out_trade_no"`
	}
	//ALQueryRequest 阿里查询预付的包装数据
	ALQueryRequest struct {
		AppID      string `json:"app_id"`
		Method     string `json:"method"`
		Charset    string `json:"charset"`
		SignType   string `json:"sign_type"`
		Sign       string `json:"sign"`
		Timestamp  string `json:"timestamp"`
		BizContent string `json:"biz_content"`
	}
	//ALQueryResult 查询交易的结果
	ALQueryResult struct {
		AlipayTradeQueryResponse json.RawMessage `json:"alipay_trade_query_response"`
		Sign                     string          `json:"sign"`
	}
	//AlipayTradeQueryResponse 查询结果中的应答
	//TradeStatus交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
	AlipayTradeQueryResponse struct {
		Code           string          `bson:"code" json:"code"`
		Msg            string          `bson:"msg" json:"msg"`
		SubCode        string          `bson:"sub_code" json:"sub_code"`
		SubMsg         string          `bson:"sub_msg" json:"sub_msg"`
		TradeNo        string          `bson:"trade_no" json:"trade_no"`
		OutTradeNo     string          `bson:"out_trade_no" json:"out_trade_no"`
		OpenID         string          `bson:"open_id" json:"open_id"`
		BuyerLogonID   string          `bson:"buyer_logon_id" json:"buyer_logon_id"`
		TradeStatus    string          `bson:"trade_status" json:"trade_status"`
		TotalAmount    string          `bson:"total_amount" json:"total_amount"`
		ReceiptAmount  string          `bson:"receipt_amount" json:"receipt_amount"`
		BuyerPayAmount string          `bson:"buyer_pay_amount" json:"buyer_pay_amount"`
		PointAmount    string          `bson:"point_amount" json:"point_amount"`
		InvoiceAmount  string          `bson:"invoice_amount" json:"invoice_amount"`
		SendPayDate    string          `bson:"send_pay_date" json:"send_pay_date"`
		AlipayStoreID  string          `bson:"alipay_store_id" json:"alipay_store_id"`
		StoreID        string          `bson:"store_id" json:"store_id"`
		TerminalID     string          `bson:"terminal_id" json:"terminal_id"`
		StoreName      string          `bson:"store_name" json:"store_name"`
		BuyerUserID    string          `bson:"buyer_user_id" json:"buyer_user_id"`
		FundBillList   []TradeFundBill `bson:"fund_bill_list" json:"fund_bill_list"`
		//		Discount_goods_detail           string `json:"discount_goods_detail"`
		//		Industry_sepc_detail         string `json:"industry_sepc_detail"`
	}
	//TradeFundBill ...
	TradeFundBill struct {
		Amount      string `json:"amount"`
		FundChannel string `json:"fund_channel"`
		RealAmount  string `json:"real_amount"`
	}
)

//CreateTrade 创建交易
func (client *AlipayClient) CreateTrade(content CreateContent) (Trade, error) {
	yuan := content.TotalAmount / 100
	fen := content.TotalAmount % 100
	totalAmount := fmt.Sprintf("%v.%v", yuan, fen)
	if fen < 10 {
		totalAmount = fmt.Sprintf("%v.0%v", yuan, fen)
	}
	bizContent, err := json.Marshal(ALCreateContent{
		OutTradeNo:  content.OutTradeNo,
		TotalAmount: totalAmount,
		Subject:     content.Subject,
	})
	if err != nil {
		return nil, err
	}
	data := &ALCreateRequest{
		AppID:      client.AppID,
		Method:     alCreateMethod,
		Charset:    alCharset,
		SignType:   alSignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		BizContent: string(bizContent),
		NotifyURL:  client.NotifyURL,
	}
	result := &ALCreateResult{}
	err = client.request(data, result)
	if err != nil {
		return nil, err
	}
	//check sign
	err = client.verify(result.AlipayTradePrecreateResponse, result.Sign)
	if err != nil {
		return nil, err
	}
	response := &AlipayTradePrecreateResponse{}
	err = json.Unmarshal(result.AlipayTradePrecreateResponse, response)
	if err != nil {
		return nil, err
	}
	if response.Code != "10000" {
		return nil, errors.New(response.SubMsg)
	}
	return response, nil
}

//QueryTrade 查询交易,返回交易状态
func (client *AlipayClient) QueryTrade(outTradeNo string) (Trade, error) {
	bizContent, err := json.Marshal(ALQueryContent{
		OutTradeNo: outTradeNo,
	})
	if err != nil {
		return nil, err
	}
	data := &ALQueryRequest{
		AppID:      client.AppID,
		Method:     alQueryMethod,
		Charset:    alCharset,
		SignType:   alSignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		BizContent: string(bizContent),
	}
	result := &ALQueryResult{}
	err = client.request(data, result)
	if err != nil {
		return nil, err
	}
	//check sign
	err = client.verify(result.AlipayTradeQueryResponse, result.Sign)
	if err != nil {
		return nil, err
	}
	response := new(AlipayTradeQueryResponse)
	err = json.Unmarshal(result.AlipayTradeQueryResponse, response)
	if err != nil {
		return nil, err
	}
	if response.Code != "10000" {
		return nil, errors.New(response.SubMsg)
	}
	return response, nil
}

//request 向阿里服务器发起请求
func (client *AlipayClient) request(data interface{}, result interface{}) error {
	client.sign(data)
	query := alMarshal(data)
	url := fmt.Sprintf("%s?%s", alRUL, query)
	// fmt.Printf("request: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("response: %s\n", body)
	return json.Unmarshal(body, result)
}

//alMarshal 把数据序列化成url格式
func alMarshal(data interface{}) string {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	d := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		var name = t.Field(i).Tag.Get("json")
		var value = v.Field(i).Interface()
		if v, ok := value.(string); ok && v != "" {
			d.Set(name, v)
		}
	}
	return d.Encode()
}

//sign 对要发送的数据按阿里规则签名
func (client *AlipayClient) sign(data interface{}) error {
	t := reflect.TypeOf(data).Elem()
	v := reflect.ValueOf(data).Elem()
	pairs := make([]string, 0, 64)
	for i := 0; i < t.NumField(); i++ {
		var name = t.Field(i).Tag.Get("json")
		var value = v.Field(i).Interface()
		if v, ok := value.(string); ok && v != "" && name != "sign" {
			pairs = append(pairs, name+"="+v)
		}
	}
	sort.Strings(pairs)
	var str = strings.Join(pairs, "&")
	//	fmt.Println(len(pairs), cap(pairs), str)
	hashed := sha1.Sum([]byte(str))
	s, err := rsa.SignPKCS1v15(rand.Reader, client.Key, crypto.SHA1, hashed[:])
	if err != nil {
		return err
	}
	v.FieldByName("Sign").SetString(base64.StdEncoding.EncodeToString(s))
	return nil
}

//verify 验证应答的签名
func (client *AlipayClient) verify(data []byte, sign string) error {
	hashed := sha1.Sum(data)
	s, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(client.PubKey, crypto.SHA1, hashed[:], s)
}
