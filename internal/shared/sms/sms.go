package sms

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	txSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/algorithm9/flash-deal/internal/model"
)

type SMS struct {
	client     *txSMS.Client
	sdkAppID   string
	signName   string
	templateID string
}

func NewSMSClient(cfg *model.SMS) (*SMS, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	client, err := txSMS.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, err
	}
	return &SMS{client: client, sdkAppID: cfg.SdkAppID, signName: cfg.SignName, templateID: cfg.TemplateID}, nil
}

func (s *SMS) Send(phone string, code string) error {
	req := txSMS.NewSendSmsRequest()
	req.PhoneNumberSet = []*string{&phone}
	req.SmsSdkAppId = &s.sdkAppID
	req.SignName = &s.signName
	req.TemplateId = &s.templateID
	req.TemplateParamSet = []*string{&code}

	_, err := s.client.SendSms(req)
	if tErr, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("TencentCloud SDK Error: %v", tErr)
	}
	return err
}
