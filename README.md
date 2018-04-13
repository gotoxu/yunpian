[![Build Status](https://travis-ci.org/gotoxu/yunpian.svg?branch=master)](https://travis-ci.org/gotoxu/yunpian)

# YunPian
云片网Go SDK，暂时只实现了短信发送相关接口

## Installation
```go
go get -u github.com/gotoxu/yunpian
```

## Example usages
```go
package main

import (
	"testing"

	"github.com/gotoxu/assert"
)

func main() {
	sms := NewClient(DefaultConfig().WithAPIKey("you api key").WithUseSSL(true)).SMS()
	input := &SingleSendRequest{
		Mobile: "13800138000",
		Text:   "您的验证码为332211,请不要告诉其他人哦",
	}

	resp, err := sms.SingleSend(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("短信发送是否成功: %t\n", resp.IsSuccess())
}
```
