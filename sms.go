package yunpian

import (
	"errors"
	"net/http"
)

// SMS 是短信发送客户端
type SMS struct {
	*Client
}

// SMS 返回一个新的短信客户端
func (c *Client) SMS() *SMS {
	return &SMS{c}
}

// SingleSendRequest 单条短信发送请求
type SingleSendRequest struct {
	Mobile      string `url:"mobile,omitempty"`
	Text        string `url:"text,omitempty"`
	Extend      string `url:"extend,omitempty"`
	UID         string `url:"uid,omitempty"`
	CallbackURL string `url:"callback_url,omitempty"`
	Register    bool   `url:"register,omitempty"`
}

// Verify 用于验证请求参数的正确性
func (r *SingleSendRequest) Verify() error {
	if len(r.Mobile) == 0 {
		return errors.New("Miss param: mobile")
	}
	if len(r.Text) == 0 {
		return errors.New("Miss param: text")
	}
	return nil
}

// SingleSendResponse 单条短信发送响应
type SingleSendResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Count   int     `json:"count"`
	Fee     float64 `json:"fee"`
	Unit    string  `json:"unit"`
	Mobile  string  `json:"mobile"`
	SID     int64   `json:"sid"`
}

// IsSuccess 用于验证短信发送是否成功
func (resp *SingleSendResponse) IsSuccess() bool {
	return resp.Code == 0
}

// SingleSend 发送单条短信
func (sms *SMS) SingleSend(input *SingleSendRequest) (*SingleSendResponse, error) {
	if input == nil {
		input = &SingleSendRequest{}
	}

	var result SingleSendResponse
	err := sms.sendRequest(http.MethodPost, "/v2/sms/single_send.json", input, &result)
	return &result, err
}

// BatchSendRequest 批量发送请求参数
type BatchSendRequest struct {
	Mobile      string `url:"mobile,omitempty"`
	Text        string `url:"text,omitempty"`
	Extend      string `url:"extend,omitempty"`
	CallbackURL string `url:"callback_url,omitempty"`
}

// Verify 用于检查请求参数的正确性
func (req *BatchSendRequest) Verify() error {
	if len(req.Mobile) == 0 {
		return errors.New("Miss param: mobile")
	}
	if len(req.Text) == 0 {
		return errors.New("Miss param: text")
	}

	return nil
}

// BatchSendResponse 批量发送响应结构
type BatchSendResponse struct {
	TotalCount int                  `json:"total_count"`
	TotalFee   string               `json:"total_fee"`
	Unit       string               `json:"unit"`
	Data       []SingleSendResponse `json:"data"`
}

// BatchSend 批量发送接口
func (sms *SMS) BatchSend(input *BatchSendRequest) (*BatchSendResponse, error) {
	if input == nil {
		input = &BatchSendRequest{}
	}

	var result BatchSendResponse
	err := sms.sendRequest(http.MethodPost, "/v2/sms/batch_send.json", input, &result)
	return &result, err
}

// MultiSendRequest 批量个性化发送请求参数
type MultiSendRequest struct {
	Mobile      string `url:"mobile,omitempty"`
	Text        string `url:"text,omitempty"`
	Extend      string `url:"extend,omitempty"`
	CallbackURL string `url:"callback_url,omitempty"`
}

// Verify 用于检查请求参数的有效性
func (req *MultiSendRequest) Verify() error {
	if len(req.Mobile) == 0 {
		return errors.New("Miss param: mobile")
	}
	if len(req.Text) == 0 {
		return errors.New("Miss param: text")
	}

	return nil
}

// MultiSendResponse 批量个性化发送响应
type MultiSendResponse struct {
	TotalCount int                  `json:"total_count"`
	TotalFee   string               `json:"total_fee"`
	Unit       string               `json:"unit"`
	Data       []SingleSendResponse `json:"data"`
}

// MultiSend 批量个性化发送接口
func (sms *SMS) MultiSend(input *MultiSendRequest) (*MultiSendResponse, error) {
	if input == nil {
		input = &MultiSendRequest{}
	}

	var result MultiSendResponse
	err := sms.sendRequest(http.MethodPost, "/v2/sms/multi_send.json", input, &result)
	return &result, err
}

// TPLSingleSendRequest 指定模板单发请求参数
type TPLSingleSendRequest struct {
	Mobile   string `url:"mobile,omitempty"`
	TPLID    int64  `url:"tpl_id,omitempty"`
	TPLValue string `url:"tpl_value,omitempty"`
	Extend   string `url:"extend,omitempty"`
	UID      string `url:"uid,omitempty"`
}

// Verify 用于检查请求参数的有效性
func (req *TPLSingleSendRequest) Verify() error {
	if len(req.Mobile) == 0 {
		return errors.New("Miss param: mobile")
	}
	if req.TPLID == 0 {
		return errors.New("Miss param: tpl_id")
	}
	if len(req.TPLValue) == 0 {
		return errors.New("Miss param: tpl_value")
	}

	return nil
}

// TPLSingleSend 指定模板单发接口
func (sms *SMS) TPLSingleSend(input *TPLSingleSendRequest) (*SingleSendResponse, error) {
	if input == nil {
		input = &TPLSingleSendRequest{}
	}

	var result SingleSendResponse
	err := sms.sendRequest(http.MethodPost, "/v2/sms/tpl_single_send.json", input, &result)
	return &result, err
}

// TPLBatchSendRequest 指定模板群发请求参数
type TPLBatchSendRequest struct {
	Mobile   string `url:"mobile,omitempty"`
	TPLID    int64  `url:"tpl_id,omitempty"`
	TPLValue string `url:"tpl_value,omitempty"`
	Extend   string `url:"extend,omitempty"`
	UID      string `url:"uid,omitempty"`
}

// Verify 检查请求参数的有效性
func (req *TPLBatchSendRequest) Verify() error {
	if len(req.Mobile) == 0 {
		return errors.New("Miss param: mobile")
	}
	if req.TPLID == 0 {
		return errors.New("Miss param: tpl_id")
	}
	if len(req.TPLValue) == 0 {
		return errors.New("Miss param: tpl_value")
	}

	return nil
}

// TPLBatchSend 指定模板群发
func (sms *SMS) TPLBatchSend(input *TPLBatchSendRequest) (*BatchSendResponse, error) {
	if input == nil {
		input = &TPLBatchSendRequest{}
	}

	var result BatchSendResponse
	err := sms.sendRequest(http.MethodPost, "/v2/sms/tpl_batch_send.json", input, &result)
	return &result, err
}

func (sms *SMS) sendRequest(method, path string, input clientRequest, output interface{}) error {
	err := input.Verify()
	if err != nil {
		return err
	}

	r := sms.newRequest(method, sms.config.smsHost, path)
	r.header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	reader, err := sms.encodeFormBody(input)
	if err != nil {
		return err
	}
	r.body = reader

	resp, err := sms.doRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = sms.handleResponse(resp, output); err != nil {
		return err
	}

	return nil
}
