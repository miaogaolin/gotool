package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"github.com/miaogaolin/gotool/db/redis"
	"github.com/miaogaolin/gotool/errorx"
	"github.com/miaogaolin/gotool/limit"
	"time"

	"github.com/miaogaolin/gotool/logx"
)

var (
	DefaultRetryNum    = 3
	retryDuration      = time.Second
	clientTimeout      = time.Second * 5
	Client             *Rest
	periodLimitTimeout = time.Second * 5

	ErrRedisClose           = errorx.New("redis is disabled")
	ErrorPeriodLimitTimeout = errors.New("proxy take timeout")
)

type (
	Option       func(proxyName string, periodLimits map[string]*limit.PeriodLimit)
	ProxyFunc    func(*http.Request) (*url.URL, error)
	callbackFunc func(response *http.Response, err error) error
	Request      struct {
		ProxyName string
		Url       string
		Headers   map[string]string
		Body      io.Reader
		Cookies   []*http.Cookie
	}

	ProxyPac struct {
		Pac      []string
		Func     ProxyFunc
		Callback callbackFunc
	}

	Rest struct {
		proxyLimits map[string]*limit.PeriodLimit
		proxies     map[string]ProxyPac
	}
)

func New() *Rest {
	Client = &Rest{
		proxyLimits: make(map[string]*limit.PeriodLimit),
		proxies:     make(map[string]ProxyPac),
	}

	return Client
}

func (c *Rest) SetProxy(proxyName string, proxyPac ProxyPac, opts ...Option) *Rest {
	c.proxies[proxyName] = proxyPac
	for i := range opts {
		opts[i](proxyName, c.proxyLimits)
	}
	return c
}

func (c *Rest) GetProxy(proxyName string) *ProxyPac {
	if v, ok := c.proxies[proxyName]; ok {
		return &v
	}
	return nil
}

func (c *Rest) Pac(proxyName, url string) (ProxyPac, error) {
	res := ProxyPac{}
	proxyPac, err := c.getProxyLimit(proxyName)
	if err != nil {
		return res, err
	}
	if proxyPac == nil || proxyPac.Func == nil {
		return res, nil
	}

	for _, v := range proxyPac.Pac {
		if strings.Contains(url, v) {
			return *proxyPac, nil
		}
	}
	return res, nil
}

func (c *Rest) PostForm(r Request) (*http.Response, error) {
	proxyPac, err := c.Pac(r.ProxyName, r.Url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxyPac.Func,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: clientTimeout,
	}

	req, err := http.NewRequest(http.MethodPost, r.Url, r.Body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	return c.do(client, req, proxyPac.Callback)
}

func (c *Rest) Get(r Request) (*http.Response, error) {
	proxyPac, err := c.Pac(r.ProxyName, r.Url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: proxyPac.Func,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: clientTimeout,
	}

	req, err := http.NewRequest(http.MethodGet, r.Url, r.Body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}
	return c.do(client, req, proxyPac.Callback)
}

func (c *Rest) PostFile(r Request) (*http.Response, error) {
	proxyPac, err := c.Pac(r.ProxyName, r.Url)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy: proxyPac.Func,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   clientTimeout,
	}
	req, err := http.NewRequest(http.MethodPost, r.Url, r.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	for _, c := range r.Cookies {
		req.AddCookie(c)
	}

	return c.do(client, req, proxyPac.Callback)
}

func (c *Rest) getProxyLimit(proxyName string) (*ProxyPac, error) {
	if proxyName == "" {
		return nil, nil
	}

	prox := c.GetProxy(proxyName)
	if prox == nil {
		return nil, nil
	}

	l, ok := c.proxyLimits[proxyName]
	if !ok || l == nil {
		return prox, nil
	}

	for {
		timeout := time.After(periodLimitTimeout)
		select {
		case <-timeout:
			return prox, errorx.Wrapf(ErrorPeriodLimitTimeout, "timeout=%ds", periodLimitTimeout)
		default:
			code, err := l.Take(proxyName)
			if err != nil {
				logx.Errorf("PeriodLimit Take: %v", err)
				return nil, err
			}

			switch code {
			case limit.OverQuota:
				logx.Errorf("OverQuota key: %v", proxyName)
				<-time.After(time.Millisecond * 100)
				continue
			case limit.Allowed:
				//logx.Infof("AllowedQuota key: %v", proxyName)
				return prox, nil
			case limit.HitQuota:
				//logx.Errorf("HitQuota key: %v", proxyName)
				return prox, nil
			default:
				logx.Errorf("DefaultQuota key: %v", proxyName)
				return nil, nil
			}
		}

	}
}

func (c *Rest) do(client *http.Client, req *http.Request, callback callbackFunc) (*http.Response, error) {
	resp, err := client.Do(req)
	if callback == nil {
		return resp, err
	}
	return resp, callback(resp, err)
}

func Retry(fun func() (*http.Response, error), isRetry func(*http.Response, error) bool) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < DefaultRetryNum; i++ {
		resp, err = fun()
		if isRetry(resp, err) {
			<-time.After(retryDuration)
			continue
		}
		break
	}
	return resp, err
}

func WithPeriodLimit(redis *redis.RedisDB, period, quota int) Option {
	if !redis.IsEnabled() {
		logx.Fatal(ErrRedisClose)
		return nil
	}
	if period < 1 || quota < 1 {
		return nil
	}
	return func(proxyName string, periodLimits map[string]*limit.PeriodLimit) {
		periodLimits[proxyName] = limit.NewPeriodLimit(period, quota, redis, "period_limit")
	}
}

func PostJson(url string, headers map[string]string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
