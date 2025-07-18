package service

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	captcha *base64Captcha.Captcha
}

func NewCaptcha() *Captcha {
	// 自定义验证码驱动
	driver := &base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      6,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          6,
		Source:          "1234567890ABCDEFGHJKLMNPQRSTUVWXYZ",
		BgColor:         &color.RGBA{R: 240, G: 240, B: 246, A: 255},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	driver = driver.ConvertFonts()

	return &Captcha{
		captcha: base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore),
	}
}

// Generate 生成验证码
func (s *Captcha) Generate() (id, b64s, answer string, err error) {
	return s.captcha.Generate()
}

// Verify 验证验证码
func (s *Captcha) Verify(id, answer string, clear bool) bool {
	return s.captcha.Store.Verify(id, answer, clear)
}
