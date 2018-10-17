# aliyunSms

## 阿里云短信服务golang版本
### 示例
```
go get "github.com/vfxh/aliyunSms"
```

```
args := aliyunSms.UserParameters{
	SignName:      "XXXX",
	PhoneNumber:   "13700000000",
	AccessKeyId:   "accessKeyId",
	TemplateCode:  "SMS_00000000",
	AccessSecret:  "AccessSecret",
	TemplateParam: "param",																	
}

code, message, err := aliyunSms.SendMessage(args)
if err != nil {
	fmt.Println(err)
	return						
}
```
