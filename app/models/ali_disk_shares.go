package models

import (
	"bytes"
	"encoding/json"
	orm2 "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AliDiskShares struct {
	orm.Model
	ShareId  string              `gorm:"column:share_id"`
	Password string              `gorm:"column:password"`
	SyncFlag string              `gorm:"column:flag"`
	ProxyIp  ProxyIp             `gorm:"-"`
	Files    []*AliDiskShareFile `gorm:"foreignKey:AliDiskShareId;references:ShareId"`
	orm.SoftDeletes
}

func (t *AliDiskShares) RandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (t *AliDiskShares) GetHttpClient() *http.Client {
	if t.ProxyIp.Ip != "" {
		proxyStr := "http://" + t.ProxyIp.Ip + ":" + strconv.Itoa(t.ProxyIp.Port)
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		return &http.Client{
			Transport: transport,
		}
	}
	return &http.Client{}
}

func (t *AliDiskShares) GetToken() string {
	var cacheName string
	if t.ProxyIp.Ip != "" {
		cacheName = "aliyunpansharetoken:" + t.ProxyIp.Ip + ":" + t.ShareId
	} else {
		cacheName = "aliyunpansharetoken:" + t.ShareId
	}
	value := facades.Cache().Get(cacheName, func() any {
		//https://api.aliyundrive.com/v2/share_link/get_share_token
		url1 := "https://api.aliyundrive.com/v2/share_link/get_share_token"
		// 准备请求体数据
		requestData := map[string]interface{}{
			"share_id":  t.ShareId,
			"share_pwd": t.Password,
		}
		// 将请求体数据编码为JSON
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			panic(err) // 处理JSON编码错误
		}
		req, err := http.NewRequest("POST", url1, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := t.GetHttpClient().Do(req)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(resp.Body)
		response := &GetShareTokenResponse{}
		if err := json.Unmarshal(body, response); err != nil {
			panic(err)
		}
		facades.Cache().Add(cacheName, response, time.Duration(response.ExpiresIn-300)*time.Second)
		return facades.Cache().Get(cacheName)
	})
	response1 := &GetShareTokenResponse{}
	if err := response1.UnmarshalBinary([]byte(value.(string))); err != nil {
		panic(err)
	}
	return response1.ShareToken
}

func (t *AliDiskShares) SetProxyIp(ip ProxyIp) bool {
	t.ProxyIp = ip
	return true
}

func (t *AliDiskShares) GetProxyIp() ProxyIp {
	var ipPool IpPool
	return ipPool.GetIp(t.ShareId, "get_share_file_list")
}

func (t *AliDiskShares) InitProxyIp() bool {
	return t.SetProxyIp(t.GetProxyIp())
}

func (t *AliDiskShares) DispatchesEvents() map[orm2.EventType]func(orm2.Event) error {
	return map[orm2.EventType]func(orm2.Event) error{
		orm2.EventCreated: func(event orm2.Event) error {
			_ = event
			share := &AliDiskShares{}
			_ = event.Query().Where("share_id = ?", event.GetAttribute("share_id")). /*.With("Files")*/ Find(&share)
			if len(share.Files) == 0 {
				var file AliDiskShareFile
				file.FileId = "root"
				file.AliDiskShareId = share.ShareId
				file.Type = "folder"
				_ = facades.Orm().Query().Save(&file)
			}
			return nil
		},
	}
}

func (t *AliDiskShares) GetFileListByFileId(fileId string, marker string) ([]AliDiskShareFile, string) {
	url1 := "https://api.aliyundrive.com/adrive/v2/file/list_by_share"
	// 准备请求体数据
	requestData := map[string]interface{}{
		"limit":           100,
		"order_by":        "name",
		"order_direction": "DESC",
		"parent_file_id":  fileId,
		"share_id":        t.ShareId,
		"marker":          marker,
	}
	// 将请求体数据编码为JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(err) // 处理JSON编码错误
	}
	req, err := http.NewRequest("POST", url1, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Share-Token", t.GetToken())
	resp, err := t.GetHttpClient().Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	response := &ListByShareResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		panic(err)
	}
	facades.Log().Info(string(body))
	facades.Log().Info(response)
	for _, item := range response.Items {
		facades.Log().Info(item.DriveId)
		facades.Log().Info(item.DomainId)
	}
	return response.Items, response.NextMarker
}

type ListByShareResponse struct {
	Items      []AliDiskShareFile `json:"items"`
	NextMarker string             `json:"next_marker"`
}

type AliDiskShareFile struct {
	DriveId        string          `gorm:"column:drive_id" json:"drive_id"`
	DomainId       string          `gorm:"column:domain_id" json:"domain_id"`
	FileId         string          `gorm:"column:file_id" json:"file_id"`
	AliDiskShareId string          `gorm:"column:share_id" json:"share_id"`
	Name           string          `gorm:"column:name" json:"name"`
	Type           string          `gorm:"column:type" json:"type"`
	ParentFileId   string          `gorm:"column:parent_file_id" json:"parent_file_id"`
	NextMarker     string          `gorm:"column:next_marker" json:"next_marker"`
	Json           string          `gorm:"column:json"`
	CompletedAt    carbon.DateTime `gorm:"column:completed_at"`
	Share          AliDiskShares   `gorm:"foreignKey:ShareId;references:AliDiskShareId"`
	orm.Model
	orm.SoftDeletes
}

type GetShareTokenResponse struct {
	ExpireTime string `json:"expire_time"`
	ExpiresIn  int    `json:"expires_in"`
	ShareToken string `json:"share_token"`
}

func (p GetShareTokenResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *GetShareTokenResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
