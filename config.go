package yunpian

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Config 是云片sdk的相关配置项
type Config struct {
	UseSSL     *bool
	HTTPClient *http.Client
	APIKey     *string
	Context    context.Context

	smsHost string
}

// WithAPIKey 设置sdk的API key
func (c *Config) WithAPIKey(key string) *Config {
	c.APIKey = &key
	return c
}

// WithUseSSL 设置调用API时是否使用HTTPS
func (c *Config) WithUseSSL(use bool) *Config {
	c.UseSSL = &use
	return c
}

// WithHTTPClient 设置发送请求的Client
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// MergeIn 合并多个配置
func (c *Config) MergeIn(cfgs ...*Config) {
	for _, other := range cfgs {
		mergeInConfig(c, other)
	}
}

func mergeInConfig(dst *Config, other *Config) {
	if other == nil {
		return
	}

	if other.APIKey != nil {
		dst.APIKey = other.APIKey
	}
	if other.UseSSL != nil {
		dst.UseSSL = other.UseSSL
	}
	if other.HTTPClient != nil {
		dst.HTTPClient = other.HTTPClient
	}
	if other.Context != nil {
		dst.Context = other.Context
	}
}

// DefaultConfig 返回默认的sdk配置
func DefaultConfig() *Config {
	cfg := &Config{
		smsHost: "sms.yunpian.com",
	}
	return cfg.WithUseSSL(true).WithHTTPClient(defaultHTTPClient())
}

func defaultHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{Transport: transport}
}
