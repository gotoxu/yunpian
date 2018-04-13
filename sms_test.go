package yunpian

// import (
// 	"testing"

// 	"github.com/gotoxu/assert"
// )

// func TestSingleSend(t *testing.T) {
// 	sms := NewClient(DefaultConfig().WithAPIKey("").WithUseSSL(true)).SMS()
// 	input := &SingleSendRequest{
// 		Mobile: "13320942172",
// 		Text:   "您的验证码为332211,请不要告诉其他人哦",
// 	}

// 	resp, err := sms.SingleSend(input)
// 	assert.Nil(t, err)
// 	assert.True(t, resp.IsSuccess())
// }

// func TestInternationalSingleSend(t *testing.T) {
// 	sms := NewClient(DefaultConfig().WithAPIKey("").WithUseSSL(true)).SMS()
// 	input := &SingleSendRequest{
// 		Mobile: "+19782268687",
// 		Text:   "【BlockCDN】Verification code 332211",
// 	}

// 	resp, err := sms.SingleSend(input)
// 	assert.Nil(t, err)
// 	assert.True(t, resp.IsSuccess())
// }
