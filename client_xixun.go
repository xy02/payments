package payments

// import (
// 	"bytes"
// 	"crypto/md5"
// 	"encoding/xml"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"reflect"
// 	"sort"
// 	"strings"
// 	"time"
// )

// const (
// 	wx_query_url    = "https://api.mch.weixin.qq.com/pay/orderquery"
// 	wx_create_url   = "https://api.mch.weixin.qq.com/pay/unifiedorder"
// 	wx_body         = "aips上网卡扫码充值"
// 	wx_trade_type   = "NATIVE"
// 	wx_content_type = "application/xml"
// )

// type (
// 	WX_QueryData struct {
// 		XMLName      xml.Name `xml:"xml"`
// 		Appid        string   `xml:"appid"`
// 		Mch_id       string   `xml:"mch_id"`
// 		Nonce_str    string   `xml:"nonce_str"`
// 		Out_trade_no string   `xml:"out_trade_no"`
// 		Sign         string   `xml:"sign"`
// 	}

// 	WX_QueryResult struct {
// 		Return_code string `xml:"return_code"`
// 		Return_msg  string `xml:"return_msg"`
// 		Appid       string `xml:"appid"`
// 		Mch_id      string `xml:"mch_id"`
// 		Nonce_str   string `xml:"nonce_str"`
// 		Sign        string `xml:"sign"`
// 		Result_code string `xml:"result_code"`
// 		//when something worng
// 		Err_code         string `xml:"err_code"`
// 		Err_code_des     string `xml:"err_code_des"`
// 		Device_info      string `xml:"device_info"`
// 		Openid           string `xml:"openid"`
// 		Is_subscribe     string `xml:"is_subscribe"`
// 		Trade_type       string `xml:"trade_type"`
// 		Trade_state      string `xml:"trade_state"`
// 		Bank_type        string `xml:"bank_type"`
// 		Total_fee        string `xml:"total_fee"`
// 		Fee_type         string `xml:"fee_type"`
// 		Transaction_id   string `xml:"transaction_id"`
// 		Out_trade_no     string `xml:"out_trade_no"`
// 		Attach           string `xml:"attach"`
// 		Time_end         string `xml:"time_end"`
// 		Trade_state_desc string `xml:"trade_state_desc"`
// 		Cash_fee         string `xml:"cash_fee"`
// 		Cash_fee_type    string `xml:"cash_fee_type"`
// 	}

// 	WX_CreateData struct {
// 		XMLName          xml.Name `xml:"xml"`
// 		Appid            string   `xml:"appid"`
// 		Mch_id           string   `xml:"mch_id"`
// 		Nonce_str        string   `xml:"nonce_str"`
// 		Out_trade_no     string   `xml:"out_trade_no"`
// 		Sign             string   `xml:"sign"`
// 		Body             string   `xml:"body"`
// 		Total_fee        string   `xml:"total_fee"`
// 		Spbill_create_ip string   `xml:"spbill_create_ip"`
// 		Trade_type       string   `xml:"trade_type"`
// 		Notify_url       string   `xml:"notify_url"`
// 	}
// 	WX_CreateResult struct {
// 		Return_code string `xml:"return_code"`
// 		Return_msg  string `xml:"return_msg"`
// 		Appid       string `xml:"appid"`
// 		Mch_id      string `xml:"mch_id"`
// 		Device_info string `xml:"device_info"`
// 		Nonce_str   string `xml:"nonce_str"`
// 		Sign        string `xml:"sign"`
// 		Result_code string `xml:"result_code"`
// 		Trade_type  string `xml:"trade_type"`
// 		Prepay_id   string `xml:"prepay_id"`
// 		Code_url    string `xml:"code_url"`
// 		//when something worng
// 		Err_code     string `xml:"err_code"`
// 		Err_code_des string `xml:"err_code_des"`
// 	}
// 	WeiXinPayment struct {
// 		Key string
// 		//公众号APPID
// 		Appid string
// 		//微信支付商户号
// 		Mch_id string
// 		//，Native支付填调用微信支付API的机器IP。
// 		Spbill_create_ip string
// 		//回调地址
// 		Notify_url string
// 	}
// )

// func (r *WX_QueryData) SetSign(sign string) {
// 	r.Sign = sign
// }
// func (r *WX_CreateData) SetSign(sign string) {
// 	r.Sign = sign
// }

// func (r *WX_QueryResult) CheckSign(sign string) bool {
// 	return r.Sign == sign
// }
// func (r *WX_CreateResult) CheckSign(sign string) bool {
// 	return r.Sign == sign
// }

// func (r WeiXinPayment) Query(out_trade_no string) (string, error) {
// 	data := WX_QueryData{
// 		Appid:        r.Appid,
// 		Mch_id:       r.Mch_id,
// 		Nonce_str:    fmt.Sprintf("%d", time.Now().UnixNano()),
// 		Out_trade_no: out_trade_no,
// 	}
// 	result := WX_QueryResult{}
// 	err := r.request(wx_query_url, &data, &result)
// 	if err != nil {
// 		return "", err
// 	} else {
// 		return result.Trade_state, nil
// 	}
// }

// func (r WeiXinPayment) Create(out_trade_no string, amount int64) (string, error) {
// 	data := WX_CreateData{
// 		Appid:            r.Appid,
// 		Mch_id:           r.Mch_id,
// 		Nonce_str:        fmt.Sprintf("%d", time.Now().UnixNano()),
// 		Out_trade_no:     out_trade_no,
// 		Trade_type:       wx_trade_type,
// 		Body:             wx_body,
// 		Spbill_create_ip: r.Spbill_create_ip,
// 		Notify_url:       r.Notify_url,
// 		Total_fee:        fmt.Sprintf("%v", amount),
// 	}
// 	result := WX_CreateResult{}
// 	err := r.request(wx_create_url, &data, &result)
// 	if err != nil {
// 		return "", err
// 	} else {
// 		return result.Code_url, nil
// 	}
// }

// func (r WeiXinPayment) request(url string, data RequestData, result ResponseResult) error {
// 	data.SetSign(r.sign(data))
// 	buf, err := xml.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	//	fmt.Printf("request: %s\n", buf)
// 	resp, err := http.Post(url, wx_content_type, bytes.NewReader(buf))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	//	fmt.Printf("response: %s\n", body)
// 	err = xml.Unmarshal(body, result)
// 	if err != nil {
// 		return err
// 	}
// 	var sign = r.sign(result)
// 	if result.CheckSign(sign) {
// 		//		fmt.Println("sign ok")
// 		return nil
// 	} else {
// 		return errors.New(fmt.Sprintf("wrong res sign(weixin), url:%s, buf:%s, body:%s", url, buf, body))
// 	}
// }

// func (r WeiXinPayment) sign(data interface{}) string {
// 	t := reflect.TypeOf(data)
// 	v := reflect.ValueOf(data)
// 	if v.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 		v = v.Elem()
// 	}
// 	pairs := make([]string, 0, 64)
// 	for i := 0; i < t.NumField(); i++ {
// 		var name = t.Field(i).Tag.Get("xml")
// 		var value = v.Field(i).Interface()
// 		if v, ok := value.(string); ok && v != "" && name != "sign" {
// 			pairs = append(pairs, name+"="+v)
// 		}
// 	}
// 	sort.Strings(pairs)
// 	var str = strings.Join(pairs, "&") + "&key=" + r.Key
// 	//	fmt.Println(len(pairs), cap(pairs), str)
// 	var b = md5.Sum([]byte(str))
// 	return fmt.Sprintf("%X", b)
// }
