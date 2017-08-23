package payments

func (*AlipayTradePrecreateResponse) isTrade() {}
func (*AlipayTradePrecreateResponse) GetStatus() Status {
	return Precreate
}

func (*AlipayTradeQueryResponse) isTrade() {}
func (response *AlipayTradeQueryResponse) GetStatus() Status {
	switch response.TradeStatus {
	case "WAIT_BUYER_PAY":
		return Wait
	case "TRADE_CLOSED":
		return Closed
	case "TRADE_SUCCESS":
		return Success
	case "TRADE_FINISHED":
		return Finished
	default:
		return Unknown
	}
}

// //ToPB 转换为protocol buffer 定义的消息
// func (r *AlipayTradePrecreateResponse) ToPB() *pb.AlipayTrade {
// 	return &pb.AlipayTrade{
// 		Code:       r.Code,
// 		Msg:        r.Msg,
// 		SubCode:    r.Code,
// 		SubMsg:     r.SubMsg,
// 		OutTradeNo: r.OutTradeNo,
// 		QrCode:     r.QrCode,
// 	}
// }

// //ToPB 转换为protocol buffer 定义的消息
// func (r *AlipayTradeQueryResponse) ToPB() {
// 	return &pb.AlipayTrade{
// 		Code:           r.Code,
// 		Msg:            r.Msg,
// 		SubCode:        r.Code,
// 		SubMsg:         r.SubMsg,
// 		OutTradeNo:     r.OutTradeNo,
// 		TradeNo:        r.TradeNo,
// 		OpenId:         r.OpenID,
// 		BuyerLogonId:   r.BuyerLogonID,
// 		TradeStatus:    r.TradeStatus,
// 		TotalAmount:    r.TotalAmount,
// 		ReceiptAmount:  r.ReceiptAmount,
// 		BuyerPayAmount: r.BuyerPayAmount,
// 		PointAmount:    r.PointAmount,
// 		InvoiceAmount:  r.InvoiceAmount,
// 		SendPayDate:    r.SendPayDate,
// 		AlipayStoreId:  r.AlipayStoreID,
// 		StoreId:        r.StoreID,
// 		TerminalId:     r.TerminalID,
// 		StoreName:      r.StoreName,
// 		BuyerUserId:    r.BuyerUserID,
// 	}
// }
