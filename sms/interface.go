package sms

import "errors"

type UserParameters struct {
	SignName      string
	PhoneNumber   string
	AccessKeyId   string
	TemplateCode  string
	AccessSecret  string
	TemplateParam string
}

func SendMessage(userParas UserParameters) (int, string, error) {
	err := userParas.CheckFields()
	if err != nil {
		return 0, "", err
	}

	instance := Init(userParas)

	return instance.Send()
}

func (o *UserParameters) CheckFields() error {
	if o.SignName == "" {
		return errors.New("`signName` is required")
	}

	if o.PhoneNumber == "" {
		return errors.New("`phoneNumber` is required")
	}

	if o.AccessKeyId == "" {
		return errors.New("`accessKeyId` is required")
	}

	if o.TemplateCode == "" {
		return errors.New("`templateCode` is required")
	}

	if o.AccessSecret == "" {
		return errors.New("`accessSecret` is required")
	}

	if o.TemplateParam == "" {
		return errors.New("`templateParam` is required")
	}

	return nil
}
