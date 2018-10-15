package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type Parameters struct {
	Action           string
	Format           string
	Version          string
	RegionId         string
	Timestamp        string
	SignaName        string
	AccessKeyId      string
	PhoneNumner      string
	AccessSecret     string
	TemplateCode     string
	TemplateParam    string
	SignatureNonce   string
	SignatureMethod  string
	SignatureVersion string
}

type SendMessageResponseSuccess struct {
	Code string `json:"Code"`
	//	BizId     string `json:"BizId"`
	//	Message   string `json:"Message"`
	//	RequestId string `json:"RequestId"`
}

type SendMessageResponseError struct {
	Code string `json:"Code"`
	//	HostId    string `json:"HostId"`
	//	Message   string `json:"Message"`
	//	RequestId string `json:"RequestId"`
	//	Recommend string `json:"Recommend"`
}

func Init(user UserParameters) Parameters {
	return Parameters{
		Action:           "SendSms",
		Format:           "JSON",
		Version:          "2017-05-25",
		RegionId:         "cn-hangzhou",
		Timestamp:        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		SignaName:        user.SignName,
		AccessKeyId:      user.AccessKeyId,
		PhoneNumner:      user.PhoneNumber,
		AccessSecret:     user.AccessSecret,
		TemplateCode:     user.TemplateCode,
		TemplateParam:    fmt.Sprintf("{\"code\":\"%s\"}", user.TemplateParam),
		SignatureNonce:   fmt.Sprintf("%s", uuid.Must(uuid.NewV4())),
		SignatureMethod:  "HMAC-SHA1",
		SignatureVersion: "1.0",
	}
}

func (o *Parameters) Send() (int, string, error) {
	signature, paras := o.GetSignature()

	url := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s%s", signature, paras)
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}

	var response string

	if resp.StatusCode == 200 {
		sendSuc := &SendMessageResponseSuccess{}
		err := json.NewDecoder(resp.Body).Decode(sendSuc)
		if err != nil {
			return 0, "", err
		}
		response = sendSuc.Code
	} else {
		sendErr := &SendMessageResponseError{}
		err := json.NewDecoder(resp.Body).Decode(sendErr)
		if err != nil {
			return 0, "", err
		}
		response = sendErr.Code
	}

	return resp.StatusCode, response, nil
}

func (o *Parameters) GetSignature() (string, string) {
	paras := make(map[string]string)

	paras["Action"] = o.Action
	paras["Format"] = o.Format
	paras["Version"] = o.Version
	paras["RegionId"] = o.RegionId
	paras["SignName"] = o.SignaName
	paras["Timestamp"] = o.Timestamp
	paras["AccessKeyId"] = o.AccessKeyId
	paras["PhoneNumbers"] = o.PhoneNumner
	paras["TemplateCode"] = o.TemplateCode
	paras["TemplateParam"] = o.TemplateParam
	paras["SignatureNonce"] = o.SignatureNonce
	paras["SignatureMethod"] = o.SignatureMethod
	paras["SignatureVersion"] = o.SignatureVersion

	keys := make([]string, 0, 10)

	for k := range paras {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sortString := ""
	for _, k := range keys {
		sortString = sortString + "&" + CharatersReplace(k) + "=" + CharatersReplace(paras[k])
	}

	buildString := "GET" + "&" + CharatersReplace("/") + "&" + CharatersReplace(sortString[1:])

	signature := HmacSha1(o.AccessSecret+"&", buildString)

	return CharatersReplace(signature), sortString
}

func HmacSha1(secret string, buildString string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(buildString))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func CharatersReplace(str string) string {
	str = url.QueryEscape(str)
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "*", "%2A", -1)
	str = strings.Replace(str, "%7E", "~", -1)

	return str
}
