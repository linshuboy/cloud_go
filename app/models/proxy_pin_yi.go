package models

import (
	"bytes"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/json"
	"golang.org/x/exp/rand"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type ProxyPinYi struct {
	orm.Model
	UserId   string `gorm:"column:user_id"`
	Account  string `gorm:"column:account"`
	Password string `gorm:"column:password"`
	orm.SoftDeletes
}

type UserInfoResponse struct {
	PinyiResponse
	RetData struct {
		FlowAllBuy    int `json:"flow_all_buy"`
		FlowBalance   int `json:"flow_balance"`
		FlowUsedToday int `json:"flow_used_today"`
	} `json:"ret_data"`
}

func (p *ProxyPinYi) GetLoginToken() string {
	url := "https://pycn.yapi.py.cn/login"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("phone", p.Account)
	_ = writer.WriteField("password", p.Password)
	_ = writer.WriteField("remember", "0")
	_ = writer.Close()
	method := "POST"
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := (&http.Client{}).Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	response := &LoginResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	return response.RetData.Token
}

func (p *ProxyPinYi) GetToken() string {
	var cacheName = "pinyi_token:" + p.UserId
	value := facades.Cache().Get(cacheName, func() any {
		facades.Cache().Add(cacheName, p.GetLoginToken(), time.Duration(600)*time.Second)
		return facades.Cache().Get(cacheName)
	})
	return value.(string)
}

type IpifyResponse struct {
	Ip string `json:"ip"`
}

func (p *ProxyPinYi) ThisIp() string {
	url := "https://api.ipify.org?format=json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "api.ipify.org")
	req.Header.Add("Connection", "keep-alive")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	facades.Log().Info(string(body), "ipipip")
	response := &IpifyResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	return response.Ip
}

type LoginResponse struct {
	PinyiResponse
	RetData struct {
		Token string `json:"token"`
	} `json:"ret_data"`
}

type PinyiResponse struct {
	Ret       int    `json:"ret"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int    `json:"timestamp"`
}

func (p *ProxyPinYi) AddThisIpToWhiteList() bool {
	thisIp := p.ThisIp()
	if !p.AddWhiteList(thisIp) {
		return false
	}
	return true
}

type SaveWhiteResponse struct {
	PinyiResponse
}

func (p *ProxyPinYi) AddWhiteList(ip string) bool {
	url := "https://pycn.yapi.py.cn/user/save_white_ip"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("ip", ip)
	_ = writer.Close()
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("Authorization", "Bearer "+p.GetToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, _ := (&http.Client{}).Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	response := &SaveWhiteResponse{}
	facades.Log().Info(string(body), "add white ip")
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	if response.Ret == 0 && response.Code == 1 {
		return true
	}
	return false
}

func (p *ProxyPinYi) HasFlowBalance() bool {
	url := "https://pycn.yapi.py.cn/user/user_info"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+p.GetToken())
	res, _ := (&http.Client{}).Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, _ := io.ReadAll(res.Body)
	response := &UserInfoResponse{}
	facades.Log().Info(string(body), "HasFlowBalance")
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	return response.RetData.FlowAllBuy > response.RetData.FlowBalance
}

type ProxyIp struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	ExpireTime string `json:"expire_time"`
}

func (p ProxyIp) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProxyIp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

type GetIpResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Success bool      `json:"success"`
	Data    []ProxyIp `json:"data"`
}

func (p *ProxyPinYi) GetProxyIp(quantity int) []ProxyIp {
	url := "http://zltiqu.pyhttp.taolop.com/getip?count=" + strconv.Itoa(quantity) + "&neek=" + p.UserId + "&type=2&yys=100026&port=2&sb=&mr=1&sep=0&ts=1&ys=1&cs=1"
	method := "GET"
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "zltiqu.pyhttp.taolop.com")
	req.Header.Add("Connection", "keep-alive")
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//facades.Log().Info(string(body))
	response := &GetIpResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	return response.Data
}

type IpPool struct {
	Ips []string `json:"ips"`
}

func (p *IpPool) GetIp(uniqueIdentification string, operationType string) ProxyIp {
	ip := p.GetRandomIp()
	cacheName := "pinyi_proxy_ip:" + ip
	var proxyIp ProxyIp
	err := proxyIp.UnmarshalBinary([]byte(facades.Cache().Get(cacheName).(string)))
	if err != nil {
		panic(err)
	}
	return proxyIp
}

func (p *IpPool) SetIpPool(ips []ProxyIp) error {
	cacheName1 := "pinyi_proxy_ips"
	allIps := facades.Cache().Get(cacheName1, func() any {
		var ipPool = &IpPool{}
		for _, ip := range ips {
			ipPool.Ips = append(ipPool.Ips, ip.Ip)
		}
		facades.Cache().Add(cacheName1, ipPool, time.Duration(3600)*time.Second)
		return facades.Cache().Get(cacheName1)
	})
	if err := p.UnmarshalBinary([]byte(allIps.(string))); err != nil {
		panic(err)
	}
	for _, ip := range ips {
		p.Ips = append(p.Ips, ip.Ip)
	}
	p.Ips = p.RemoveRepByLoop(p.Ips)
	err := facades.Cache().Put(cacheName1, p, time.Duration(3600)*time.Second)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *IpPool) GetIpPool() *IpPool {
	cacheName := "pinyi_proxy_ips"
	if !facades.Cache().Has(cacheName) {
		return p
	}
	allIps := facades.Cache().Get(cacheName)
	if err := p.UnmarshalBinary([]byte(allIps.(string))); err != nil {
		panic(err)
	}
	return p
}

func (p *IpPool) GetRandomIp() string {
	cacheName1 := "pinyi_proxy_ips"
	cacheName := "pinyi_proxy_ip:"
	if len(p.Ips) == 0 {
		p.GetIpPool()
	}
	if len(p.Ips) == 0 {
		return ""
	}
	rand.Seed(uint64(time.Now().Nanosecond()))
	for {
		index := rand.Intn(len(p.Ips))
		ip := p.Ips[index]
		if facades.Cache().Has(cacheName + ip) {
			return ip
		} else {
			p.Ips = append(p.Ips[:index], p.Ips[index+1:]...)
			_ = facades.Cache().Put(cacheName1, p, time.Duration(3600)*time.Second)
		}
		if len(p.Ips) == 0 {
			break
		}
	}
	return ""
}

func (p *IpPool) RemoveInvalidIp() {
	cacheName1 := "pinyi_proxy_ips"
	cacheName := "pinyi_proxy_ip:"
	if len(p.Ips) == 0 {
		p.GetIpPool()
	}
	var newIps []string
	for _, ip := range p.Ips {
		if facades.Cache().Has(cacheName + ip) {
			newIps = append(newIps, ip)
		}
	}
	p.Ips = newIps
	_ = facades.Cache().Put(cacheName1, p, time.Duration(3600)*time.Second)
}

func (p *IpPool) RemoveRepByLoop(ips []string) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range ips {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func (p IpPool) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IpPool) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *ProxyPinYi) AddProxyToIpPool() bool {
	ips := p.GetProxyIp(1000)
	var ipPool1 = &IpPool{}
	_ = ipPool1.SetIpPool(ips)
	for _, ip := range ips {
		cacheName := "pinyi_proxy_ip:" + ip.Ip
		if !facades.Cache().Has(cacheName) {
			facades.Cache().Add(cacheName, ip, time.Duration(carbon.Parse(ip.ExpireTime).DiffAbsInSeconds()-30)*time.Second)
		} else {
			err := facades.Cache().Put(cacheName, ip, time.Duration(carbon.Parse(ip.ExpireTime).DiffAbsInSeconds()-30)*time.Second)
			if err != nil {
				panic(err)
			}
		}
	}
	return true
}
