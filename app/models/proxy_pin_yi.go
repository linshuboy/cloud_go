package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/json"
	"io"
	"net/http"
)

type ProxyPinYi struct {
	orm.Model
	UserId string `gorm:"column:user_id"`
	Appkey string `gorm:"column:appkey"`
	orm.SoftDeletes
}

type WhiteListResponse struct {
	ExpireTime string `json:"expire_time"`
	ExpiresIn  int    `json:"expires_in"`
	ShareToken string `json:"share_token"`
}

type FlowBalanceResponse struct {
	Ret       int    `json:"ret"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int    `json:"timestamp"`
	RetData   struct {
		AllBuy        int `json:"all_buy"`
		Balance       int `json:"balance"`
		FlowUsedToday int `json:"flow_used_today"`
	} `json:"ret_data"`
}

func (c *ProxyPinYi) AddThisIpToWhiteList() string {
	return ""
}

func (c *ProxyPinYi) HasFlowBalance() bool {
	url := "https://pycn.yapi.py.cn/index/users/flow_balance?neek=" + c.UserId + "&appkey=" + c.Appkey
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	response := &FlowBalanceResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	return response.RetData.AllBuy > response.RetData.Balance
}
