package service

import "testing"

func TestCaptcha_Generate(t *testing.T) {
	captcha := NewCaptcha()
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return
	}
	t.Log(id, b64s, answer)
	captcha.Verify(id, answer, false)
}
