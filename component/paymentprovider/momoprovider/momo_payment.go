package momoprovider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"foodlive/common"
	"foodlive/component/paymentprovider"
	"foodlive/config"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
)

type momoProvider struct {
	EndPointMomo  string
	PartnerCode   string
	AccessKey     string
	SecretKey     string
	RequestType   string
	NotifyUrl     string
	BaseReturnUrl string
}

func NewMomoProvider(cf config.MomoConfig) *momoProvider {
	if cf.EndPointMomo == "" || cf.SecretKey == "" || cf.AccessKey == "" || cf.NotifyUrl == "" || cf.BaseReturnUrl == "" || cf.PartnerCode == "" || cf.RequestType == "" {
		log.Fatal("Momo config is empty")
	}
	return &momoProvider{
		EndPointMomo:  cf.EndPointMomo,
		PartnerCode:   cf.PartnerCode,
		RequestType:   cf.RequestType,
		NotifyUrl:     cf.NotifyUrl,
		BaseReturnUrl: cf.BaseReturnUrl,
		SecretKey:     cf.SecretKey,
		AccessKey:     cf.AccessKey,
	}
}

func (m *momoProvider) SendRequestPayment(ctx context.Context, data paymentprovider.OrderRequester, dataExtra string) (*paymentprovider.TransactionResp, error) {
	b := guuid.New()

	var orderId = fmt.Sprintf("%v", data.GetOrderId())
	var requestId = fmt.Sprintf("%v", b)
	var endpoint = m.EndPointMomo
	var partnerCode = m.PartnerCode
	var accessKey = m.AccessKey
	var secretKey = m.SecretKey
	var orderInfo = dataExtra
	var returnUrl = m.BaseReturnUrl + orderId
	var notifyUrl = m.NotifyUrl
	var amount = fmt.Sprintf("%v", data.GetPrice())
	var requestType = m.RequestType
	var extraData = dataExtra

	//build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&returnUrl=")
	rawSignature.WriteString(returnUrl)
	rawSignature.WriteString("&notifyUrl=")
	rawSignature.WriteString(notifyUrl)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(secretKey))

	// Write Data to it
	hmac.Write(rawSignature.Bytes())

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))
	var payload = paymentprovider.TransactionReq{
		PartnerCode: partnerCode,
		AccessKey:   accessKey,
		RequestID:   requestId,
		Amount:      amount,
		OrderID:     orderId,
		OrderInfo:   orderInfo,
		ReturnURL:   returnUrl,
		NotifyURL:   notifyUrl,
		ExtraData:   extraData,
		RequestType: requestType,
		Signature:   signature,
	}

	var jsonPayload []byte
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, common.ErrInvalidRequest(err)
	}
	//send HTTP to momo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	defer resp.Body.Close()
	//decode result
	var respTransaction paymentprovider.TransactionResp
	err = json.NewDecoder(resp.Body).Decode(&respTransaction)
	if err != nil {
		return nil, common.ErrInternal(err)
	}
	return &respTransaction, nil
}
